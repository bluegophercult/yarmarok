package structerror

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	t.Run("normal format", func(t *testing.T) {
		err := New("code")
		assert.Equal(t, "code", err.Error())
	})

	t.Run("empty code", func(t *testing.T) {
		err := New("")
		assert.Equal(t, "", err.Error())
	})
}

func TestWrap(t *testing.T) {
	structured := New("code")
	wrapper := fmt.Errorf("wrapped error: %w", structured)

	t.Run("as", func(t *testing.T) {
		t.Run("type", func(t *testing.T) {
			var wrappedErr *Error
			require.True(t, errors.As(wrapper, &wrappedErr))
			assert.Equal(t, "code", wrappedErr.Code())
		})

		t.Run("interface", func(t *testing.T) {
			var wrappedErr Coder
			require.False(t, errors.As(fmt.Errorf("test error"), &wrappedErr))
			assert.Nil(t, wrappedErr)
		})
	})

	t.Run("is", func(t *testing.T) {
		assert.True(t, errors.Is(wrapper, structured))
	})

	t.Run("unwrap", func(t *testing.T) {
		assert.Equal(t, structured, errors.Unwrap(wrapper))
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
			var wrappedErr *Error
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
		err := New("code")
		withValue := WithLabel(err, "key1", "value1")
		withValue = WithLabel(withValue, "key2", "value2")

		assert.Equal(t, "key2=value2: key1=value1: code", withValue.Error())
		assert.ErrorIs(t, withValue, err)
	})

	t.Run("with labels", func(t *testing.T) {
		err := New("code")
		withValue := WithLabels(err, KV("key1", "value1"), KV("key2", "value2"))
		assert.Equal(t, "key1=value1, key2=value2: code", withValue.Error())

		t.Run("nil", func(t *testing.T) {
			withValue := WithLabels(nil, KV("key1", "value1"))
			assert.Nil(t, withValue)
		})
	})
}

func TestAPIUsecases(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		t.Run("Code error", func(t *testing.T) {
			err := New("code")
			expected := `{"errorCode":"code"}`

			data, jErr := AsJSON(err)
			require.NoError(t, jErr)

			assert.JSONEq(t, expected, string(data))
		})

		t.Run("Default error", func(t *testing.T) {
			err := fmt.Errorf("test error")
			expected := fmt.Sprintf(`{"errorCode":"%s"}`, DefaultCode)

			data, jErr := AsJSON(err)
			require.NoError(t, jErr)

			assert.JSONEq(t, expected, string(data))
		})
	})

	t.Run("Wrapped JSON", func(t *testing.T) {
		base := New("code")
		wrapped := fmt.Errorf("wrapped error: %w", base)

		expected := `{"errorCode":"code"}`

		data, jErr := AsJSON(wrapped)
		require.NoError(t, jErr)

		assert.JSONEq(t, expected, string(data))
	})

	t.Run("Response", func(t *testing.T) {
		t.Run("Code error", func(t *testing.T) {
			err := New("code")
			expected := Response{
				Code: "code",
			}

			response := AsResponse(err)
			assert.Equal(t, expected, response)
		})

		t.Run("Wrapped code error", func(t *testing.T) {
			base := New("code")
			wrapped := fmt.Errorf("wrapped error: %w", base)

			expected := Response{
				Code: "code",
			}

			response := AsResponse(wrapped)
			assert.Equal(t, expected, response)
		})

		t.Run("Default error", func(t *testing.T) {
			err := fmt.Errorf("test error")
			expected := Response{
				Code: DefaultCode,
			}

			response := AsResponse(err)
			assert.Equal(t, expected, response)
		})

		t.Run("Wrapped default error", func(t *testing.T) {
			base := fmt.Errorf("test error")
			wrapped := fmt.Errorf("wrapped error: %w", base)

			expected := Response{
				Code: DefaultCode,
			}

			response := AsResponse(wrapped)
			assert.Equal(t, expected, response)
		})
	})
}
