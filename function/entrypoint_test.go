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

		t.Setenv(projectIDEnvVar, firestoreInstance.ProjectID())

		router, err := LoadRouter(log)
		require.NoError(t, err)

		req := dummyRequest(t)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		resp := service.InitResult{}

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp), recorder.Body.String())

		assert.NotEmpty(t, resp.ID)
	})

}

func TestEntrypoint(t *testing.T) {
	t.Run("no_project_id", func(t *testing.T) {
		req := dummyRequest(t)

		recorder := httptest.NewRecorder()

		Entrypoint(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
	t.Run("ok", func(t *testing.T) {
		firestoreInstance, err := fsemulator.RunInstance(t)
		require.NoError(t, err)

		t.Setenv(projectIDEnvVar, firestoreInstance.ProjectID())

		req := dummyRequest(t)

		recorder := httptest.NewRecorder()

		Entrypoint(recorder, req)

		resp := service.InitResult{}

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp), recorder.Body.String())

		assert.NotEmpty(t, resp.ID)
	})
}

func dummyRequest(t *testing.T) *http.Request {
	t.Helper()
	yarmarokInit := service.YarmarokInitRequest{
		Name: "yarmarok_1",
		Note: "note_1",
	}

	jsonBody, err := json.Marshal(yarmarokInit)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/create-yarmarok", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	req.Header.Set(web.GoogleUserIDHeader, "user_id_1")

	return req
}
