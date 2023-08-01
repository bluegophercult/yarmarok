package web

import (
	"net/http"
	"path"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
)

// Handler represents application handler function.
type Handler func(rw http.ResponseWriter, req *http.Request) error

// Web represents a controller for handling HTTP requests.
type Web struct {
	mux *chi.Mux
	mws []Middleware
	log *logger.Entry
	svc service.OrganizerService
}

// NewWeb creates a new Web instance.
func NewWeb(log *logger.Logger, svc service.OrganizerService, mws ...Middleware) *Web {
	return &Web{
		mux: chi.NewRouter(),
		mws: []Middleware{
			WithLogging(log),
			WithErrors(log),
			WithCORS,
			WithJSON,
			WithRecover,
			WithOrganizer(svc, log),
		},
		log: log.WithFields(logger.Fields{
			"component": "web",
			"trace_id":  uuid.New().String(),
		}),
		svc: svc,
	}
}

// Handle adds a new route to the Web's router. It takes an HTTP method, a group,
// a pattern, a Handler, and optional middleware. It wraps the Handler with the
// provided middleware and the Web's default middleware, then adds the route to the router.
func (w *Web) Handle(method string, group, pattern string, h Handler, mws ...Middleware) {
	h = WrapMiddlewares(h, append(w.mws, mws...)...)

	fn := func(rw http.ResponseWriter, req *http.Request) {
		if err := h(rw, req); err != nil {
			w.log.WithError(err).Error("handling request")
			return
		}
	}

	w.mux.MethodFunc(method, path.Join(group, pattern), fn)
}

// ServeHTTP makes the Web instance satisfy the http.Handler interface.
// It delegates the handling of the HTTP request to the Web's router.
func (w *Web) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w.mux.ServeHTTP(rw, req)
}
