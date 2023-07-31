package web

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	tests := map[string]struct {
		err  error
		code int
		msg  Message
		log  Log
	}{
		"log and message fields": {
			err:  errors.New("test error 1"),
			code: 500,
		},
		"message field": {
			err:  errors.New("test error 2"),
			code: 404,
			msg:  Message{"client": "Client error message 2", "detail": "Detailed client error message 2"},
		},
		"log field": {
			err:  errors.New("test error 3"),
			code: 400,
			log:  Log{"log": "Log error message 3", "debug": "Debug log message 3"},
		},
		"fields of both types": {
			err:  errors.New("test error 5"),
			code: 403,
			msg:  Message{"client": "Client error message 5", "detail": "Detailed client error message 5"},
			log:  Log{"log": "Log error message 5", "debug": "Debug log message 5"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err, ok := ErrorAs(NewError(tc.err, tc.code, tc.msg, tc.log))
			require.True(t, ok, "Expected Error type")

			require.True(t, ok, "Expected error, got nil")

			require.NotNil(t, err, "Expected error, got nil")
			require.Equal(t, tc.err.Error(), err.Error(), "Error messages do not match")
			require.Equal(t, tc.code, err.StatusCode(), "Status codes do not match")

			require.EqualValues(t, tc.msg, err.Message, "Client messages do not match")
			require.EqualValues(t, tc.log, err.Log, "Log messages do not match")
		})
	}
}
