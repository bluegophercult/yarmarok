package web

import (
	"net/http"
	"path"
)

// CRUD functions accept I/O parameters,
// which correspond to new/updated or existing/created entity.
type (
	Create[I any] func(I) (id string, err error)
	Get[O any]    func(id string) (O, error)
	Edit[I any]   func(id string, upd I) error
	Delete        func(id string) error
	List[O any]   func() ([]O, error)
)

// CreateResponse represents the response structure containing an item ID.
type CreateResponse struct {
	ID string `json:"id"`
}

// ListResponse represents a generic response containing an array of items.
type ListResponse[O any] struct {
	Items []O `json:"items"`
}

// newHandler creates a new instance of CRUD-func.
func newHandler[T any](fn T) T { return fn }

func (m Create[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := decodeBody(req.Body, &in); err != nil {
		respondErr(rw, err)
		return
	}

	id, err := m(in)
	if err != nil {
		respondErr(rw, err)
		return
	}

	respond(rw, CreateResponse{id})
}

func (m Get[_]) Handle(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.String())

	out, err := m(id)
	if err != nil {
		respondErr(rw, err)
		return
	}

	respond(rw, out)
}

// Handle handles the HTTP request for the Edit operation.
func (m Edit[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := decodeBody(req.Body, &in); err != nil {
		respondErr(rw, err)
		return
	}

	id := path.Base(req.URL.String())
	if err := m(id, in); err != nil {
		respondErr(rw, err)
		return
	}
}

// Handle handles the HTTP request for the Delete operation.
func (m Delete) Handle(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.String())

	if err := m(id); err != nil {
		respondErr(rw, err)
	}
}

// Handle handles the HTTP request for the List operation.
func (m List[O]) Handle(rw http.ResponseWriter, _ *http.Request) {
	out, err := m()
	if err != nil {
		respondErr(rw, err)
		return
	}

	respond(rw, ListResponse[O]{out})
}
