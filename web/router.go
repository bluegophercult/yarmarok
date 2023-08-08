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
)

const (
	raffleIDParam      = "raffle_id"
	participantIDParam = "participant_id"
	prizeIDParam       = "prize_id"
)

const (
	raffleIDPlaceholder      = "/{" + raffleIDParam + "}"
	participantIDPlaceholder = "/{" + participantIDParam + "}"
	prizeIDPlaceholder       = "/{" + prizeIDParam + "}"
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
				r.Get("/download-xlsx", router.downloadRaffleXLSX)

				// "/api/raffles/{raffle_id}/participants"
				r.Route(ParticipantsPath, func(r chi.Router) {
					r.Post("/", router.createParticipant)
					r.Get("/", router.listParticipants)

					// "/api/raffles/{raffle_id}/participants/{participant_id}"
					r.Put(participantIDPlaceholder, router.editParticipant)
					r.Delete(participantIDPlaceholder, router.deleteParticipant)
				})

				// "/api/raffles/{raffle_id}/prizes"
				r.Route(PrizesPath, func(r chi.Router) {
					r.Post("/", router.createPrize)
					r.Get("/", router.listPrizes)

					// "/api/raffles/{raffle_id}/prizes/{prize_id}"
					r.Get(prizeIDPlaceholder, router.getPrize)
					r.Put(prizeIDPlaceholder, router.editPrize)
					r.Delete(prizeIDPlaceholder, router.deletePrize)
				})
			})
		})
	})

	return router, nil
}

func (r *Router) createRaffle(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newCreate(svc.Create).Handle(w, req)
}

func (r *Router) listRaffles(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getRaffleService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newList(svc.List).Handle(w, req)
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

func (r *Router) createParticipant(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newCreate(svc.Create).Handle(w, req)
}

func (r *Router) editParticipant(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newEdit(svc.Edit).Handle(w, req)
}

func (r *Router) deleteParticipant(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newDelete(svc.Delete).Handle(w, req)
}

func (r *Router) listParticipants(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getParticipantService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newList(svc.List).Handle(w, req)
}

func (r *Router) createPrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newCreate(svc.Create).Handle(w, req)
}

func (r *Router) getPrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newGet(svc.Get).Handle(w, req)
}

func (r *Router) editPrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newEdit(svc.Edit).Handle(w, req)
}

func (r *Router) deletePrize(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newDelete(svc.Delete).Handle(w, req)
}

func (r *Router) listPrizes(w http.ResponseWriter, req *http.Request) {
	svc, err := r.getPrizeService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newList(svc.List).Handle(w, req)
}
