package function

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra"
	fsemulator "github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/kaznasho/yarmarok/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadRouter(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)
	log := logger.NewNoOpLogger()

	t.Run("empty_project_id", func(t *testing.T) {
		_, err := LoadRouter(log)
		require.ErrorIs(t, err, ErrEmptyProjectID)
	})

	t.Run("ok", func(t *testing.T) {
		firestoreInstance, err := fsemulator.RunInstance(t)
		require.NoError(t, err)

		t.Setenv(ProjectIDEnvVar, firestoreInstance.ProjectID())

		router, err := LoadRouter(log)
		require.NoError(t, err)

		req := dummyRequest(t)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		var resp web.CreateResponse

		require.Equal(t, http.StatusOK, recorder.Code)
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp), recorder.Body.String())

		require.NotEmpty(t, resp.ID)
	})

}

func TestEntrypoint(t *testing.T) {
	testinfra.SkipIfNotIntegrationRun(t)

	t.Run("no_project_id", func(t *testing.T) {
		req := dummyRequest(t)

		recorder := httptest.NewRecorder()

		Entrypoint(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
	t.Run("ok", func(t *testing.T) {
		firestoreInstance, err := fsemulator.RunInstance(t)
		require.NoError(t, err)

		t.Setenv(ProjectIDEnvVar, firestoreInstance.ProjectID())

		req := dummyRequest(t)

		recorder := httptest.NewRecorder()

		Entrypoint(recorder, req)

		var resp web.CreateResponse

		require.Equal(t, http.StatusOK, recorder.Code)
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp), recorder.Body.String())
		require.NotEmpty(t, resp.ID)
	})
}

func dummyRequest(t *testing.T) *http.Request {
	t.Helper()
	raffleInit := service.RaffleRequest{
		Name: "raffle_1",
		Note: "note_1",
	}

	jsonBody, err := json.Marshal(raffleInit)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, web.ApiPath+web.RafflesPath, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	req.Header.Set(web.GoogleUserIDHeader, "organizer_id_1")

	return req
}
