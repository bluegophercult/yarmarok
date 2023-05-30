package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/web/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const googleUserIDHeader = "X-Goog-Authenticated-User-Id"

// ErrAmbiguousUserIDHeader is returned when
// the user id header is not set or is ambiguous.
var ErrAmbiguousUserIDHeader = errors.New("ambiguous user id format")

//go:generate mockgen -destination=mocks/mock_user.go -package=mocks github.com/kaznasho/yarmarok/service UserService
//go:generate mockgen -destination=mocks/mock_yarmarok.go -package=mocks github.com/kaznasho/yarmarok/service YarmarokService
func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	us := mocks.NewMockUserService(ctrl)
	userID := "user_id_1"

	router, err := NewRouter(us)
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("create_yarmarok", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			initRequest := &service.YarmarokInitRequest{
				Name: "yarmarok_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(initRequest)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := http.NewRequest("POST", "/create-yarmarok", body)
			require.NoError(t, err)

			req.Header.Set(googleUserIDHeader, userID)
			us.EXPECT().InitUserIfNotExists(userID).Return(nil)

			ysMock := mocks.NewMockYarmarokService(ctrl)
			us.EXPECT().YarmarokService(userID).Return(ysMock)

			ysMock.EXPECT().Init(initRequest).Return(&service.InitResult{}, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
		})

	})
}

// A convenient alias for chi.Router
type Router struct {
	chi.Router
	userService service.UserService
}

// NewRouter creates a new Router
func NewRouter(us service.UserService) (*Router, error) {
	router := &Router{
		Router:      chi.NewRouter(),
		userService: us,
	}

	router.Use(router.applyUserIDMiddleware)

	router.Post("/create-yarmarok", router.createYarmarok)

	return router, nil
}

func (r *Router) createYarmarok(w http.ResponseWriter, req *http.Request) {
	userID, err := extractUserID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	yarmarokService := r.userService.YarmarokService(userID)

	initRequest := &service.YarmarokInitRequest{}
	err = json.NewDecoder(req.Body).Decode(initRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := yarmarokService.Init(initRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (r *Router) applyUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, err := extractUserID(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = r.userService.InitUserIfNotExists(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func extractUserID(r *http.Request) (string, error) {
	ids := r.Header.Values(googleUserIDHeader)

	if len(ids) != 1 {
		return "", ErrAmbiguousUserIDHeader
	}

	id := ids[0]
	if id == "" {
		return "", ErrAmbiguousUserIDHeader
	}

	return id, nil
}
