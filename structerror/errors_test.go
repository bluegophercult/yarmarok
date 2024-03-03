package structerror

import (
	"errors"
	"fmt"
	"maps"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	err := New("code", "test error")
	assert.Error(t, err)
}

func TestCoder(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		err := New("code", "test error")

		var coder Coder = Coder(nil)
		require.True(t, errors.As(err, &coder))
		assert.Equal(t, "code", coder.Code())
	})

	t.Run("No code", func(t *testing.T) {
		err := New("", "test error")

		var coder Coder = Coder(nil)
		require.True(t, errors.As(err, &coder))
		assert.Equal(t, "", coder.Code())
	})
}

func TestError(t *testing.T) {
	t.Run("normal format", func(t *testing.T) {
		err := New("code", "test error")
		assert.Equal(t, "code: test error", err.Error())
	})

	t.Run("empty code", func(t *testing.T) {
		err := New("", "test error")
		assert.Equal(t, "test error", err.Error())
	})

	t.Run("empty error", func(t *testing.T) {
		err := New("code", "")
		assert.Equal(t, "code", err.Error())
	})

	t.Run("empty code and error", func(t *testing.T) {
		err := New("", "")
		assert.Equal(t, "", err.Error())
	})

	t.Run("fmt", func(t *testing.T) {
		err := New("code", "test error: %s", "param")
		assert.Equal(t, "code: test error: param", err.Error())
	})
}

func TestWrap(t *testing.T) {
	base := &net.AddrError{Err: "test error", Addr: "test addr"}
	structured := New("code", "test error: %w", base)
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

	t.Run("with value", func(t *testing.T) {
		err := New("code", "test error")
		withValue := WithValue(err, "key1", "value1")
		withValue = WithValue(withValue, "key2", "value2")

		assert.Equal(t, "key2=value2: key1=value1: code: test error", withValue.Error())
		assert.ErrorIs(t, withValue, err)
	})

}

type Coder interface {
	Code() string
}

type CodeError struct {
	error
	code string
}

func New(code, format string, args ...any) error {
	return &CodeError{code: code, error: fmt.Errorf(format, args...)}
}

func (e *CodeError) Error() string {
	if e.code == "" {
		return e.error.Error()
	}

	if e.error.Error() == "" {
		return e.code
	}

	return e.code + ": " + e.error.Error()
}

func (e *CodeError) Unwrap() error {
	return e.error
}

func (e *CodeError) Code() string {
	return e.code
}

func WithCode(code string, err error) error {
	if err == nil {
		return nil
	}
	return &CodeError{code: code, error: err}
}

type ValueError struct {
	error
	values map[string]string
}

func WithValue(err error, key, value string) *ValueError {
	if err == nil {
		return nil
	}

	return &ValueError{error: err, values: map[string]string{key: value}}
}

func (e *ValueError) Error() string {
	if len(e.values) == 0 {
		return e.error.Error()
	}

	return fmt.Sprintf("%s: %s", e.formatValues(), e.error.Error())
}

func (e *ValueError) Unwrap() error {
	return e.error
}

func (e *ValueError) formatValues() string {
	strs := make([]string, 0, len(e.values))
	for k, v := range e.values {
		strs = append(strs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(strs, ", ")
}

func (e *ValueError) Values() map[string]string {
	values := maps.Clone(e.values)
	var valueErr *ValueError

	if !errors.As(e.error, &valueErr) {
		return values
	}

	maps.Copy(values, valueErr.Values())

	return values
}
