package web

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
)

// ErrAmbiguousUserIDHeader is returned when
// the user id header is not set or is ambiguous.
var ErrAmbiguousUserIDHeader = errors.New("ambiguous user id format")

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

	router.Post("/create-yarmarok", router.createYarmarok)
	router.Get("/list-yarmaroks", router.listYarmaroks)

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
