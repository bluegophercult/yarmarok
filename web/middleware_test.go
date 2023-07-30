package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestWithErrors(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	log := logger.NewNoOpLogger()

	web := NewWeb(logger.NewNoOpLogger(), osMock)
	require.NotNil(t, web)

	t.Run("default error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := newHandlerStub()
		mockedError := errors.New("mocked error")
		h.WithError(mockedError)
		rw := httptest.NewRecorder()

		require.Nil(t, WithErrors(log)(h.ServeHTTP)(rw, req))

		require.Equal(t, http.StatusInternalServerError, rw.Code)
		require.True(t, h.Called())
	})

	t.Run("custom error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := newHandlerStub()
		mockedError := NewError(errors.New("mocked error"), http.StatusBadRequest)
		h.WithError(mockedError)
		rw := httptest.NewRecorder()

		require.Nil(t, WithErrors(log)(h.ServeHTTP)(rw, req))
		require.Equal(t, http.StatusBadRequest, rw.Code)
		require.True(t, h.Called())
	})

	t.Run("no error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := newHandlerStub()
		rw := httptest.NewRecorder()

		require.Nil(t, WithErrors(log)(h.ServeHTTP)(rw, req))
		require.Equal(t, http.StatusOK, rw.Code)
		require.True(t, h.Called())
	})
}

func TestWithOrganizer(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	web := NewWeb(logger.NewNoOpLogger(), osMock)
	require.NotNil(t, web)

	log := logger.NewNoOpLogger()

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

		h := newHandlerStub()
		rw := httptest.NewRecorder()

		require.NoError(t, WithOrganizer(osMock, log)(h.ServeHTTP)(rw, req))
		require.Equal(t, http.StatusOK, rw.Code)
		require.True(t, h.Called())
	})

	t.Run("no_organizer_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		h := newHandlerStub()
		rw := httptest.NewRecorder()

		e, ok := ErrorAs(WithOrganizer(osMock, log)(h.ServeHTTP)(rw, req))
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, e.StatusCode())

		require.False(t, h.Called())
	})

	t.Run("error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		mockedErr := errors.New("mocked error")
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(mockedErr)

		h := newHandlerStub()
		rw := httptest.NewRecorder()

		err = WithOrganizer(osMock, log)(h.ServeHTTP)(rw, req)
		e, ok := ErrorAs(err)
		require.True(t, ok)
		require.Error(t, e)
		require.Equal(t, http.StatusInternalServerError, e.StatusCode())

		require.False(t, h.Called())
	})
}

func TestWithCORS(t *testing.T) {
	web := NewWeb(logger.NewNoOpLogger(), nil)
	require.NotNil(t, web)

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		h := newHandlerStub()
		rw := httptest.NewRecorder()

		require.NoError(t, WithCORS(h.ServeHTTP)(rw, req))
		require.Equal(t, http.StatusOK, rw.Code)
		require.Equal(t, defaultOrigin, rw.Header().Get("Access-Control-Allow-Origin"))
		require.True(t, h.Called())
	})

	t.Run("no_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)

		h := func(rw http.ResponseWriter, r *http.Request) error { return nil }
		rw := httptest.NewRecorder()

		require.NoError(t, WithCORS(h)(rw, req))
		require.Equal(t, http.StatusOK, rw.Code)
		require.Equal(t, "", rw.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("wrong_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)
		req.Header.Set("Origin", "wrong_origin")

		h := func(rw http.ResponseWriter, r *http.Request) error { return nil }
		rw := httptest.NewRecorder()

		require.NoError(t, WithCORS(h)(rw, req))
		require.Equal(t, http.StatusOK, rw.Code)
		require.Equal(t, "", rw.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestWithRecover(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	log := logger.NewNoOpLogger()

	web := NewWeb(logger.NewNoOpLogger(), osMock)
	require.NotNil(t, web)

	t.Run("panic recovery", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := func(rw http.ResponseWriter, r *http.Request) error { panic("test panic") }

		rw := httptest.NewRecorder()

		e, ok := ErrorAs(WithRecover(log)(h)(rw, req))
		require.True(t, ok)
		require.ErrorIs(t, e.Value, ErrRecoveredFromPanic)

		require.Equal(t, http.StatusInternalServerError, e.StatusCode())
	})

	t.Run("no panic", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := func(rw http.ResponseWriter, r *http.Request) error {
			return nil
		}

		rw := httptest.NewRecorder()

		err = WithRecover(log)(h)(rw, req)
		require.Nil(t, err)

		require.Equal(t, http.StatusOK, rw.Code)
	})
}
