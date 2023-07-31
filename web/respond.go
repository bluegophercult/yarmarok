package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Respond responds with converted data to the client with the given status code.
func Respond(rw http.ResponseWriter, data any) error {
	val, ok := data.(interface{ StatusCode() int })
	if !ok || data == nil {
		rw.WriteHeader(http.StatusNoContent)
		return nil
	}

	var buf bytes.Buffer
	if err := EncodeBody(&buf, data); err != nil {
		return fmt.Errorf("encoding to buffer: %w", err)
	}

	rw.WriteHeader(val.StatusCode())

	if _, err := buf.WriteTo(rw); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}

// DecodeBody reads data from a body and converts it to any.
func DecodeBody(body io.Reader, data any) error {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	return nil
}

// EncodeBody writes data to a writer after converting it to JSON.
func EncodeBody(rw io.Writer, data any) error {
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	return nil
}
