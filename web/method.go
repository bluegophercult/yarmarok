package web

import (
	"net/http"
	"path"
)

type (
	M0[Resp any]      func() (Resp, error)
	M1[Req, Resp any] func(Req) (Resp, error)
	M2[Req, Resp any] func(string, Req) (Resp, error)
)

func newM0[Resp any](m M0[Resp]) M0[Resp]                { return m }
func newM1[Req, Resp any](m M1[Req, Resp]) M1[Req, Resp] { return m }
func newM2[Req, Resp any](m M2[Req, Resp]) M2[Req, Resp] { return m }

func (m M0[Resp]) Handle(rw http.ResponseWriter, _ *http.Request) error {
	resp, err := m()
	if err != nil {
		return err
	}

	return EncodeBody(rw, &resp)
}

func (m M1[Req, Resp]) Handle(rw http.ResponseWriter, req *http.Request) error {
	var request Req
	if err := DecodeBody(req.Body, &request); err != nil {
		return err
	}

	resp, err := m(request)
	if err != nil {
		return err
	}

	return EncodeBody(rw, &resp)
}

func (m M2[Req, Resp]) Handle(rw http.ResponseWriter, req *http.Request) error {
	id := path.Base(req.URL.String())

	var request Req
	if err := DecodeBody(req.Body, &request); err != nil {
		return err
	}

	resp, err := m(id, request)
	if err != nil {
		return err
	}

	return EncodeBody(rw, &resp)
}
