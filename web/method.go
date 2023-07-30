package web

import (
	"net/http"
	"path"
)

type (
	List[Resp any]        func() (Resp, error)
	Create[Req, Resp any] func(Req) (Resp, error)
	Delete                func(string) error
	Update[Req any]       func(string, Req) error
)

func newList[Resp any](m List[Resp]) List[Resp]                      { return m }
func newCreate[Req, Resp any](m Create[Req, Resp]) Create[Req, Resp] { return m }
func newUpdate[Req any](m Update[Req]) Update[Req]                   { return m }
func newDelete(m Delete) Delete                                      { return m }

func (m List[Resp]) Handle(rw http.ResponseWriter, _ *http.Request) error {
	resp, err := m()
	if err != nil {
		return err
	}

	return EncodeBody(rw, resp)
}

func (m Create[Req, Resp]) Handle(rw http.ResponseWriter, req *http.Request) error {
	var request Req
	if err := DecodeBody(req.Body, &request); err != nil {
		return err
	}

	resp, err := m(request)
	if err != nil {
		return err
	}

	return EncodeBody(rw, resp)
}

func (m Update[Req]) Handle(_ http.ResponseWriter, req *http.Request) error {
	var request Req
	if err := DecodeBody(req.Body, &request); err != nil {
		return err
	}

	id := path.Base(req.URL.String())

	return m(id, request)
}

func (m Delete) Handle(_ http.ResponseWriter, req *http.Request) error {
	id := path.Base(req.URL.String())

	return m(id)
}
