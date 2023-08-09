package web

import (
	"errors"
	"net/http"
)

var ErrNotImplemented = errors.New("not implemented")

// CRUD interfaces mirror the CRUD functions,
// that must be implemented in the business logic layer.
// I - new/updated entity; O - output entity
type (
	Creator[I any] interface{ Create(I) (string, error) }
	Getter[O any]  interface{ Get(string) (O, error) }
	Editor[I any]  interface{ Edit(string, I) error }
	Deleter        interface{ Delete(string) error }
	Lister[O any]  interface{ List() ([]O, error) }
)

// Service wrapper around service function,
// which fetches from request services implemented in the business logic layer.
// It checks if underlying business layer service implements a collection of CRUD interfaces
// and maps them to corresponding functions that handle requests.
type Service[T, I, O any] ServiceFunc[T]

// ServiceFunc is a function that fetches from request services implemented in the business logic layer.
type ServiceFunc[T any] func(*http.Request) (T, error)

func newService[T, I, O any](fn Service[T, I, O]) Service[T, I, O] { return fn }

// Func is a method that casts service functions to ServiceFunc without I/O parameters.
func (s Service[T, _, _]) Func() ServiceFunc[T] { return ServiceFunc[T](s) }

func (s Service[T, I, _]) Create(rw http.ResponseWriter, req *http.Request) {
	svc, err := extractService[T, Creator[I]](s.Func(), req)
	if err != nil {
		respondErr(rw, err)
		return
	}
	newHandler[Create[I]](svc.Create).Handle(rw, req)
}

func (s Service[T, _, O]) Get(rw http.ResponseWriter, req *http.Request) {
	svc, err := extractService[T, Getter[O]](s.Func(), req)
	if err != nil {
		respondErr(rw, err)
		return
	}
	newHandler[Get[O]](svc.Get).Handle(rw, req)
}

func (s Service[T, I, _]) Edit(rw http.ResponseWriter, req *http.Request) {
	svc, err := extractService[T, Editor[I]](s.Func(), req)
	if err != nil {
		respondErr(rw, err)
		return
	}
	newHandler[Edit[I]](svc.Edit).Handle(rw, req)
}

func (s Service[T, _, _]) Delete(rw http.ResponseWriter, req *http.Request) {
	svc, err := extractService[T, Deleter](s.Func(), req)
	if err != nil {
		respondErr(rw, err)
		return
	}
	newHandler[Delete](svc.Delete).Handle(rw, req)
}

func (s Service[T, _, O]) List(rw http.ResponseWriter, req *http.Request) {
	svc, err := extractService[T, Lister[O]](s.Func(), req)
	if err != nil {
		respondErr(rw, err)
		return
	}

	newHandler[List[O]](svc.List).Handle(rw, req)
}

// extractService fetch a service T from service func using request data
func extractService[T, S any](fn ServiceFunc[T], req *http.Request) (S, error) {
	var svc S

	val, err := fn(req)
	if err != nil {
		return svc, err
	}

	return validateService[S](val)
}

// validateService checks if service implements a specific CRUD interface S.
func validateService[S any](val any) (S, error) {
	svc, ok := val.(S)
	if !ok {
		return svc, ErrNotImplemented
	}

	return svc, nil
}
