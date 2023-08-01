package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/mocks"
	"github.com/stretchr/testify/require"
)

func TestChainMiddlewares(t *testing.T) {
	const mwNum = 10

	var (
		key  = new(any)
		want string
	)

	mw := func(str string) Middleware {
		return func(next Handler) Handler {
			return func(rw http.ResponseWriter, req *http.Request) error {
				val, ok := req.Context().Value(key).(string)
				require.True(t, ok)
				val += str
				ctx := context.WithValue(req.Context(), key, val)
				return next(rw, req.WithContext(ctx))
			}
		}
	}

	mws := make([]Middleware, mwNum)

	for i := 0; i < mwNum; i++ {
		str := strconv.Itoa(i)
		mws[i] = mw(str)
		want += str
	}

	h := func(rw http.ResponseWriter, req *http.Request) error {
		_, err := rw.Write([]byte(fmt.Sprint(req.Context().Value(key))))
		require.NoError(t, err)
		return nil
	}

	h = WrapMiddlewares(h, mws...)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), key, "")
	err := h(rw, req.WithContext(ctx))
	require.NoError(t, err)
	require.Equal(t, want, rw.Body.String())
}

func TestWithErrors(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	log := logger.NewNoOpLogger()

	web := NewWeb(logger.NewNoOpLogger(), osMock)
	require.NotNil(t, web)

	t.Run("default_error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)

		mockedErr := errors.New("mocked error")
		h := newStubHandler(mockedErr)
		rw := httptest.NewRecorder()

		err = WithErrors(log)(h.ServeHTTP)(rw, req)
		require.Nil(t, err)

		require.Equal(t, http.StatusInternalServerError, rw.Code)
		require.True(t, h.Called())
	})

	t.Run("custom_error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		mockedErr := errors.New("mocked error")
		errweb := NewError(mockedErr, http.StatusBadRequest)
		h := newStubHandler(errweb)
		rw := httptest.NewRecorder()

		require.Nil(t, WithErrors(log)(h.ServeHTTP)(rw, req))
		require.Equal(t, http.StatusBadRequest, rw.Code)
		require.True(t, h.Called())
	})

	t.Run("no_error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := newStubHandler()
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

		h := newStubHandler()
		rw := httptest.NewRecorder()

		require.NoError(t, WithOrganizer(osMock, log)(h.ServeHTTP)(rw, req))
		require.Equal(t, http.StatusOK, rw.Code)
		require.True(t, h.Called())
	})

	t.Run("no_organizer_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		h := newStubHandler()
		rw := httptest.NewRecorder()

		err = WithOrganizer(osMock, log)(h.ServeHTTP)(rw, req)
		errweb, ok := ErrorAs(err)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, errweb.StatusCode())
		require.False(t, h.Called())
	})

	t.Run("error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		mockedErr := errors.New("mocked error")
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(mockedErr)

		h := newStubHandler()
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

		h := newStubHandler()
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

		err = WithCORS(h)(rw, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rw.Code)
		require.Equal(t, "", rw.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("wrong_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)
		req.Header.Set("Origin", "wrong_origin")

		h := func(rw http.ResponseWriter, r *http.Request) error { return nil }
		rw := httptest.NewRecorder()

		err = WithCORS(h)(rw, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rw.Code)
		require.Equal(t, "", rw.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestWithRecover(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	web := NewWeb(logger.NewNoOpLogger(), osMock)
	require.NotNil(t, web)

	t.Run("panic recovery", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := func(rw http.ResponseWriter, r *http.Request) error { panic("test panic") }
		rw := httptest.NewRecorder()

		err = WithRecover(h)(rw, req)
		errweb, ok := ErrorAs(err)
		require.True(t, ok)
		require.ErrorIs(t, errweb.value, ErrRecoveredFromPanic)
		require.Equal(t, http.StatusInternalServerError, errweb.StatusCode())
	})

	t.Run("no panic", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)
		req.Header.Set(GoogleUserIDHeader, organizerID)

		h := func(rw http.ResponseWriter, r *http.Request) error { return nil }
		rw := httptest.NewRecorder()

		err = WithRecover(h)(rw, req)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, rw.Code)
	})
}

type handlerStub struct {
	err    error
	called bool
	once   sync.Once
}

func newStubHandler(errs ...error) *handlerStub {
	return &handlerStub{
		once:   sync.Once{},
		called: false,
		err:    errors.Join(errs...),
	}
}

func (h *handlerStub) ServeHTTP(_ http.ResponseWriter, _ *http.Request) error {
	h.once.Do(func() {
		h.called = true
	})

	return h.err
}

func (h *handlerStub) Called() bool {
	return h.called
}
