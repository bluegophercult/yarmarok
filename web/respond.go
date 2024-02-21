package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// respond writes minimalistic response.
// function signature and error/status handling may be different.
func (r *Router) respond(rw http.ResponseWriter, data any) {
	if data == nil {
		return
	}

	var buf bytes.Buffer
	if err := r.encodeBody(&buf, data); err != nil {
		err = fmt.Errorf("encoding to buffer: %w", err)
		r.respondErr(rw, err)
		return
	}

	if _, err := buf.WriteTo(rw); err != nil {
		err = fmt.Errorf("writing response: %w", err)
		r.respondErr(rw, err)
		return
	}
}

func (r *Router) respondErr(rw http.ResponseWriter, err error) {
	r.logger.WithError(err).Warn("responding with error")
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}

// decodeBody reads data from a body and converts it to any.
func (r *Router) decodeBody(body io.Reader, data any) error {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	return nil
}

// encodeBody writes data to a writer after converting it to JSON.
func (r *Router) encodeBody(rw io.Writer, data any) error {
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	return nil
}
