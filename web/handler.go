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

// CreateHandler is a wrapper around a service method
// that creates a new object.
type CreateHandler[I any] struct {
	Create[I]
	*Router
}

// NewCreateHandler creates a new CreateHandler.
func NewCreateHandler[I any](router *Router, fn Create[I]) CreateHandler[I] {
	return CreateHandler[I]{
		Create: fn,
		Router: router,
	}
}

// Handle handles a create request.
func (h CreateHandler[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := h.decodeBody(req.Body, &in); err != nil {
		h.respondErr(rw, err)
		return
	}

	id, err := h.Create(in)
	if err != nil {
		h.respondErr(rw, err)
		return
	}

	h.respond(rw, CreateResponse{id})
}

// CreateResponse represents the response structure containing an item ID.
type CreateResponse struct {
	ID string `json:"id"`
}

// GetHandler is a wrapper around a service method
// that returns an object by ID.
type GetHandler[O any] struct {
	Get[O]
	*Router
}

// NewGetHandler creates a new GetHandler.
func NewGetHandler[O any](router *Router, fn Get[O]) GetHandler[O] {
	return GetHandler[O]{
		Get:    fn,
		Router: router,
	}
}

// Handle handles a get request.
func (h GetHandler[O]) Handle(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.String())

	out, err := h.Get(id)
	if err != nil {
		h.respondErr(rw, err)
		return
	}

	h.respond(rw, out)
}

// EditHandler is a wrapper around a service method
// that updates an object.
type EditHandler[I any] struct {
	Edit[I]
	*Router
}

// NewEditHandler creates a new EditHandler.
func NewEditHandler[I any](router *Router, fn Edit[I]) EditHandler[I] {
	return EditHandler[I]{
		Edit:   fn,
		Router: router,
	}
}

// Handle handles an edit request.
func (h EditHandler[I]) Handle(rw http.ResponseWriter, req *http.Request) {
	var in I
	if err := h.decodeBody(req.Body, &in); err != nil {
		h.respondErr(rw, err)
		return
	}

	id := path.Base(req.URL.String())
	if err := h.Edit(id, in); err != nil {
		h.respondErr(rw, err)
		return
	}
}

// DeleteHandler is a wrapper around a service method
// that deletes an object by ID.
type DeleteHandler struct {
	Delete
	*Router
}

// NewDeleteHandler creates a new DeleteHandler.
func NewDeleteHandler(router *Router, fn Delete) DeleteHandler {
	return DeleteHandler{
		Delete: fn,
		Router: router,
	}
}

// Handle handles a delete request.
func (h DeleteHandler) Handle(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.String())

	if err := h.Delete(id); err != nil {
		h.respondErr(rw, err)
	}
}

// ListHandler is a wrapper around a service method
// that returns a all objects of kind.
type ListHandler[O any] struct {
	List[O]
	*Router
}

// NewListHandler creates a new ListHandler.
func NewListHandler[O any](router *Router, fn List[O]) ListHandler[O] {
	return ListHandler[O]{
		List:   fn,
		Router: router,
	}
}

// Handle handles a list request.
func (h ListHandler[O]) Handle(rw http.ResponseWriter, _ *http.Request) {
	out, err := h.List()
	if err != nil {
		h.respondErr(rw, err)
		return
	}

	h.respond(rw, ListResponse[O]{out})
}

// ListResponse represents a generic response containing an array of items.
type ListResponse[O any] struct {
	Items []O `json:"items"`
}
