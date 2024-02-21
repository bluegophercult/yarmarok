package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ParticipantSuite struct {
	suite.Suite

	ctrl     *gomock.Controller
	storage  *MockParticipantStorage
	manager  *ParticipantManager
	mockUUID string
	mockTime time.Time
}

func TestParticipant(t *testing.T) {
	suite.Run(t, &ParticipantSuite{})
}

func (s *ParticipantSuite) SetupTest() {
	s.mockUUID = "participant_id"
	s.mockTime = time.Now().UTC()
	setTimeNowMock(s.mockTime)
	setUUIDMock(s.mockUUID)

	s.ctrl = gomock.NewController(s.T())
	s.storage = NewMockParticipantStorage(s.ctrl)
	s.manager = NewParticipantManager(s.storage)
}

func (s *ParticipantSuite) TestCreateParticipant() {
	participantRequest := dummyParticipantRequest()

	mockedParticipant := &Participant{
		ID:        s.mockUUID,
		Name:      participantRequest.Name,
		Phone:     participantRequest.Phone,
		Note:      participantRequest.Note,
		CreatedAt: s.mockTime,
	}

	s.storage.EXPECT().Create(mockedParticipant).Return(nil)

	resID, err := s.manager.Create(participantRequest)
	require.NoError(s.T(), err)
	require.Equal(s.T(), mockedParticipant.ID, resID)

	s.Run("error", func() {
		s.storage.EXPECT().Create(mockedParticipant).Return(errors.New("test error"))

		_, err := s.manager.Create(participantRequest)
		require.Error(s.T(), err)
	})

	s.Run("invalid name", func() {
		participantRequest := dummyParticipantRequest()
		participantRequest.Name = "a"

		_, err := s.manager.Create(participantRequest)
		require.Error(s.T(), err)
	})

	s.Run("invalid phone", func() {
		participantRequest := dummyParticipantRequest()
		participantRequest.Phone = "123"

		_, err := s.manager.Create(participantRequest)
		require.Error(s.T(), err)
	})

	s.Run("invalid note", func() {
		participantRequest := dummyParticipantRequest()
		participantRequest.Note = "///"

		_, err := s.manager.Create(participantRequest)
		require.Error(s.T(), err)
	})
}

func (s *ParticipantSuite) TestEditParticipant() {
	participantRequest := dummyParticipantRequest()
	participant := &Participant{
		ID:        s.mockUUID,
		Name:      participantRequest.Name,
		Phone:     participantRequest.Phone,
		Note:      participantRequest.Note,
		CreatedAt: s.mockTime,
	}

	s.storage.EXPECT().Get(participant.ID).Return(participant, nil)
	s.storage.EXPECT().Update(participant).Return(nil)

	err := s.manager.Edit(participant.ID, participantRequest)
	require.NoError(s.T(), err)

	s.Run("error", func() {
		s.storage.EXPECT().Get(participant.ID).Return(nil, errors.New("test error"))

		err := s.manager.Edit(participant.ID, participantRequest)
		require.Error(s.T(), err)
	})

	s.Run("error_in_update", func() {
		s.storage.EXPECT().Get(participant.ID).Return(participant, nil)
		s.storage.EXPECT().Update(participant).Return(errors.New("test error"))

		err := s.manager.Edit(participant.ID, participantRequest)
		require.Error(s.T(), err)
	})

	s.Run("invalid name", func() {
		participantRequest := dummyParticipantRequest()
		participantRequest.Name = "a"

		err := s.manager.Edit(participant.ID, participantRequest)
		require.Error(s.T(), err)
	})

	s.Run("invalid phone", func() {
		participantRequest := dummyParticipantRequest()
		participantRequest.Phone = "123"

		err := s.manager.Edit(participant.ID, participantRequest)
		require.Error(s.T(), err)
	})

	s.Run("invalid note", func() {
		participantRequest := dummyParticipantRequest()
		participantRequest.Note = "///"

		err := s.manager.Edit(participant.ID, participantRequest)
		require.Error(s.T(), err)
	})
}

func (s *ParticipantSuite) TestDeleteParticipant() {
	participant := dummyParticipant()

	s.storage.EXPECT().Delete(participant.ID).Return(nil)

	err := s.manager.Delete(participant.ID)
	require.NoError(s.T(), err)

	s.Run("error", func() {
		s.storage.EXPECT().Delete(participant.ID).Return(errors.New("test error"))

		err := s.manager.Delete(participant.ID)
		require.Error(s.T(), err)
	})
}

func (s *ParticipantSuite) TestListParticipant() {
	participant := dummyParticipant()
	participants := []Participant{*participant}

	s.storage.EXPECT().GetAll().Return(participants, nil)

	res, err := s.manager.List()
	require.NoError(s.T(), err)
	require.Equal(s.T(), participants, res)

	s.Run("error", func() {
		s.storage.EXPECT().GetAll().Return(nil, errors.New("test error"))

		_, err := s.manager.List()
		require.Error(s.T(), err)
	})
}

func dummyParticipantRequest() *ParticipantRequest {
	return &ParticipantRequest{
		Name:  "John Doe",
		Phone: "+380123456789",
		Note:  "Test participant",
	}
}

func dummyParticipant() *Participant {
	return &Participant{
		ID:        stringUUID(),
		Name:      "John Doe",
		Phone:     "+380123456789",
		Note:      "Test participant",
		CreatedAt: timeNow(),
	}
}
