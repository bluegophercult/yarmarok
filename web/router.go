package web

import (
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
)

const (
	YarmaroksPath    = "/yarmaroks"
	ParticipantsPath = "/participants"
)

const (
	yarmarokIDParam    = "yarmarok_id"
	participantIDParam = "participant_id"
)

const (
	yarmarokIDPlaceholder = "{" + yarmarokIDParam + "}"
)

var (
	// ErrAmbiguousUserIDHeader is returned when
	// the user id header is not set or is ambiguous.
	ErrAmbiguousUserIDHeader = errors.New("ambiguous user id format")

	// ErrMissingID is returned when id is missing.
	ErrMissingID = errors.New("missing id")
)

// Router is responsible for routing requests
// to the corresponding services.
type Router struct {
	chi.Router
	userService service.UserService
	logger      *logger.Entry
}

// NewRouter creates a new Router
func NewRouter(us service.UserService, log *logger.Logger) (*Router, error) {
	router := &Router{
		Router:      chi.NewRouter(),
		userService: us,
		logger: log.WithFields(
			logger.Fields{
				"component": "router",
				"trace_id":  uuid.New().String(),
			},
		),
	}

	router.Use(router.loggingMiddleware)
	router.Use(router.recoverMiddleware)
	router.Use(router.userMiddleware)

	router.Route(YarmaroksPath, func(subRouter chi.Router) {
		subRouter.Post("/", router.createYarmarok)
		subRouter.Get("/", router.listYarmaroks)
	})

	router.Route(joinPath(YarmaroksPath, yarmarokIDPlaceholder, ParticipantsPath), func(subRouter chi.Router) {
		subRouter.Post("/", router.createParticipant)
		subRouter.Put("/", router.updateParticipant)
		subRouter.Get("/", router.listParticipants)
	})

	return router, nil
}

func (r *Router) createYarmarok(w http.ResponseWriter, req *http.Request) {
	userID, err := extractUserID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	yarmarokService := r.userService.YarmarokService(userID)

	m := newMethodHandler(yarmarokService.Init, r.logger.Logger)

	m.ServeHTTP(w, req)
}

func (r *Router) listYarmaroks(w http.ResponseWriter, req *http.Request) {
	userID, err := extractUserID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	yarmarokService := r.userService.YarmarokService(userID)

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

func joinPath(args ...string) string {
	return path.Clean("/" + path.Join(args...))
}

func (r *Router) getParticipantService(req *http.Request) (service.ParticipantService, error) {
	userID, err := extractUserID(req)
	if err != nil {
		return nil, err
	}

	yarmarokID, err := extractParam(req, yarmarokIDParam)
	if err != nil {
		return nil, ErrMissingID
	}

	participantService := r.userService.YarmarokService(userID).ParticipantService(yarmarokID)

	return participantService, nil
}

func extractUserID(r *http.Request) (string, error) {
	ids := r.Header.Values(GoogleUserIDHeader)

	if len(ids) != 1 {
		return "", ErrAmbiguousUserIDHeader
	}

	id := ids[0]
	if id == "" {
		return "", ErrAmbiguousUserIDHeader
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
