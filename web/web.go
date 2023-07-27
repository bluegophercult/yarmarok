package web

import (
	"net/http"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/rs/cors"
)

type Handler func(rw http.ResponseWriter, req *http.Request) error

type Web struct {
	mux *chi.Mux
	mws []Middleware
	log *logger.Entry
	svc service.OrganizerService
}

type Middleware = func(Handler) Handler

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
	w.log.Println(req)
	w.mux.ServeHTTP(rw, req)
}

func WrapMiddlewares(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

func WithJSON(h Handler) Handler {
	return func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Content-Type", "application/json")

		return h(rw, req)
	}
}

func WithXLSX(h Handler) Handler {
	return func(rw http.ResponseWriter, req *http.Request) error {
		rw.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		return h(rw, req)
	}
}

func WithOrganizer(svc service.OrganizerService, log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(rw http.ResponseWriter, req *http.Request) error {
			organizerID, err := extractOrganizerID(req)
			if err != nil {
				log.WithError(err).Error("extract organizer id")
				return err
			}

			err = svc.CreateOrganizerIfNotExists(organizerID)
			if err != nil {
				log.WithError(err).Error("creating organizer")
				return err
			}

			return h(rw, req)
		}
	}
}

func WithLogging(log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(rw http.ResponseWriter, req *http.Request) error {
			start := time.Now()

			defer func() {
				log.WithFields(logger.Fields{
					"uri":          req.RequestURI,
					"method":       req.Method,
					"status":       rw.Header().Get("Status"),
					"duration":     time.Since(start),
					"size":         rw.Header().Get("Content-Length"),
					"organizer_id": req.Header.Get(GoogleUserIDHeader),
				}).Info("request completed")
			}()

			return h(rw, req)
		}
	}
}

func WithRecover(log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(rw http.ResponseWriter, req *http.Request) error {
			defer func() {
				if err := recover(); err != nil {
					log.WithFields(logger.Fields{
						"uri":    req.RequestURI,
						"method": req.Method,
						"error":  err,
					}).Error("panic recovered")
				}
			}()

			return h(rw, req)
		}
	}
}

func WithCORS(h Handler) Handler {
	return func(rw http.ResponseWriter, req *http.Request) error {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{defaultOrigin},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Goog-Authenticated-User-Id"},
			ExposedHeaders:   []string{},
			AllowCredentials: true,
			MaxAge:           0,
		})
		c.HandlerFunc(rw, req)
		return h(rw, req)
	}
}

// WithErrors is a middleware that handles errors if they occur during the execution.
// If the error is not of the type Error, as determined by the ErrorIs function, then a default error is created.
func WithErrors(log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(rw http.ResponseWriter, req *http.Request) error {
			if err := h(rw, req); err != nil {
				log.WithFields(logger.Fields{"error": err})

				if !ErrorIs(err) {
					err = NewError(err, http.StatusInternalServerError).With("error", ErrUnknownError)
				}

				return Respond(rw, err)
			}

			return nil
		}
	}
}
