package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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

//go:generate mockgen -destination=mocks/mock_organizer.go -package=mocks github.com/kaznasho/yarmarok/service OrganizerService
//go:generate mockgen -destination=mocks/mock_raffle.go -package=mocks github.com/kaznasho/yarmarok/service RaffleService
//go:generate mockgen -destination=mocks/mock_contributor.go -package=mocks github.com/kaznasho/yarmarok/service ContributorService

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelError))
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("panic_in_handler", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleOrganizerIDHeader, organizerID)
		osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

		rsMock := mocks.NewMockRaffleService(ctrl)
		osMock.EXPECT().RaffleService(organizerID).Return(rsMock).Do(func(string) { panic("panic in handler") })

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("YARMAROK_ENDPOINT", func(t *testing.T) {
		t.Run("create_raffle", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				initRequest := &service.RaffleInitRequest{
					Name: "raffle_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(initRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				rsMock.EXPECT().Init(initRequest).Return(&service.InitResult{}, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)
			})

			t.Run("error", func(t *testing.T) {
				initRequest := &service.RaffleInitRequest{
					Name: "raffle_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(initRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				mockedErr := assert.AnError
				rsMock.EXPECT().Init(initRequest).Return(nil, mockedErr)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusInternalServerError, w.Code)
			})

			t.Run("empty_body", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, bytes.NewBuffer([]byte{}))
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusBadRequest, w.Code)
			})
		})

		t.Run("list_raffles", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				dummyTime := time.Now().UTC()
				expected := &service.RaffleListResponse{
					Raffles: []service.Raffle{
						{
							ID:        "raffle_id_1",
							Name:      "raffle_1",
							Note:      "note_1",
							CreatedAt: dummyTime,
						},
						{
							ID:        "raffle_id_2",
							Name:      "raffle_2",
							Note:      "note_2",
							CreatedAt: dummyTime,
						},
						{
							ID:        "raffle_id_3",
							Name:      "raffle_3",
							Note:      "note_3",
							CreatedAt: dummyTime,
						},
					},
				}

				req, err := newRequestWithOrigin(http.MethodGet, RafflesPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				rsMock.EXPECT().List().Return(expected, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)

				assertJSONResponse(t, expected, w.Body)

			})

			t.Run("error", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodGet, RafflesPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				mockedErr := assert.AnError
				rsMock.EXPECT().List().Return(nil, mockedErr)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusInternalServerError, w.Code)
			})
		})
	})

	t.Run("PARTICIPANT_ENDPOINT", func(t *testing.T) {
		organizerID := "organizer_id_1"
		raffleID := "raffle_id_1"
		contributorPath := "/raffles/raffle_id_1/contributors"

		t.Run("create_contributor", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				contributorInitRequest := &service.ContributorAddRequest{
					Name: "contributor_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(contributorInitRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, contributorPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				expect := &service.InitResult{ID: "contributor_id_1"}
				csMock.EXPECT().Add(contributorInitRequest).Return(expect, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)
			})

			t.Run("error", func(t *testing.T) {
				contributorInitRequest := &service.ContributorAddRequest{
					Name: "contributor_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(contributorInitRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, contributorPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				csMock.EXPECT().Add(contributorInitRequest).Return(nil, errors.New("test error"))

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusInternalServerError, w.Code)
			})

			t.Run("empty_body", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodPost, contributorPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusBadRequest, w.Code)
			})
		})

		t.Run("edit_contributor", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				contributorEditRequest := &service.ContributorEditRequest{
					ID:   "contributor_id_1",
					Name: "contributor_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(contributorEditRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPut, contributorPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				csMock.EXPECT().Edit(contributorEditRequest).Return(&service.Result{}, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)
			})

			t.Run("error", func(t *testing.T) {
				contributorEditRequest := &service.ContributorEditRequest{
					ID:   "contributor_id_1",
					Name: "contributor_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(contributorEditRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPut, contributorPath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				csMock.EXPECT().Edit(contributorEditRequest).Return(nil, errors.New("test error"))

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusInternalServerError, w.Code)
			})

			t.Run("empty_body", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodPut, contributorPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusBadRequest, w.Code)
			})
		})

		t.Run("list_contributors", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				dummyTime := time.Now().UTC()
				expected := &service.ContributorListResult{
					Contributors: []service.Contributor{
						{
							ID:        "contributor_id_1",
							Name:      "contributor_1",
							Phone:     "1323456789",
							Note:      "",
							CreatedAt: dummyTime,
						},
						{
							ID:        "contributor_id_2",
							Name:      "contributor_2",
							Phone:     "1323456789",
							Note:      "-",
							CreatedAt: dummyTime,
						},
						{
							ID:        "contributor_id_3",
							Name:      "contributor_3",
							Phone:     "1323456789",
							Note:      "bla bla bla",
							CreatedAt: dummyTime,
						},
					},
				}

				req, err := newRequestWithOrigin(http.MethodGet, contributorPath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				csMock.EXPECT().List().Return(expected, nil)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusOK, w.Code)
			})

			t.Run("error", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodGet, contributorPath, nil)
				require.NoError(t, err)

				req.Header.Set(GoogleOrganizerIDHeader, organizerID)
				osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				csMock := mocks.NewMockContributorService(ctrl)
				rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

				csMock.EXPECT().List().Return(nil, errors.New("test error"))

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusInternalServerError, w.Code)
			})
		})
	})

}

func TestApplyOrganizerMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	us := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(us, logger.NewNoOpLogger())
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleOrganizerIDHeader, organizerID)
		us.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		w := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		assert.True(t, stub.Called())
	})

	t.Run("no_organizer_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		w := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.False(t, stub.Called())
	})

	t.Run("error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		mockedErr := errors.New("mocked error")

		req.Header.Set(GoogleOrganizerIDHeader, organizerID)
		us.EXPECT().InitOrganizerIfNotExists(organizerID).Return(mockedErr)

		w := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
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

		w := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, defaultOrigin, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("no_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("wrong_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)

		req.Header.Set("Origin", "wrong_origin")

		w := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestGeContributorService(t *testing.T) {
	ctrl := gomock.NewController(t)

	organizerID := "organizer_id_1"
	raffleID := "raffle_id_1"

	osMock := mocks.NewMockOrganizerService(ctrl)
	rsMock := mocks.NewMockRaffleService(ctrl)
	csMock := mocks.NewMockContributorService(ctrl)

	osMock.EXPECT().InitOrganizerIfNotExists(organizerID).Return(nil).AnyTimes()
	osMock.EXPECT().RaffleService(organizerID).Return(rsMock).AnyTimes()

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelDebug))

	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodGet, "/raffles/raffle_id_1/contributors", nil)
		require.NoError(t, err)

		req.Header.Set(GoogleOrganizerIDHeader, organizerID)

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add(raffleIDParam, raffleID)

		rsMock.EXPECT().ContributorService(raffleID).Return(csMock)

		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(ctx)

		svc, err := router.getContributorService(req)
		assert.NoError(t, err)
		assert.Equal(t, svc, csMock)
	})

	t.Run("missing_organizer_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodGet, "/raffles/raffle_id_1/contributors", nil)
		require.NoError(t, err)

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add(raffleIDParam, raffleID)

		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing_raffle_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodGet, "/raffles//contributors", nil)
		require.NoError(t, err)

		req.Header.Set(GoogleOrganizerIDHeader, organizerID)

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add(raffleIDParam, "")

		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(ctx)

		svc, err := router.getContributorService(req)
		assert.Nil(t, svc)
		assert.ErrorIs(t, err, ErrMissingID)
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

	actualJSON, err := io.ReadAll(body)
	require.NoError(t, err)

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))

}

func newRequestWithOrigin(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Origin", defaultOrigin)

	return req, nil
}
