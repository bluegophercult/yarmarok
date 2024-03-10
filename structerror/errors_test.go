package structerror

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	err := NewCodeError("code", "test error")
	assert.Error(t, err)
}

func TestError(t *testing.T) {
	t.Run("normal format", func(t *testing.T) {
		err := NewCodeError("code", "test error")
		assert.Equal(t, "code: test error", err.Error())
	})

	t.Run("empty code", func(t *testing.T) {
		err := NewCodeError("", "test error")
		assert.Equal(t, "test error", err.Error())
	})

	t.Run("empty error", func(t *testing.T) {
		err := NewCodeError("code", "")
		assert.Equal(t, "code", err.Error())
	})

	t.Run("empty code and error", func(t *testing.T) {
		err := NewCodeError("", "")
		assert.Equal(t, "", err.Error())
	})

	t.Run("fmt", func(t *testing.T) {
		err := NewCodeError("code", "test error: %s", "param")
		assert.Equal(t, "code: test error: param", err.Error())
	})
}

func TestWrap(t *testing.T) {
	base := &net.AddrError{Err: "test error", Addr: "test addr"}
	structured := NewCodeError("code", "test error: %w", base)
	wrapper := fmt.Errorf("wrapped error: %w", structured)

	t.Run("as", func(t *testing.T) {
		var wrappedErr *CodeError
		require.True(t, errors.As(wrapper, &wrappedErr))
		assert.Equal(t, "code", wrappedErr.Code())

		t.Run("nested", func(t *testing.T) {
			var addrErr *net.AddrError
			require.True(t, errors.As(wrapper, &addrErr))
			assert.Equal(t, base, addrErr)
		})
	})

	t.Run("is", func(t *testing.T) {
		t.Run("base", func(t *testing.T) {
			assert.True(t, errors.Is(structured, base))
		})

		t.Run("wrapper", func(t *testing.T) {
			assert.True(t, errors.Is(wrapper, structured))
		})
	})

	t.Run("unwrap", func(t *testing.T) {
		t.Run("base", func(t *testing.T) {
			fmtWrapped := errors.Unwrap(structured)
			assert.Equal(t, base, errors.Unwrap(fmtWrapped))
		})

		t.Run("wrapper", func(t *testing.T) {
			assert.Equal(t, structured, errors.Unwrap(wrapper))
		})
	})

	t.Run("wrap nil", func(t *testing.T) {
		err := WithCode("code", nil)
		assert.Nil(t, err)
	})

	t.Run("double wrap", func(t *testing.T) {
		e1 := WithCode("code1", fmt.Errorf("test error"))
		e2 := WithCode("code2", e1)
		assert.Equal(t, "code2: code1: test error", e2.Error())

		t.Run("as", func(t *testing.T) {
			var wrappedErr *CodeError
			require.True(t, errors.As(e2, &wrappedErr))
			assert.Equal(t, "code2", wrappedErr.Code())
		})

		t.Run("is", func(t *testing.T) {
			assert.True(t, errors.Is(e2, e1))
		})

		t.Run("join", func(t *testing.T) {
			e3 := errors.Join(e2, e1)
			assert.True(t, errors.Is(e3, e1))
			assert.True(t, errors.Is(e3, e2))
		})
	})
}

func TestDeveloperUsecases(t *testing.T) {
	t.Run("Wrap", func(t *testing.T) {
		err := WithCode("code", fmt.Errorf("test error"))
		assert.Equal(t, "code: test error", err.Error())
	})

	t.Run("with label", func(t *testing.T) {
		err := NewCodeError("code", "test error")
		withValue := WithLabel(err, "key1", "value1")
		withValue = WithLabel(withValue, "key2", "value2")

		assert.Equal(t, "key2=value2: key1=value1: code: test error", withValue.Error())
		assert.ErrorIs(t, withValue, err)
	})

	t.Run("with labels", func(t *testing.T) {
		err := NewCodeError("code", "test error")
		withValue := WithLabels(err, KV("key1", "value1"), KV("key2", "value2"))
		assert.Equal(t, "key1=value1, key2=value2: code: test error", withValue.Error())

		t.Run("nil", func(t *testing.T) {
			withValue := WithLabels(nil, KV("key1", "value1"))
			assert.Nil(t, withValue)
		})
	})
}

func TestAPIUsecases(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		err := NewCodeError("code", "test error")
		expected := `{"code":"code","error":"test error"}`

		data, jErr := AsJSON(err)
		require.NoError(t, jErr)

		assert.JSONEq(t, expected, string(data))
	})

	t.Run("Wrapped JSON", func(t *testing.T) {
		base := NewCodeError("code", "test error")
		wrapped := fmt.Errorf("wrapped error: %w", base)

		expected := `{"code":"code","error":"test error"}`

		data, jErr := AsJSON(wrapped)
		require.NoError(t, jErr)

		assert.JSONEq(t, expected, string(data))
	})

	t.Run("Response", func(t *testing.T) {
		t.Run("Code error", func(t *testing.T) {
			err := NewCodeError("code", "test error")
			expected := Response{
				Code:    "code",
				Message: DefaultMessage,
			}

			response := AsResponse(err)
			assert.Equal(t, expected, response)
		})

		t.Run("Wrapped code error", func(t *testing.T) {
			base := NewCodeError("code", "test error")
			wrapped := fmt.Errorf("wrapped error: %w", base)

			expected := Response{
				Code:    "code",
				Message: DefaultMessage,
			}

			response := AsResponse(wrapped)
			assert.Equal(t, expected, response)
		})

		t.Run("Default error", func(t *testing.T) {
			err := fmt.Errorf("test error")
			expected := Response{
				Code:    DefaultCode,
				Message: DefaultMessage,
			}

			response := AsResponse(err)
			assert.Equal(t, expected, response)
		})

		t.Run("Wrapped default error", func(t *testing.T) {
			base := fmt.Errorf("test error")
			wrapped := fmt.Errorf("wrapped error: %w", base)

			expected := Response{
				Code:    DefaultCode,
				Message: DefaultMessage,
			}

			response := AsResponse(wrapped)
			assert.Equal(t, expected, response)
		})
	})
}
