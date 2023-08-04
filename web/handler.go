package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// Service is an interface for CRUD operations.
type Service[I In, O Out] interface {
	Create(I) (ID, error)
	Edit(ID, I) error
	Delete(ID) error
	List() (O, error)
}

// ServiceFunc is responsible for fetching a service from request data.
type ServiceFunc[S any] func(*http.Request) (S, error)

// ServiceHandler is responsible for handling HTTP requests.
type ServiceHandler[S Service[I, O], I In, O Out] struct{ ServiceFunc[S] }

// newServiceHandler creates a new instance of ServiceHandler with the given ServiceFunc.
func newServiceHandler[S Service[I, O], I In, O Out](fn ServiceFunc[S]) ServiceHandler[S, I, O] {
	return ServiceHandler[S, I, O]{
		ServiceFunc: fn,
	}
}

// Create is responsible for creating a resource.
func (m ServiceHandler[_, I, _]) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}

		var in I
		if err := decodeBody(req.Body, &in); err != nil {
			respond(rw, err)
		}

		id, err := svc.Create(in)
		if err != nil {
			respond(rw, err)
		}

		resp := struct{ ID string }{ID: id}
		respond(rw, resp)
	}
}

// Edit is responsible for editing a resource.
func (m ServiceHandler[_, I, _]) Edit() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}

		var in I
		if err := decodeBody(req.Body, &in); err != nil {
			respond(rw, err)
		}

		id := path.Base(req.URL.String())

		if err := svc.Edit(id, in); err != nil {
			respond(rw, err)
		}
	}
}

// Delete is responsible for deleting a resource.
func (m ServiceHandler[_, _, _]) Delete() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}

		id := path.Base(req.URL.String())

		if err := svc.Delete(id); err != nil {
			respond(rw, err)
		}
	}
}

// List is responsible for listing resources.
func (m ServiceHandler[_, _, _]) List() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}

		result, err := svc.List()
		if err != nil {
			respond(rw, err)
		}

		respond(rw, result)
	}
}

// respond writes minimalistic response.
// function signature and error/status handling may be different.
func respond(rw http.ResponseWriter, data any) {
	if data == nil {
		return
	}

	code := http.StatusOK
	if _, ok := data.(error); ok {
		code = http.StatusInternalServerError
	}

	var buf bytes.Buffer
	if err := encodeBody(&buf, data); err != nil {
		err = fmt.Errorf("encoding to buffer: %w", err)
		respond(rw, err)
	}

	rw.WriteHeader(code)

	if _, err := buf.WriteTo(rw); err != nil {
		err = fmt.Errorf("writing response: %w", err)
		respond(rw, err)
	}
}

// decodeBody reads data from a body and converts it to any.
func decodeBody(body io.Reader, data any) error {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	return nil
}

// encodeBody writes data to a writer after converting it to JSON.
func encodeBody(rw io.Writer, data any) error {
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	return nil
}
