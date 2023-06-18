package web

import (
	"net/http"
	"time"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/rs/cors"
)

const (
	// GoogleUserIDHeader is the header that contains the user id
	// set by google identity aware proxy.
	GoogleUserIDHeader = "X-Goog-Authenticated-User-Id"

	defaultOrigin = "https://yarmarock.com.ua"
)

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

func (r *Router) corsMiddleware(next http.Handler) http.Handler {
	return cors.New(
		cors.Options{
			AllowedOrigins: []string{defaultOrigin},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},

			AllowedHeaders: []string{
				"Accept",
				"Authorization",
				"Content-Type",
				"X-CSRF-Token",
				"X-Goog-Authenticated-User-Id",
			},
			ExposedHeaders:       []string{},
			MaxAge:               0,
			AllowPrivateNetwork:  false,
			OptionsPassthrough:   false,
			OptionsSuccessStatus: 0,
			Debug:                false,
		},
	).Handler(next)
}
