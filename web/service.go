/*
DISCLAIMER: Optional feature designed for mapping busines logic services to transport layer.
*/
package web

import "net/http"

// Service is an interface for CRUD operations.
type Service[I In, O Out] interface {
	Create(I) (ID, error)
	Get(ID) (O, error)
	Edit(ID, I) error
	Delete(ID) error
	List() ([]O, error)
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

func (m ServiceHandler[_, I, _]) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}

		newMethod[Create[I]](svc.Create).Handle(rw, req)
	}
}

func (m ServiceHandler[_, _, O]) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}
		newMethod[Get[O]](svc.Get).Handle(rw, req)
	}
}

func (m ServiceHandler[_, I, _]) Update() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}
		newMethod[Update[I]](svc.Edit).Handle(rw, req)
	}
}

func (m ServiceHandler[_, _, _]) Delete() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}

		newMethod[Delete](svc.Delete).Handle(rw, req)
	}
}

func (m ServiceHandler[_, _, O]) List() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		svc, err := m.ServiceFunc(req)
		if err != nil {
			respond(rw, err)
		}

		newMethod[List[O]](svc.List).Handle(rw, req)
	}
}
