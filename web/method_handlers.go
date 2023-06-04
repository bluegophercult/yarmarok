package web

import (
	"encoding/json"
	"net/http"

	"github.com/kaznasho/yarmarok/logger"
)

// noRequestMethodHandler is a wrapper around a service method
// that converts a no request method to an http handler.
type noRequestMethodHandler[Response any] struct {
	method func() (Response, error)
	logger *logger.Logger
}

func newNoRequestMethodHandler[Response any](
	method func() (Response, error),
	log *logger.Logger,
) noRequestMethodHandler[Response] {
	return noRequestMethodHandler[Response]{
		method: method,
		logger: log,
	}
}

func (m noRequestMethodHandler[Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := m.method()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		m.logger.WithError(err).Error("request execution failed")
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		m.logger.WithError(err).Error("encoding response failed")
		return
	}

	m.logger.WithFields(
		logger.Fields{
			"response": resp,
		},
	).Debug("response")

	w.WriteHeader(http.StatusOK)
}

// methodHandler is a wrapper around a service method
// that converts that method to an http handler.
type methodHandler[Request any, Response any] struct {
	method func(Request) (Response, error)
	logger *logger.Logger
}

func newMethodHandler[Request any, Response any](
	method func(Request) (Response, error),
	log *logger.Logger,
) methodHandler[Request, Response] {
	return methodHandler[Request, Response]{
		method: method,
		logger: log,
	}
}

func (m methodHandler[Request, Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := new(Request)

	err := json.NewDecoder(req.Body).Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		m.logger.WithError(err).Error("decoding request body failed")
		return
	}

	m.logger.WithFields(
		logger.Fields{
			"request": request,
		},
	).Debug("request")

	resp, err := m.method(*request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		m.logger.WithError(err).Error("request execution failed")
		return
	}

	m.logger.WithFields(
		logger.Fields{
			"response": resp,
		},
	).Debug("response")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		m.logger.WithError(err).Error("encoding response failed")
		return
	}

	w.WriteHeader(http.StatusOK)
}
