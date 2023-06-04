package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
)

const GoogleUserIDHeader = "X-Goog-Authenticated-User-Id"

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

	return router, nil
}

func (r *Router) createYarmarok(w http.ResponseWriter, req *http.Request) {
	userID, err := extractUserID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	yarmarokService := r.userService.YarmarokService(userID)

	m := methodHandler[
		*service.YarmarokInitRequest,
		*service.InitResult,
	](
		yarmarokService.Init,
	).WithLogger(r.logger)

	m.ServeHTTP(w, req)
}

// noRequestMethodHandler is a wrapper around a service method
// that converts a no request method to an http handler.
type noRequestMethodHandler[Response any] func() (Response, error)

func (m noRequestMethodHandler[Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	emptyRequestWrapper := func(struct{}) (Response, error) {
		return m()
	}

	h := methodHandler[
		struct{},
		Response,
	](
		emptyRequestWrapper,
	)

	h.ServeHTTP(w, req)
}

// methodHandler is a wrapper around a service method
// that converts that method to an http handler.
type methodHandler[Request any, Response any] func(Request) (Response, error)

func (m methodHandler[Request, Response]) WithLogger(log *logger.Entry) methodHandler[Request, Response] {
	return func(req Request) (Response, error) {
		log.WithFields(
			logger.Fields{
				"request": req,
			},
		).Debug("request")

		resp, err := m(req)
		if err != nil {
			log.WithError(err).Error("request failed")
		}

		log.WithFields(
			logger.Fields{
				"response": resp,
			},
		).Debug("response")

		return resp, err
	}
}

func (m methodHandler[Request, Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := new(Request)

	err := json.NewDecoder(req.Body).Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := m(*request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (r *Router) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, err := extractUserID(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = r.userService.InitUserIfNotExists(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, _ := extractUserID(req)

		start := time.Now()
		duration := time.Since(start)

		lrw := logger.NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, req)

		responseMetric := lrw.ResponseMetric()

		r.logger.WithFields(
			logger.Fields{
				"uri":      req.RequestURI,
				"method":   req.Method,
				"status":   responseMetric.Status,
				"duration": duration,
				"size":     responseMetric.Size,
				"user_id":  userID,
			},
		).Info("request completed")
	})
}

func (R *Router) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				R.logger.WithFields(logger.Fields{
					"uri":    req.RequestURI,
					"method": req.Method,
					"error":  err,
				}).Error("panic recovered")
			}
			w.WriteHeader(http.StatusInternalServerError)
		}()

		next.ServeHTTP(w, req)
	})
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
