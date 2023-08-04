package web

import (
	"net/http"
	"path"
)

type ( // Input/Output types.
	In  any
	Out any
	ID  = string
)

type ( // CRUD functions.
	Create[I In] func(I) (ID, error)
	Get[O Out]   func(ID) (O, error)
	Update[I In] func(ID, I) error
	Delete       func(ID) error
	List[O Out]  func() ([]O, error)
)

// newMethod creates a new instance of CRUD-Func.
func newMethod[F any](fn F) F { return fn }

func (m Create[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := decodeBody(req.Body, &in); err != nil {
		respond(rw, err)
	}

	id, err := m(in)
	if err != nil {
		respond(rw, err)
	}

	respond(rw, struct{ ID ID }{id})
}

func (m Get[_]) Handle(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.String())

	out, err := m(id)
	if err != nil {
		respond(rw, err)
	}

	respond(rw, out)
}

// Handle handles the HTTP request for the Update operation.
func (m Update[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := decodeBody(req.Body, &in); err != nil {
		respond(rw, err)
	}

	id := path.Base(req.URL.String())

	if err := m(id, in); err != nil {
		respond(rw, err)
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
	}

	respond(rw, struct{ Data []O }{out})
}
