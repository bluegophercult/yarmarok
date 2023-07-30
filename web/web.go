package web

import (
	"net/http"
	"path"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
)

type Handler func(rw http.ResponseWriter, req *http.Request) error

type Web struct {
	mux *chi.Mux
	mws []Middleware
	log *logger.Entry
	svc service.OrganizerService
}

func NewWeb(log *logger.Logger, svc service.OrganizerService) *Web {
	return &Web{
		mux: chi.NewRouter(),
		mws: []Middleware{
			WithOrganizer(svc, log),
			WithLogging(log),
			WithCORS,
			WithJSON,
			WithErrors(log),
			WithRecover(log),
		},
		log: log.WithFields(
			logger.Fields{
				"component": "web",
				"trace_id":  uuid.New().String(),
			},
		),
		svc: svc,
	}
}

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

func (w *Web) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w.mux.ServeHTTP(rw, req)
}
