package web

import (
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
