package web

import (
	"net/http"
	"path"
)

type (
	// In represents the request type.
	In any
	// Out represents the response type.
	Out any
	// ID represents the identifier type.
	ID = string
)

type ( // Define generic types for List, Create, Delete, and Update operations.
	List[O Out]         func() (O, error)
	Create[I In, O Out] func(I) (O, error)
	Delete              func(ID) error
	Update[I In]        func(ID, I) error
)

// newList returns a new List function.
func newList[O Out](m List[O]) List[O] { return m }

// newCreate returns a new Create function.
func newCreate[I In, O Out](m Create[I, O]) Create[I, O] { return m }

// newUpdate returns a new Update function.
func newUpdate[I In](m Update[I]) Update[I] { return m }

// newDelete returns a new Delete function.
func newDelete(m Delete) Delete { return m }

// Handle handles the HTTP request for the List operation.
func (m List[O]) Handle(rw http.ResponseWriter, _ *http.Request) error {
	resp, err := m()
	if err != nil {
		return err
	}

	return Respond(rw, resp)
}

// Handle handles the HTTP request for the Create operation.
func (m Create[I, O]) Handle(rw http.ResponseWriter, req *http.Request) error {
	var in I
	if err := DecodeBody(req.Body, &in); err != nil {
		return err
	}

	out, err := m(in)
	if err != nil {
		return err
	}

	return Respond(rw, out)
}

// Handle handles the HTTP request for the Update operation.
func (m Update[I]) Handle(_ http.ResponseWriter, req *http.Request) error {
	var in I
	if err := DecodeBody(req.Body, &in); err != nil {
		return err
	}

	id := path.Base(req.URL.String())

	return m(id, in)
}

// Handle handles the HTTP request for the Delete operation.
func (m Delete) Handle(_ http.ResponseWriter, req *http.Request) error {
	id := path.Base(req.URL.String())

	return m(id)
}
