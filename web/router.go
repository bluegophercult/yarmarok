package web

import (
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
	raffleGroup      = ApiPath + RafflesPath
	participantGroup = raffleGroup + raffleIDPlaceholder + ParticipantsPath
	prizeGroup       = raffleGroup + raffleIDPlaceholder + PrizesPath
	donationGroup    = prizeGroup + prizeIDPlaceholder + DonationsPath
)

const (
	raffleIDPlaceholder      = "/{" + raffleIDParam + "}"
	participantIDPlaceholder = "/{" + participantIDParam + "}"
	prizeIDPlaceholder       = "/{" + prizeIDParam + "}"
	donationIDPlaceholder    = "/{" + donationIDParam + "}"
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
	r := &Router{
		Router:           chi.NewRouter(),
		organizerService: os,
		logger: log.WithFields(
			logger.Fields{
				"component": "router",
				"trace_id":  uuid.New().String(),
			},
		),
	}

	r.Use(r.corsMiddleware)
	r.Use(r.loggingMiddleware)
	r.Use(r.recoverMiddleware)
	r.Use(r.headerMiddleware)
	r.Use(r.organizerMiddleware)

	r.Map(loginRoute)
	r.Map(raffleRoute)
	r.Map(participantRoute)
	r.Map(prizeRoute)
	r.Map(donationRoute)

	return r, nil
}

func (r *Router) Map(route func(*Router)) { route(r) }

func loginRoute(r *Router) {
	r.Route(ApiPath, func(r chi.Router) {
		r.Handle("/login", http.RedirectHandler("/", http.StatusSeeOther))
	})
}

func raffleRoute(r *Router) {
	raffle := newService[service.RaffleService, *service.RaffleRequest, service.Raffle](r.getRaffleService)

	r.Route(raffleGroup, func(router chi.Router) {
		router.Post("/", raffle.Create)
		router.Get("/", raffle.List)
		router.Route(raffleIDPlaceholder, func(router chi.Router) {
			router.Put("/", raffle.Edit)
			router.Get("/download-xlsx", r.downloadRaffleXLSX)
		})
	})
}

func participantRoute(r *Router) {
	participant := newService[service.ParticipantService, *service.ParticipantRequest, service.Participant](r.getParticipantService)

	r.Route(participantGroup, func(router chi.Router) {
		router.Post("/", participant.Create)
		router.Get("/", participant.List)
		router.Route(participantIDPlaceholder, func(router chi.Router) {
			router.Put("/", participant.Edit)
			router.Delete("/", participant.Delete)
		})
	})
}

func prizeRoute(r *Router) {
	prize := newService[service.PrizeService, *service.PrizeRequest, service.Prize](r.getPrizeService)

	r.Route(prizeGroup, func(router chi.Router) {
		router.Post("/", prize.Create)
		router.Get("/", prize.List)
		router.Route(prizeIDPlaceholder, func(router chi.Router) {
			router.Get("/", prize.Get)
			router.Put("/", prize.Edit)
			router.Delete("/", prize.Delete)
		})
	})
}

func donationRoute(r *Router) {
	donation := newService[service.DonationService, *service.DonationRequest, service.Donation](r.getDonationService)

	r.Route(donationGroup, func(router chi.Router) {
		router.Post("/", donation.Create)
		router.Get("/", donation.List)
		router.Route(donationIDPlaceholder, func(router chi.Router) {
			router.Get("/", donation.Get)
			router.Put("/", donation.Edit)
			router.Delete("/", donation.Delete)
		})
	})
}

/* Handlers */

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
