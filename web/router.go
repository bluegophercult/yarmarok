package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
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
					r.Put("/", router.updateParticipant)
					r.Get("/", router.listParticipants)
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
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	raffleService := r.organizerService.RaffleService(organizerID)

	m := newNoRequestMethodHandler(raffleService.List, r.logger.Logger)

	m.ServeHTTP(w, req)
}

func (r *Router) downloadRaffleXLSX(w http.ResponseWriter, req *http.Request) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	raffleService := r.organizerService.RaffleService(organizerID)

	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := raffleService.Export(raffleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	participantService, err := r.getParticipantService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMethodHandler(participantService.Create, r.logger.Logger).ServeHTTP(w, req)
}

func (r *Router) updateParticipant(w http.ResponseWriter, req *http.Request) {
	participantService, err := r.getParticipantService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMethodHandler(participantService.Edit, r.logger.Logger).ServeHTTP(w, req)
}

func (r *Router) listParticipants(w http.ResponseWriter, req *http.Request) {
	participantService, err := r.getParticipantService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newNoRequestMethodHandler(participantService.List, r.logger.Logger).ServeHTTP(w, req)
}

func (r *Router) getParticipantService(req *http.Request) (service.ParticipantService, error) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		return nil, err
	}

	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil || raffleID == "" {
		return nil, ErrMissingID
	}

	participantService := r.organizerService.RaffleService(organizerID).ParticipantService(raffleID)

	return participantService, nil
}

func extractOrganizerID(r *http.Request) (id string, err error) {
	defer func() {
		if localRun && err != nil {
			err = nil
			id = "dummy_test_user"
		}
	}()

	ids := r.Header.Values(GoogleUserIDHeader)

	if len(ids) != 1 {
		return "", ErrAmbiguousOrganizerIDHeader
	}

	id = ids[0]
	if id == "" {
		return "", ErrAmbiguousOrganizerIDHeader
	}

	return id, nil
}

func extractParam(req *http.Request, param string) (string, error) {
	val := chi.URLParam(req, param)
	if param == "" {
		return "", fmt.Errorf("missing param: %s", param)
	}

	return val, nil
}
