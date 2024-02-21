package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/web/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ParticipantSuite struct {
	suite.Suite
	organizerService   service.OrganizerService
	raffleService      *mocks.MockRaffleService
	participantService *mocks.MockParticipantService
	router             *Router
	organizerID        string
	raffleID           string
	participantID      string
}

func TestParticipant(t *testing.T) {
	suite.Run(t, &ParticipantSuite{})
}

func (s *ParticipantSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.organizerService = mocks.NewMockOrganizerService(ctrl)
	s.raffleService = mocks.NewMockRaffleService(ctrl)
	s.participantService = mocks.NewMockParticipantService(ctrl)
	s.organizerID = "organizer_id_1"
	s.raffleID = "raffle_id_1"
	s.participantID = "participant_id_1"

	s.organizerService.(*mocks.MockOrganizerService).EXPECT().CreateOrganizerIfNotExists(s.organizerID).Return(nil).AnyTimes()
	s.organizerService.(*mocks.MockOrganizerService).EXPECT().RaffleService(s.organizerID).Return(s.raffleService).AnyTimes()
	s.raffleService.EXPECT().ParticipantService(s.raffleID).Return(s.participantService).AnyTimes()

	var err error
	s.router, err = NewRouter(s.organizerService, logger.NewLogger(logger.LevelDebug))
	s.Require().NoError(err)
}

func (s *ParticipantSuite) TestCreate() {
	participantPath := joinPath(ApiPath, RafflesPath, s.raffleID, ParticipantsPath)

	s.Run("success", func() {
		participantCreateRequest := &service.ParticipantRequest{
			Name: "participant_1",
			Note: "note_1",
		}

		req, err := newRequestJSON(http.MethodPost, participantPath, s.organizerID, participantCreateRequest)
		s.Require().NoError(err)

		s.participantService.EXPECT().Create(participantCreateRequest).Return(s.participantID, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		participantCreateRequest := &service.ParticipantRequest{Name: "participant_1", Phone: "phone_1", Note: "note_1"}

		req, err := newRequestJSON(http.MethodPost, participantPath, s.organizerID, participantCreateRequest)
		s.Require().NoError(err)
		s.participantService.EXPECT().Create(participantCreateRequest).Return("", errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPost, participantPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *ParticipantSuite) TestEdit() {
	participantPath := joinPath(ApiPath, RafflesPath, s.raffleID, ParticipantsPath, s.participantID)

	s.Run("success", func() {
		participantEditRequest := &service.ParticipantRequest{
			Name: "participant_1",
			Note: "note_1",
		}

		req, err := newRequestJSON(http.MethodPut, participantPath, s.organizerID, participantEditRequest)
		s.Require().NoError(err)

		s.participantService.EXPECT().Edit(s.participantID, participantEditRequest).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		participantEditRequest := &service.ParticipantRequest{Name: "participant_1", Phone: "phone_1", Note: "note_1"}

		req, err := newRequestJSON(http.MethodPut, participantPath, s.organizerID, participantEditRequest)
		s.Require().NoError(err)
		s.participantService.EXPECT().Edit(s.participantID, participantEditRequest).Return(errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPut, participantPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *ParticipantSuite) TestDelete() {
	participantPath := joinPath(ApiPath, RafflesPath, s.raffleID, ParticipantsPath, s.participantID)

	s.Run("success", func() {
		req, err := newRequestWithOrigin(http.MethodDelete, participantPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		s.participantService.EXPECT().Delete(s.participantID).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestWithOrigin(http.MethodDelete, participantPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		s.participantService.EXPECT().Delete(s.participantID).Return(errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *ParticipantSuite) TestList() {
	participantPath := joinPath(ApiPath, RafflesPath, s.raffleID, ParticipantsPath)

	s.Run("success", func() {
		req, err := newRequestWithOrigin(http.MethodGet, participantPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		participants := []service.Participant{
			{
				ID:   "participant_id_1",
				Name: "participant_1",
				Note: "note_1",
			},
			{
				ID:   "participant_id_2",
				Name: "participant_2",
				Note: "note_2",
			},
		}
		s.participantService.EXPECT().List().Return(participants, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestWithOrigin(http.MethodGet, participantPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		s.participantService.EXPECT().List().Return(nil, errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}
