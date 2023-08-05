/*
DISCLAIMER: Optional feature designed for mapping busines logic services to HTTP handlers.
*/
package web

import (
	"errors"
	"net/http"
)

var ErrNotImplemented = errors.New("not implemented")

type Creator[T any] interface {
	Create(T) (ID, error)
}

type Getter[T any] interface {
	Get(ID) (T, error)
}

type Editor[T any] interface {
	Edit(ID, T) error
}

type Deleter interface {
	Delete(ID) error
}

type Lister[T any] interface {
	List() ([]T, error)
}

// ServiceFunc is responsible for fetching a service from request data.
type ServiceFunc[T any] func(*http.Request) (T, error)

func newServiceFunc[T any](fn ServiceFunc[T]) ServiceFunc[T] { return fn }

func (sf ServiceFunc[T]) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		val, err := sf(req)
		if err != nil {
			respond(rw, err)
		}

		svc, ok := any(val).(Creator[T])
		if !ok {
			respond(rw, ErrNotImplemented)
		}

		newMethod[Create[T]](svc.Create).Handle(rw, req)
	}
}

func (sf ServiceFunc[T]) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		val, err := sf(req)
		if err != nil {
			respond(rw, err)
		}

		svc, ok := any(val).(Creator[T])
		if !ok {
			respond(rw, ErrNotImplemented)
		}

		newMethod[Create[T]](svc.Create).Handle(rw, req)
	}
}

func (sf ServiceFunc[T]) Edit() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		val, err := sf(req)
		if err != nil {
			respond(rw, err)
		}

		svc, ok := any(val).(Editor[T])
		if !ok {
			respond(rw, ErrNotImplemented)
		}

		newMethod[Update[T]](svc.Edit).Handle(rw, req)
	}
}

func (sf ServiceFunc[T]) Delete() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		val, err := sf(req)
		if err != nil {
			respond(rw, err)
		}

		svc, ok := any(val).(Deleter)
		if !ok {
			respond(rw, ErrNotImplemented)
		}

		newMethod[Delete](svc.Delete).Handle(rw, req)
	}
}

func (sf ServiceFunc[T]) List() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		val, err := sf(req)
		if err != nil {
			respond(rw, err)
		}

		svc, ok := any(val).(Lister[T])
		if !ok {
			respond(rw, ErrNotImplemented)
		}

		newMethod[List[T]](svc.List).Handle(rw, req)
	}
}
