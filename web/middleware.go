package web

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/cors"

	"github.com/kaznasho/yarmarok/logger"
)

const (
	// GoogleUserIDHeader is the header that contains the organizer id
	// set by google identity aware proxy.
	GoogleUserIDHeader = "X-Goog-Authenticated-User-Id"

	defaultOrigin = "https://yarmarock.com.ua"
)

func (r *Router) organizerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		organizerID, err := extractOrganizerID(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			r.logger.WithError(err).Warn("failed to extract organizer id")
			return
		}

		err = r.organizerService.CreateOrganizerIfNotExists(organizerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			r.logger.WithError(err).Warn("failed to init organizer")
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		organizerID, _ := extractOrganizerID(req)

		start := time.Now()
		duration := time.Since(start)

		lrw := logger.NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, req)

		responseMetric := lrw.ResponseMetric()

		r.logger.WithFields(
			logger.Fields{
				"uri":          req.RequestURI,
				"method":       req.Method,
				"status":       responseMetric.Status,
				"duration":     duration,
				"size":         responseMetric.Size,
				"organizer_id": organizerID,
			},
		).Info("request completed")
	})
}

func (r *Router) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				r.logger.WithFields(logger.Fields{
					"uri":    req.RequestURI,
					"method": req.Method,
					"rec":    rec,
					"trace":  string(debug.Stack()),
				}).Error("panic recovered")

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, req)
	})
}

var allowedOrigins = []string{defaultOrigin}

func (r *Router) corsMiddleware(next http.Handler) http.Handler {
	return cors.New(
		cors.Options{
			AllowedOrigins: allowedOrigins,
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

func (r *Router) headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}
