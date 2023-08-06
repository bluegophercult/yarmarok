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

// ServiceFunc is responsible for fetching a service from request data, where T - service type/
type ServiceFunc[T, I, O any] func(*http.Request) (T, error)

func newServiceFunc[T, I, O any](fn ServiceFunc[T, I, O]) ServiceFunc[T, I, O] { return fn }

func (sf ServiceFunc[T, I, _]) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := extractService[T, Creator[I]](sf, req)
		if err != nil {
			respondErr(rw, err)
			return
		}
		newHandler[Create[I]](svc.Create).Handle(rw, req)
	}
}

func (sf ServiceFunc[T, _, O]) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := extractService[T, Getter[O]](sf, req)
		if err != nil {
			respondErr(rw, err)
			return
		}
		newHandler[Get[O]](svc.Get).Handle(rw, req)
	}
}

func (sf ServiceFunc[T, I, _]) Edit() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := extractService[T, Editor[I]](sf, req)
		if err != nil {
			respondErr(rw, err)
			return
		}
		newHandler[Edit[I]](svc.Edit).Handle(rw, req)
	}
}

func (sf ServiceFunc[T, _, _]) Delete() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := extractService[T, Deleter](sf, req)
		if err != nil {
			respondErr(rw, err)
			return
		}
		newHandler[Delete](svc.Delete).Handle(rw, req)
	}
}

func (sf ServiceFunc[T, _, O]) List() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := extractService[T, Lister[O]](sf, req)
		if err != nil {
			respondErr(rw, err)
			return
		}

		newHandler[List[O]](svc.List).Handle(rw, req)
	}
}

// extractService fetch a service T from service func using request data
// checking if it implements a specific CRUD interface S.
func extractService[T, S any](fn func(*http.Request) (T, error), req *http.Request) (S, error) {
	var svc S
	t, err := fn(req)
	if err != nil {
		return svc, err
	}

	svc, ok := any(t).(S)
	if !ok {
		return svc, ErrNotImplemented
	}

	return svc, nil
}
