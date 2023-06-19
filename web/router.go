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

const (
	YarmaroksPath    = "/yarmaroks"
	ParticipantsPath = "/participants"

	yarmarokIDParam    = "yarmarok_id"
	participantIDParam = "participant_id"

	yarmarokIDPlaceholder = "/{" + yarmarokIDParam + "}"
)

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
	router.Use(router.organizerMiddleware)

	router.Route(
		YarmaroksPath,
		func(yarmaroksRouter chi.Router) { // "/yarmaroks"
			yarmaroksRouter.Post("/", router.createYarmarok)
			yarmaroksRouter.Get("/", router.listYarmaroks)
			yarmaroksRouter.Route(
				yarmarokIDPlaceholder,
				func(yarmarokIDRouter chi.Router) { // "/yarmaroks/{yarmarok_id}"
					yarmarokIDRouter.Route(
						ParticipantsPath,
						func(participantsRouter chi.Router) { // "/yarmaroks/{yarmarok_id}/participants"
							participantsRouter.Post("/", router.createParticipant)
							participantsRouter.Put("/", router.updateParticipant)
							participantsRouter.Get("/", router.listParticipants)
						},
					)
				},
			)
		},
	)

	return router, nil
}

func (r *Router) createYarmarok(w http.ResponseWriter, req *http.Request) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	yarmarokService := r.organizerService.YarmarokService(organizerID)

	m := newMethodHandler(yarmarokService.Init, r.logger.Logger)

	m.ServeHTTP(w, req)
}

func (r *Router) listYarmaroks(w http.ResponseWriter, req *http.Request) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	yarmarokService := r.organizerService.YarmarokService(organizerID)

	m := newNoRequestMethodHandler(yarmarokService.List, r.logger.Logger)

	m.ServeHTTP(w, req)
}

func (r *Router) createParticipant(w http.ResponseWriter, req *http.Request) {
	participantService, err := r.getParticipantService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMethodHandler(participantService.Add, r.logger.Logger).ServeHTTP(w, req)
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

	yarmarokID, err := extractParam(req, yarmarokIDParam)
	if err != nil || yarmarokID == "" {
		return nil, ErrMissingID
	}

	participantService := r.organizerService.YarmarokService(organizerID).ParticipantService(yarmarokID)

	return participantService, nil
}

func extractOrganizerID(r *http.Request) (string, error) {
	ids := r.Header.Values(GoogleOrganizerIDHeader)

	if len(ids) != 1 {
		return "", ErrAmbiguousOrganizerIDHeader
	}

	id := ids[0]
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
