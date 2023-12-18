package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"sync"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/web/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=mocks/mock_organizer.go -package=mocks github.com/kaznasho/yarmarok/service OrganizerService
//go:generate mockgen -destination=mocks/mock_raffle.go -package=mocks github.com/kaznasho/yarmarok/service RaffleService
//go:generate mockgen -destination=mocks/mock_participant.go -package=mocks github.com/kaznasho/yarmarok/service ParticipantService
//go:generate mockgen -destination=mocks/mock_prize.go -package=mocks github.com/kaznasho/yarmarok/service PrizeService
//go:generate mockgen -destination=mocks/mock_donation.go -package=mocks github.com/kaznasho/yarmarok/service DonationService

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelError))
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		loginPath := joinPath(ApiPath, "/login")

		req, err := newRequestWithOrigin(http.MethodPost, loginPath, emptyBody())
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, req)
		require.Equal(t, http.StatusSeeOther, writer.Code)
		require.Equal(t, "/", writer.Header().Get("Location"))
	})
}

func TestRecoverMiddleware(t *testing.T) {
	router, err := NewRouter(nil, logger.NewLogger(logger.LevelDebug))
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("panic_recovery", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, "/", nil)
		require.NoError(t, err)

		h := func(rw http.ResponseWriter, r *http.Request) { panic("test panic") }
		rw := httptest.NewRecorder()

		router.recoverMiddleware(http.HandlerFunc(h)).ServeHTTP(rw, req)
		require.Equal(t, http.StatusInternalServerError, rw.Code)
	})

	t.Run("no_panic", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, "/", nil)
		require.NoError(t, err)

		h := func(rw http.ResponseWriter, r *http.Request) {}
		rw := httptest.NewRecorder()

		router.recoverMiddleware(http.HandlerFunc(h)).ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Code)
	})
}

func TestApplyOrganizerMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewNoOpLogger())
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		writer := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		assert.True(t, stub.Called())
	})

	t.Run("no_organizer_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		writer := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusBadRequest, writer.Code)
		assert.False(t, stub.Called())
	})

	t.Run("error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		mockedErr := errors.New("mocked error")

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(mockedErr)

		writer := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusInternalServerError, writer.Code)
		assert.False(t, stub.Called())
	})
}

func TestCORSMiddleware(t *testing.T) {
	router, err := NewRouter(nil, logger.NewNoOpLogger())
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		writer := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		require.Equal(t, defaultOrigin, writer.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("no_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)

		writer := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		require.Equal(t, "", writer.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("wrong_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)

		req.Header.Set("Origin", "wrong_origin")

		writer := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		require.Equal(t, "", writer.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestJoinPath(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{"path", "subpath", "subsubpath"}, "/path/subpath/subsubpath"},
		{[]string{"/", "path", "/subpath/", "/subsubpath"}, "/path/subpath/subsubpath"},
		{[]string{"path"}, "/path"},
		{[]string{"/"}, "/"},
	}

	for _, testCase := range testCases {
		result := joinPath(testCase.input...)
		if result != testCase.expected {
			t.Errorf("joinPath(%v) = %v, expected %v", testCase.input, result, testCase.expected)
		}
	}
}

func joinPath(args ...string) string {
	return path.Clean("/" + path.Join(args...))
}

type HandlerStub struct {
	called bool
	once   sync.Once
}

func (h *HandlerStub) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	h.once.Do(func() {
		h.called = true
	})
}

func (h *HandlerStub) Called() bool {
	return h.called
}

func newHandlerStub() *HandlerStub {
	handler := &HandlerStub{
		once:   sync.Once{},
		called: false,
	}

	return handler
}

func emptyBody() io.Reader {
	return bytes.NewReader([]byte{})
}

func assertJSONResponse(t *testing.T, expected interface{}, body io.Reader) {
	t.Helper()

	actualJSON, err := io.ReadAll(body)
	require.NoError(t, err)

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))

}

func newRequestJSON(method, url, organizerID string, raffleNew any) (*http.Request, error) {
	encoded, err := json.Marshal(raffleNew)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(encoded)

	req, err := newRequestWithOrigin(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(GoogleUserIDHeader, organizerID)
	return req, nil
}

func newRequestWithOrigin(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Origin", defaultOrigin)

	return req, nil
}
