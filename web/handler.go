package web

import (
	"net/http"
	"path"
)

// CRUD functions accept parameter T,
// which can be a new, updated or created entity.
type (
	Create[T any] func(T) (id string, err error)
	Get[T any]    func(id string) (T, error)
	Update[T any] func(id string, upd T) error
	Delete        func(id string) error
	List[T any]   func() ([]T, error)
)

// newMethod creates a new instance of CRUD-Func.
func newMethod[F any](fn F) F { return fn }

func (m Create[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := decodeBody(req.Body, &in); err != nil {
		respond(rw, err)
		return
	}

	id, err := m(in)
	if err != nil {
		respond(rw, err)
		return
	}

	respond(rw, struct{ ID string }{id})
}

func (m Get[_]) Handle(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.String())

	out, err := m(id)
	if err != nil {
		respond(rw, err)
		return
	}

	respond(rw, out)
}

// Handle handles the HTTP request for the Update operation.
func (m Update[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := decodeBody(req.Body, &in); err != nil {
		respond(rw, err)
		return
	}

	id := path.Base(req.URL.String())
	if err := m(id, in); err != nil {
		respond(rw, err)
		return
	}
}

// Handle handles the HTTP request for the Delete operation.
func (m Delete) Handle(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.String())

	if err := m(id); err != nil {
		respond(rw, err)
	}
}

// Handle handles the HTTP request for the List operation.
func (m List[O]) Handle(rw http.ResponseWriter, _ *http.Request) {
	out, err := m()
	if err != nil {
		respond(rw, err)
		return
	}

	respond(rw, struct{ Data []O }{out})
}
