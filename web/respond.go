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
		return
	}

	rw.WriteHeader(code)

	if _, err := buf.WriteTo(rw); err != nil {
		err = fmt.Errorf("writing response: %w", err)
		respond(rw, err)
		return
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
