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
	PlayPath         = "/play"
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
		r.Route(RafflesPath, func(r chi.Router) {
			r.Post("/", router.createRaffle)
			r.Get("/", router.listRaffles)

			// "/api/raffles/{raffle_id}"
			r.Route(raffleIDPlaceholder, func(r chi.Router) {
				r.Put("/", router.editRaffle)
				r.Delete("/", router.deleteRaffle)
				r.Get("/download-xlsx", router.downloadRaffleXLSX)

				// "/api/raffles/{raffle_id}/participants"
				r.Route(ParticipantsPath, func(r chi.Router) {
					r.Post("/", router.createParticipant)
					r.Get("/", router.listParticipants)

					// "/api/raffles/{raffle_id}/participants/{participant_id}"
					r.Route(participantIDPlaceholder, func(r chi.Router) {
						r.Put("/", router.editParticipant)
						r.Delete("/", router.deleteParticipant)
					})
				})

				// "/api/raffles/{raffle_id}/prizes"
				r.Route(PrizesPath, func(r chi.Router) {
					r.Post("/", router.createPrize)
					r.Get("/", router.listPrizes)

					// "/api/raffles/{raffle_id}/prizes/{prize_id}"
					r.Route(prizeIDPlaceholder, func(r chi.Router) {
						r.Get("/", router.getPrize)
						r.Put("/", router.editPrize)
						r.Delete("/", router.deletePrize)

						// "/api/raffles/{raffle_id}/prizes/{prize_id}/play"
						r.Route(PlayPath, func(r chi.Router) {
							r.Get("/", router.playPrize)
						})

						// "/api/raffles/{raffle_id}/prizes/{prize_id}/donations"
						r.Route(DonationsPath, func(r chi.Router) {
							r.Post("/", router.createDonation)
							r.Get("/", router.listDonations)

							// "/api/raffles/{raffle_id}/prizes/{prize_id}/donations/{donation_id}"
							r.Route(donationIDPlaceholder, func(r chi.Router) {
								r.Get("/", router.getDonation)
								r.Put("/", router.editDonation)
								r.Delete("/", router.deleteDonation)
							})
						})
					})

				})
			})
		})
	})

	return router, nil
}

func (r *Router) createRaffle(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewCreateHandler(r, svc.Create).Handle(w, req)
}

func (r *Router) editRaffle(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewEditHandler(r, svc.Edit).Handle(w, req)
}

func (r *Router) deleteRaffle(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewDeleteHandler(r, svc.Delete).Handle(w, req)
}

func (r *Router) listRaffles(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewListHandler(r, svc.List).Handle(w, req)
}

func (r *Router) downloadRaffleXLSX(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	id, err := extractParam(req, raffleIDParam)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	res, err := svc.Export(id)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+res.FileName)

	if _, err := w.Write(res.Content); err != nil {
		r.respondErr(w, err)
	}
}

func (r *Router) createParticipant(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewCreateHandler(r, svc.Create).Handle(w, req)
}

func (r *Router) editParticipant(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewEditHandler(r, svc.Edit).Handle(w, req)
}

func (r *Router) deleteParticipant(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewDeleteHandler(r, svc.Delete).Handle(w, req)
}

func (r *Router) listParticipants(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewListHandler(r, svc.List).Handle(w, req)
}

func (r *Router) createPrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewCreateHandler(r, svc.Create).Handle(w, req)
}

func (r *Router) getPrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewGetHandler(r, svc.Get).Handle(w, req)
}

func (r *Router) editPrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewEditHandler(r, svc.Edit).Handle(w, req)
}

func (r *Router) deletePrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewDeleteHandler(r, svc.Delete).Handle(w, req)
}

func (r *Router) listPrizes(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewListHandler(r, svc.List).Handle(w, req)
}

func (r *Router) playPrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewGetHandler(r, svc.Play).Handle(w, req)
}

func (r *Router) createDonation(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getDonationService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewCreateHandler(r, svc.Create).Handle(w, req)
}

func (r *Router) getDonation(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getDonationService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewGetHandler(r, svc.Get).Handle(w, req)
}

func (r *Router) listDonations(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getDonationService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewListHandler(r, svc.List).Handle(w, req)
}

func (r *Router) editDonation(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getDonationService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewEditHandler(r, svc.Edit).Handle(w, req)
}

func (r *Router) deleteDonation(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getDonationService(req)
	if err != nil {
		r.respondErr(w, err)
		return
	}

	NewDeleteHandler(r, svc.Delete).Handle(w, req)
}
