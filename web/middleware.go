package web

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/kaznasho/yarmarok/service"
	"net/http"
	"time"

	"github.com/kaznasho/yarmarok/logger"
)

// GoogleUserIDHeader is the header that contains the user id
// set by google identity aware proxy.
const GoogleUserIDHeader = "X-Goog-Authenticated-User-Id"

func (r *Router) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, err := extractUserID(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			r.logger.WithError(err).Error("failed to extract user id")
			return
		}

		err = r.userService.InitUserIfNotExists(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			r.logger.WithError(err).Error("failed to init user")
			return
		}

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

func (r *Router) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				r.logger.WithFields(logger.Fields{
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

type ctxKey int

const participantServiceKey ctxKey = iota + 1

func (r *Router) participantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, err := extractUserID(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		yarmarokID := chi.URLParam(req, yarmarokIDPath)
		if yarmarokID == "" {
			http.Error(w, ErrMissingID.Error(), http.StatusBadRequest)
			return
		}

		participantService := r.userService.
			YarmarokService(userID).
			ParticipantService(yarmarokID)

		ctx := context.WithValue(req.Context(), participantServiceKey, participantService)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func (r *Router) getParticipantService(ctx context.Context) (service.ParticipantService, error) {
	val := ctx.Value(participantServiceKey)
	svc, ok := val.(service.ParticipantService)
	if !ok {
		return nil, errors.New("participant service not found")
	}

	return svc, nil
}
