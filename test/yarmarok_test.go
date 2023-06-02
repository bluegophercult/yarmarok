package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kaznasho/yarmarok/function"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/testinfra/firestore"
	"github.com/kaznasho/yarmarok/web"
	"github.com/stretchr/testify/suite"
)

type YarmarokSuite struct {
	suite.Suite
	firestoreInstance *firestore.Instance
	userID            string
}

func TestYarmarok(t *testing.T) {
	//testinfra.SkipIfNotIntegrationRun(t)
	suite.Run(t, &YarmarokSuite{})
}

func (s *YarmarokSuite) SetupSuite() {
	var err error
	s.firestoreInstance, err = firestore.RunInstance(s.T())
	s.Require().NoError(err)

	s.T().Setenv(function.ProjectIDEnvVar, s.firestoreInstance.ProjectID())

	s.userID = "test-user-id"
}

func (s *YarmarokSuite) TestYarmarok() {
	initRequest := &service.YarmarokInitRequest{
		Name: "Розіграш на фестивалі яблук",
		Note: "Благодійний розіграш призів на фестивалі яблукб проходитиме 17 серпня 2023 року з 9:00 по 16:00.",
	}

	response := &service.InitResult{}

	s.post("/create-yarmarok", initRequest, response)
	s.Require().NotEmpty(response.ID)
}

func (s *YarmarokSuite) post(path string, body, response interface{}) {
	s.executeRequest(http.MethodPost, path, body, response)
}

func (s *YarmarokSuite) get(path string, response interface{}) {
	s.executeRequest(http.MethodPost, path, struct{}{}, response)
}

func (s *YarmarokSuite) executeRequest(method, path string, body, response interface{}) {
	s.Require().True(isPointer(response))
	data, err := json.Marshal(body)
	s.Require().NoError(err)

	req, err := http.NewRequest(method, path, bytes.NewBuffer(data))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(web.GoogleUserIDHeader, s.userID)

	rr := httptest.NewRecorder()
	function.Entrypoint(rr, req)

	s.Require().Equal(http.StatusOK, rr.Code)

	s.Require().NoError(json.Unmarshal(rr.Body.Bytes(), response))
}

func isPointer(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
