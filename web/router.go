package web

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
)

const (
	ApiPath          = "/api"
	RafflesPath      = "/raffles"
	ParticipantsPath = "/participants"
	PrizesPath       = "/prizes"
	DonationsPath    = "/donations"
)

const (
	raffleIDParam      = "raffle_id"
	participantIDParam = "participant_id"
	prizeIDParam       = "prize_id"
	donationIDParam    = "donation_id"
)

const (
	raffleIDPlaceholder      = "/{" + raffleIDParam + "}"
	participantIDPlaceholder = "/{" + participantIDParam + "}"
	prizeIDPlaceholder       = "/{" + prizeIDParam + "}"
	donationIDPlaceholder    = "/{" + donationIDParam + "}"
)

// localRun is true if app is build for local run
var localRun = false

var (
	// ErrAmbiguousOrganizerIDHeader is returned when
	// the organizer id header is not set or is ambiguous.
	ErrAmbiguousOrganizerIDHeader = errors.New("ambiguous organizer id format")

	// ErrMissingID is returned when id is missing.
	ErrMissingID = errors.New("missing id")
)

// Router is responsible for routing requests
// to the corresponding services.
type Router struct {
	chi.Router
	organizerService service.OrganizerService
	logger           *logger.Entry
}

// NewRouter creates a new Router
func NewRouter(os service.OrganizerService, log *logger.Logger) (*Router, error) {
	router := &Router{
		Router:           chi.NewRouter(),
		organizerService: os,
		logger: log.WithFields(
			logger.Fields{
				"component": "router",
				"trace_id":  uuid.New().String(),
			},
		),
	}

	router.Use(router.corsMiddleware)
	router.Use(router.loggingMiddleware)
	router.Use(router.recoverMiddleware)
	router.Use(router.headerMiddleware)
	router.Use(router.organizerMiddleware)

	// "/api"
	router.Route(ApiPath, func(r chi.Router) {
		r.Handle("/login", http.RedirectHandler("/", http.StatusSeeOther))

		// "/api/raffles"
		raffle := newService[service.RaffleService, *service.RaffleRequest, service.Raffle](router.getRaffleService)
		r.Route(RafflesPath, func(r chi.Router) {
			r.Post("/", raffle.Create)
			r.Get("/", raffle.List)

			// "/api/raffles/{raffle_id}"
			r.Route(raffleIDPlaceholder, func(r chi.Router) {
				r.Put("/", raffle.Edit)
				r.Get("/download-xlsx", router.downloadRaffleXLSX)

				// "/api/raffles/{raffle_id}/participants"
				participant := newService[service.ParticipantService, *service.ParticipantRequest, service.Participant](router.getParticipantService)
				r.Route(ParticipantsPath, func(r chi.Router) {
					r.Post("/", participant.Create)
					r.Get("/", participant.List)

					// "/api/raffles/{raffle_id}/participants/{participant_id}"
					r.Route(participantIDPlaceholder, func(r chi.Router) {
						r.Put("/", participant.Edit)
						r.Delete("/", participant.Delete)
					})
				})

				// "/api/raffles/{raffle_id}/prizes"
				prize := newService[service.PrizeService, *service.PrizeRequest, service.Prize](router.getPrizeService)
				r.Route(PrizesPath, func(r chi.Router) {
					r.Post("/", prize.Create)
					r.Get("/", prize.List)

					// "/api/raffles/{raffle_id}/prizes/{prize_id}"
					r.Route(prizeIDPlaceholder, func(r chi.Router) {
						r.Get("/", prize.Get)
						r.Put("/", prize.Edit)
						r.Delete("/", prize.Delete)

            
						// "/api/raffles/{raffle_id}/prizes/{prize_id}/donations"
            donation := newService[service.DonationService, *service.DonationRequest, service.Donation](router.geDonationService)
						r.Route(DonationsPath, func(r chi.Router) {
							r.Post("/", donation.Create)
							r.Get("/", donation.List)

							// "/api/raffles/{raffle_id}/prizes/{prize_id}/donations/{donation_id}"
							r.Route(donationIDPlaceholder, func(r chi.Router) {
								r.Get("/", donation.Get)
								r.Put("/", donation.Edit)
								r.Delete("/", donation.Delete)
							})
						})
					})

				})
			})
		})
	})

	return router, nil
}

func (r *Router) downloadRaffleXLSX(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	id, err := extractParam(req, raffleIDParam)
	if err != nil {
		respondErr(w, err)
		return
	}

	res, err := svc.Export(id)
	if err != nil {
		respondErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+res.FileName)

	if _, err := w.Write(res.Content); err != nil {
		respondErr(w, err)
		r.logger.WithError(err).Error("writing xlsx")
	}
}
