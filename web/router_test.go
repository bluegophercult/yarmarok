package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/web/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=mocks/mock_user.go -package=mocks github.com/kaznasho/yarmarok/service UserService
//go:generate mockgen -destination=mocks/mock_yarmarok.go -package=mocks github.com/kaznasho/yarmarok/service YarmarokService
//go:generate mockgen -destination=mocks/mock_participant.go -package=mocks github.com/kaznasho/yarmarok/service ParticipantService

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	us := mocks.NewMockUserService(ctrl)
	userID := "user_id_1"

	router, err := NewRouter(us, logger.NewLogger(logger.LevelDebug))
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("panic_in_handler", func(t *testing.T) {
		req, err := http.NewRequest("POST", YarmaroksPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, userID)
		us.EXPECT().InitUserIfNotExists(userID).Return(nil)

		ysMock := mocks.NewMockYarmarokService(ctrl)
		us.EXPECT().YarmarokService(userID).Return(ysMock).Do(func(string) { panic("panic in handler") })

		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, req)
		require.Equal(t, http.StatusInternalServerError, writer.Code)
	})

	t.Run("create_yarmarok", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			initRequest := &service.YarmarokInitRequest{
				Name: "yarmarok_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(initRequest)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := http.NewRequest("POST", YarmaroksPath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, userID)
			us.EXPECT().InitUserIfNotExists(userID).Return(nil)

			ysMock := mocks.NewMockYarmarokService(ctrl)
			us.EXPECT().YarmarokService(userID).Return(ysMock)

			ysMock.EXPECT().Init(initRequest).Return(&service.InitResult{}, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
		})

		t.Run("error", func(t *testing.T) {
			initRequest := &service.YarmarokInitRequest{
				Name: "yarmarok_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(initRequest)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := http.NewRequest("POST", YarmaroksPath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, userID)
			us.EXPECT().InitUserIfNotExists(userID).Return(nil)

			ysMock := mocks.NewMockYarmarokService(ctrl)
			us.EXPECT().YarmarokService(userID).Return(ysMock)

			mockedErr := assert.AnError
			ysMock.EXPECT().Init(initRequest).Return(nil, mockedErr)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})

		t.Run("empty_body", func(t *testing.T) {
			req, err := http.NewRequest("POST", YarmaroksPath, bytes.NewBuffer([]byte{}))
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, userID)
			us.EXPECT().InitUserIfNotExists(userID).Return(nil)

			ysMock := mocks.NewMockYarmarokService(ctrl)
			us.EXPECT().YarmarokService(userID).Return(ysMock)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusBadRequest, writer.Code)
		})
	})

	t.Run("list_yarmaroks", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			dummyTime := time.Now().UTC()
			expected := &service.YarmarokListResponse{
				Yarmaroks: []service.Yarmarok{
					{
						ID:        "yarmarok_id_1",
						Name:      "yarmarok_1",
						Note:      "note_1",
						CreatedAt: dummyTime,
					},
					{
						ID:        "yarmarok_id_2",
						Name:      "yarmarok_2",
						Note:      "note_2",
						CreatedAt: dummyTime,
					},
					{
						ID:        "yarmarok_id_3",
						Name:      "yarmarok_3",
						Note:      "note_3",
						CreatedAt: dummyTime,
					},
				},
			}

			req, err := http.NewRequest("GET", YarmaroksPath, emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, userID)
			us.EXPECT().InitUserIfNotExists(userID).Return(nil)

			ysMock := mocks.NewMockYarmarokService(ctrl)
			us.EXPECT().YarmarokService(userID).Return(ysMock)

			ysMock.EXPECT().List().Return(expected, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)

			assertJSONResponse(t, expected, writer.Body)

		})

		t.Run("error", func(t *testing.T) {
			req, err := http.NewRequest("GET", YarmaroksPath, emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, userID)
			us.EXPECT().InitUserIfNotExists(userID).Return(nil)

			ysMock := mocks.NewMockYarmarokService(ctrl)
			us.EXPECT().YarmarokService(userID).Return(ysMock)

			mockedErr := assert.AnError
			ysMock.EXPECT().List().Return(nil, mockedErr)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("PARTICIPANT_ENDPOINT", func(t *testing.T) {
		userID := "user_id_1"
		yarmarokID := "yarmarok_id_1"
		participantPath := joinPath(YarmaroksPath, yarmarokID, ParticipantsPath)

		t.Run("create_participant", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				participantInitRequest := &service.ParticipantInitRequest{
					Name: "participant_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(participantInitRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := http.NewRequest("POST", participantPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				expect := &service.InitResult{ID: "participant_id_1"}
				psMock.EXPECT().Add(participantInitRequest).Return(expect, nil)

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusOK, writer.Code)
			})

			t.Run("error", func(t *testing.T) {
				participantInitRequest := &service.ParticipantInitRequest{
					Name: "participant_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(participantInitRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := http.NewRequest("POST", participantPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				psMock.EXPECT().Add(participantInitRequest).Return(nil, errors.New("test error"))

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusInternalServerError, writer.Code)
			})

			t.Run("empty_body", func(t *testing.T) {
				req, err := http.NewRequest("POST", participantPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusBadRequest, writer.Code)
			})
		})

		t.Run("edit_participant", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				participantEditRequest := &service.ParticipantEditRequest{
					ID:   "participant_id_1",
					Name: "participant_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(participantEditRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := http.NewRequest("PUT", participantPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				psMock.EXPECT().Edit(participantEditRequest).Return(&service.Response{}, nil)

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusOK, writer.Code)
			})

			t.Run("error", func(t *testing.T) {
				participantEditRequest := &service.ParticipantEditRequest{
					ID:   "participant_id_1",
					Name: "participant_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(participantEditRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := http.NewRequest("PUT", participantPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				psMock.EXPECT().Edit(participantEditRequest).Return(nil, errors.New("test error"))

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusInternalServerError, writer.Code)
			})

			t.Run("empty_body", func(t *testing.T) {
				req, err := http.NewRequest("PUT", participantPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusBadRequest, writer.Code)
			})
		})

		t.Run("list_participants", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				dummyTime := time.Now().UTC()
				expected := &service.ParticipantListResponse{
					Participants: []service.Participant{
						{
							ID:        "participant_id_1",
							Name:      "participant_1",
							Phone:     "1323456789",
							Note:      "",
							CreatedAt: dummyTime,
						},
						{
							ID:        "participant_id_2",
							Name:      "participant_2",
							Phone:     "1323456789",
							Note:      "-",
							CreatedAt: dummyTime,
						},
						{
							ID:        "participant_id_3",
							Name:      "participant_3",
							Phone:     "1323456789",
							Note:      "bla bla bla",
							CreatedAt: dummyTime,
						},
					},
				}

				req, err := http.NewRequest("GET", participantPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				psMock.EXPECT().List().Return(expected, nil)

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusOK, writer.Code)
			})

			t.Run("error", func(t *testing.T) {
				req, err := http.NewRequest("GET", participantPath, nil)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, userID)
				us.EXPECT().InitUserIfNotExists(userID).Return(nil)

				ysMock := mocks.NewMockYarmarokService(ctrl)
				us.EXPECT().YarmarokService(userID).Return(ysMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				ysMock.EXPECT().ParticipantService(yarmarokID).Return(psMock)

				psMock.EXPECT().List().Return(nil, errors.New("test error"))

				writer := httptest.NewRecorder()
				router.ServeHTTP(writer, req)
				require.Equal(t, http.StatusInternalServerError, writer.Code)
			})
		})
	})

}

func TestApplyUserMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	us := mocks.NewMockUserService(ctrl)
	userID := "user_id_1"

	router, err := NewRouter(us, logger.NewNoOpLogger())
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest("POST", YarmaroksPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, userID)
		us.EXPECT().InitUserIfNotExists(userID).Return(nil)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		writer := httptest.NewRecorder()
		router.userMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		assert.True(t, stub.Called())
	})

	t.Run("no_user_id", func(t *testing.T) {
		req, err := http.NewRequest("POST", YarmaroksPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		writer := httptest.NewRecorder()
		router.userMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusBadRequest, writer.Code)
		assert.False(t, stub.Called())
	})

	t.Run("error", func(t *testing.T) {
		req, err := http.NewRequest("POST", YarmaroksPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		mockedErr := errors.New("mocked error")

		req.Header.Set(GoogleUserIDHeader, userID)
		us.EXPECT().InitUserIfNotExists(userID).Return(mockedErr)

		writer := httptest.NewRecorder()
		router.userMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusInternalServerError, writer.Code)
		assert.False(t, stub.Called())
	})
}

type HandlerStub struct {
	called bool
	once   sync.Once
}

func (h *HandlerStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	actualJSON, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))

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

func TestParticipantMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := "user_id_1"
	yarmarokID := "yarmarok_id_1"

	us := mocks.NewMockUserService(ctrl)
	ys := mocks.NewMockYarmarokService(ctrl)
	ps := mocks.NewMockParticipantService(ctrl)

	us.EXPECT().YarmarokService(userID).Return(ys).AnyTimes()
	ys.EXPECT().ParticipantService(yarmarokID).Return(ps).AnyTimes()

	router, err := NewRouter(us, logger.NewNoOpLogger())
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/yarmaroks/123/participants", nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, userID)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("yarmarok_id", yarmarokID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			participantService, ok := req.Context().Value(participantServiceKey).(service.ParticipantService)
			require.True(t, ok)
			assert.Equal(t, ps, participantService)
			w.WriteHeader(http.StatusOK)
		})

		router.participantMiddleware(nextHandler).ServeHTTP(httptest.NewRecorder(), req)
	})

	t.Run("missing_user_id", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/yarmaroks/123/participants", nil)
		require.NoError(t, err)

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("yarmarok_id", yarmarokID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			t.Error("Next handler should not be called")
		})

		recorder := httptest.NewRecorder()
		router.participantMiddleware(nextHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("missing_yarmarok_id", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/yarmaroks//participants", nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, userID)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			t.Error("Next handler should not be called")
		})

		recorder := httptest.NewRecorder()
		router.participantMiddleware(nextHandler).ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
