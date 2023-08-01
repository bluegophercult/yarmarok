package web

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/kaznasho/yarmarok/service"
	"github.com/rs/cors"

	"github.com/kaznasho/yarmarok/logger"
)

const (
	// GoogleUserIDHeader is the header that contains the organizer id set by google identity aware proxy.
	GoogleUserIDHeader = "X-Goog-Authenticated-User-Id"

	defaultOrigin = "https://yarmarock.com.ua"
)

var (
	ErrCreatingOrganizer  = errors.New("creating organizer")
	ErrRecoveredFromPanic = errors.New("recovered from panic")
)

type Middleware = func(Handler) Handler

func WrapMiddlewares(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

// WithErrors middleware in the web package handles errors during HTTP handler execution.
// It checks if an error is of type Error, and if not, creates a default error.
// The error and its associated Log fields are then logged.
// Finally, the error is passed to the Respond function to render the error to the client with Message if there is any.
func WithErrors(log *logger.Logger) Middleware {
	return func(h Handler) Handler {
		return func(rw http.ResponseWriter, req *http.Request) error {
			if err := h(rw, req); err != nil {
				errweb, ok := ErrorAs(err)
				if !ok {
					errweb = NewError(err, http.StatusInternalServerError, Message{"error": ErrUnknownError})
				}

				log.WithError(err).WithFields(logger.Fields(errweb.log))

				return Respond(rw, errweb)
			}

			return nil
		}
	}
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
				return NewError(errors.Join(ErrMissingID, err), http.StatusBadRequest)
			}

			if err = svc.CreateOrganizerIfNotExists(organizerID); err != nil {
				return NewError(errors.Join(ErrCreatingOrganizer, err), http.StatusInternalServerError)
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

func WithRecover(h Handler) Handler {
	return func(rw http.ResponseWriter, req *http.Request) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = NewError(ErrRecoveredFromPanic, http.StatusInternalServerError,
					Log{
						"uri":    req.RequestURI,
						"method": req.Method,
						"rec":    rec,
						"trace":  string(debug.Stack()),
					})
			}
		}()

		return h(rw, req)
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
