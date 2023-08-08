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
)

const (
	raffleIDParam      = "raffle_id"
	participantIDParam = "participant_id"
)

const (
	raffleIDPlaceholder      = "/{" + raffleIDParam + "}"
	participantIDPlaceholder = "/{" + participantIDParam + "}"
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

	router.Route(ApiPath, func(r chi.Router) { // "/api"
		r.Handle("/login", http.RedirectHandler("/", http.StatusSeeOther))
		r.Route(RafflesPath, func(r chi.Router) { // "/api/raffles"
			r.Post("/", router.createRaffle)
			r.Get("/", router.listRaffles)
			r.Route(raffleIDPlaceholder, func(r chi.Router) { // "/api/raffles/{raffle_id}"
				r.Get("/download-xlsx", router.downloadRaffleXLSX)
				r.Route(ParticipantsPath, func(r chi.Router) { // "/api/raffles/{raffle_id}/participants"
					r.Post("/", router.createParticipant)
					r.Get("/", router.listParticipants)
					r.Put(participantIDPlaceholder, router.editParticipant)
					r.Delete(participantIDPlaceholder, router.deleteParticipant)
				})
			})
		})
	})

	return router, nil
}

func (r *Router) createRaffle(w http.ResponseWriter, req *http.Request) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	raffleService := r.organizerService.RaffleService(organizerID)

	m := newMethodHandler(raffleService.Create, r.logger.Logger)

	m.ServeHTTP(w, req)
}

func (r *Router) listRaffles(w http.ResponseWriter, req *http.Request) {
	raffleService, err := r.getRaffleService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	newNoRequestMethodHandler(raffleService.List, r.logger.Logger).ServeHTTP(w, req)
}

func (r *Router) downloadRaffleXLSX(w http.ResponseWriter, req *http.Request) {
	raffleService, err := r.getRaffleService(req)
	if err != nil {
		respondErr(w, err)
		return
	}

	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil {
		respondErr(w, err)
		return
	}

	resp, err := raffleService.Export(raffleID)
	if err != nil {
		respondErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+resp.FileName)

	if _, err := w.Write(resp.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		r.logger.WithError(err).Error("writing xlsx failed")
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
