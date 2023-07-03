package resperr

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert := assert.New(t)
	t.Run("error", func(t *testing.T) {
		testError := NewError("code", "message")

		assert.Equal("message", testError.Error())
		assert.Equal("code", testError.Code())
	})

	t.Run("json", func(t *testing.T) {
		testError := NewError("code", "message")

		b, err := json.Marshal(testError)
		assert.NoError(err)

		assert.Equal(`{"code":"code","message":"message"}`, string(b))
	})

	t.Run("is_error", func(t *testing.T) {
		testError := NewError("code", "message")
		secondTestError := NewError("code", "message")

		t.Run("self", func(t *testing.T) {
			assert.ErrorIs(testError, testError)
		})

		t.Run("same", func(t *testing.T) {
			assert.NotErrorIs(testError, secondTestError)
		})

		t.Run("same_reversed", func(t *testing.T) {
			assert.NotErrorIs(secondTestError, testError)
		})

		t.Run("nil", func(t *testing.T) {
			assert.NotErrorIs(nil, testError)
		})

		t.Run("nil_reversed", func(t *testing.T) {
			assert.NotErrorIs(testError, nil)
		})

		t.Run("different", func(t *testing.T) {
			assert.NotErrorIs(testError, NewError("different", "message"))
		})

		t.Run("different_reversed", func(t *testing.T) {
			assert.NotErrorIs(NewError("different", "message"), testError)
		})

		t.Run("wrapped", func(t *testing.T) {
			assert.ErrorIs(fmt.Errorf("wrapped error: %w", testError), testError)
		})
	})

	t.Run("as_error", func(t *testing.T) {
		testError := NewError("code", "message")

		t.Run("self", func(t *testing.T) {
			result := error(&Error{})
			assert.ErrorAs(testError, &result)
			assert.ErrorIs(testError, result)
		})

	})

}
