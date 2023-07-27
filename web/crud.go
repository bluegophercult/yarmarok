package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

// MethodCreate is a wrapper around methods with create semantics.
type MethodCreate[Request any, Response any] func(Request) (Response, error)

// NewMethodCreate creates a new MethodCreate.
func NewMethodCreate[Request any, Response any](m MethodCreate[Request, Response]) MethodCreate[Request, Response] {
	return m
}

// ServeHTTP implements http.Handler interface.
// It decodes request, calls the wrapped method,
// and responds with the result.
func (m MethodCreate[Request, Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := new(Request)

	err := json.NewDecoder(req.Body).Decode(request)
	if err != nil {
		respondError(w, err)
		return
	}

	callAndRespond(w, func() (Response, error) {
		return m(*request)
	})
}

// MethodUpdate is a wrapper around methods with update semantics.
type MethodUpdate[Request any, Response any] func(string, Request) (Response, error)

// NewMethodUpdate creates a new MethodUpdate.
func NewMethodUpdate[Request any, Response any](m MethodUpdate[Request, Response]) MethodUpdate[Request, Response] {
	return m
}

// ServeHTTP implements http.Handler interface.
// It decodes request, extracts the id of entity from path,
// calls the wrapped method, and responds with the result.
func (m MethodUpdate[Request, Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := new(Request)
	id := extractID(req)

	err := json.NewDecoder(req.Body).Decode(request)
	if err != nil {
		respondError(w, err)
		return
	}

	callAndRespond(w, func() (Response, error) {
		return m(id, *request)
	})
}

// MethodDelete is a wrapper around methods with delete semantics.
type MethodDelete[Response any] func(string) error

// NewMethodDelete creates a new MethodDelete.
func NewMethodDelete[Response any](m MethodDelete[Response]) MethodDelete[Response] {
	return m
}

// ServeHTTP implements http.Handler interface.
// It extracts the id of entity from path,
// calls the wrapped method, and responds with the result.
func (m MethodDelete[Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := extractID(req)

	callAndRespond(w, func() (Response, error) {
		return *new(Response), m(id)
	})
}

// MethodList is a wrapper around methods with list semantics.
type MethodList[Response any] func() (Response, error)

// NewMethodList creates a new MethodList.
func NewMethodList[Response any](m MethodList[Response]) MethodList[Response] {
	return m
}

// ServeHTTP implements http.Handler interface.
// It calls the wrapped method, and responds with the result.
func (m MethodList[Response]) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	callAndRespond(w, m)
}

// MethodGet is a wrapper around methods with get semantics.
type MethodGet[Response any] func(string) (Response, error)

// NewMethodGet creates a new MethodGet.
func NewMethodGet[Response any](m MethodGet[Response]) MethodGet[Response] {
	return m
}

// ServeHTTP implements http.Handler interface.
// It extracts the id of entity from path,
// calls the wrapped method, and responds with the result.
func (m MethodGet[Response]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := extractID(req)

	callAndRespond(w, func() (Response, error) {
		return m(id)
	})
}

func callAndRespond[Response any](w http.ResponseWriter, f func() (Response, error)) {
	resp, err := f()
	if err != nil {
		respondError(w, err)
		return
	}

	respondResult(w, resp)
}

func respondError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func respondResult[Response any](w http.ResponseWriter, resp Response) {
	json.NewEncoder(w).Encode(resp)
}

func extractID(req *http.Request) string {
	id := chi.URLParam(req, "id")
	return id
}
