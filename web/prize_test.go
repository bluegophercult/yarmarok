package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PrizeSuite struct {
	suite.Suite
	organizerService service.OrganizerService
	raffleService    *mocks.MockRaffleService
	prizeService     *mocks.MockPrizeService
	router           *Router
	organizerID      string
	raffleID         string
	prizeID          string
}

func TestPrize(t *testing.T) {
	suite.Run(t, &PrizeSuite{})
}

func (s *PrizeSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.organizerService = mocks.NewMockOrganizerService(ctrl)
	s.raffleService = mocks.NewMockRaffleService(ctrl)
	s.prizeService = mocks.NewMockPrizeService(ctrl)
	s.organizerID = "organizer_id_1"
	s.raffleID = "raffle_id_1"
	s.prizeID = "prize_id_1"

	s.organizerService.(*mocks.MockOrganizerService).EXPECT().CreateOrganizerIfNotExists(s.organizerID).Return(nil).AnyTimes()
	s.organizerService.(*mocks.MockOrganizerService).EXPECT().RaffleService(s.organizerID).Return(s.raffleService).AnyTimes()
	s.raffleService.EXPECT().PrizeService(s.raffleID).Return(s.prizeService).AnyTimes()

	var err error
	s.router, err = NewRouter(s.organizerService, logger.NewLogger(logger.LevelDebug))
	s.Require().NoError(err)
}

func (s *PrizeSuite) TestCreate() {

	prizePath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath)

	s.Run("success", func() {
		prizeNew := &service.PrizeRequest{
			Name:        "prize_1",
			Description: "description_1",
			TicketCost:  100_500,
		}

		req, err := newRequestJSON(http.MethodPost, prizePath, s.organizerID, prizeNew)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Create(prizeNew).Return(s.prizeID, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		prizeNew := &service.PrizeRequest{
			Name:        "prize_1",
			Description: "description_1",
			TicketCost:  100_500,
		}

		req, err := newRequestJSON(http.MethodPost, prizePath, s.organizerID, prizeNew)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Create(prizeNew).Return("", errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPost, prizePath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *PrizeSuite) TestEdit() {
	prizePath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID)

	s.Run("success", func() {
		prizeUpdate := &service.PrizeRequest{
			Name:        "prize_1",
			Description: "description_1",
			TicketCost:  100_500,
		}

		req, err := newRequestJSON(http.MethodPut, prizePath, s.organizerID, prizeUpdate)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Edit(s.prizeID, prizeUpdate).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		prizeUpdate := &service.PrizeRequest{
			Name:        "prize_1",
			Description: "description_1",
			TicketCost:  100_500,
		}

		req, err := newRequestJSON(http.MethodPut, prizePath, s.organizerID, prizeUpdate)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Edit(s.prizeID, prizeUpdate).Return(errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPut, prizePath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *PrizeSuite) TestList() {
	prizePath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath)

	s.Run("success", func() {
		dummyTime := time.Now().UTC()
		expected := []service.Prize{
			{
				ID:          "prize_id_1",
				Name:        "prize_1",
				Description: "description_1",
				TicketCost:  100_500,
				CreatedAt:   dummyTime,
			},
			{
				ID:          "prize_id_2",
				Name:        "prize_2",
				Description: "description_2",
				TicketCost:  200_500,
				CreatedAt:   dummyTime,
			},
			{
				ID:          "prize_id_3",
				Name:        "prize_3",
				Description: "description_3",
				TicketCost:  300_500,
				CreatedAt:   dummyTime,
			},
		}

		req, err := newRequestJSON(http.MethodGet, prizePath, s.organizerID, nil)
		s.Require().NoError(err)

		s.prizeService.EXPECT().List().Return(expected, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestJSON(http.MethodGet, prizePath, s.organizerID, nil)
		s.Require().NoError(err)

		s.prizeService.EXPECT().List().Return(nil, errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *PrizeSuite) TestGet() {
	prizePath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID)

	s.Run("success", func() {
		dummyTime := time.Now().UTC()
		expected := &service.Prize{
			ID:          "prize_id_1",
			Name:        "prize_1",
			Description: "description_1",
			TicketCost:  100_500,
			CreatedAt:   dummyTime,
		}

		req, err := newRequestJSON(http.MethodGet, prizePath, s.organizerID, nil)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Get(s.prizeID).Return(expected, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestJSON(http.MethodGet, prizePath, s.organizerID, nil)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Get(s.prizeID).Return(nil, errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}
func (s *PrizeSuite) TestDelete() {
	prizePath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID)

	s.Run("success", func() {
		req, err := newRequestJSON(http.MethodDelete, prizePath, s.organizerID, nil)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Delete(s.prizeID).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestJSON(http.MethodDelete, prizePath, s.organizerID, nil)
		s.Require().NoError(err)

		s.prizeService.EXPECT().Delete(s.prizeID).Return(errors.New("test error"))

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *PrizeSuite) TestPlay() {
	playPath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID, PlayPath)

	s.Run("success", func() {
		req, err := newRequestJSON(http.MethodGet, playPath, s.organizerID, nil)
		s.NoError(err)

		mockedTime := time.Now().UTC()

		mockedResponse := &service.PrizePlayResult{
			Winners: []service.PlayParticipant{
				{
					Participant: service.Participant{
						ID:        "ID1",
						Name:      "name1",
						Phone:     "phone1",
						Note:      "note1",
						CreatedAt: mockedTime,
					},
					TotalDonation:      300,
					TotalTicketsNumber: 10,
					Donations: []service.Donation{
						{
							ID:            "dID1",
							ParticipantID: "id1",
							Amount:        300,
							CreatedAt:     time.Time{},
						},
					},
				},
			},
			PlayParticipants: []service.PlayParticipant{
				{
					Participant: service.Participant{
						ID:        "ID2",
						Name:      "name2",
						Phone:     "phone2",
						Note:      "note2",
						CreatedAt: mockedTime,
					},
					TotalDonation:      200,
					TotalTicketsNumber: 5,
					Donations: []service.Donation{
						{
							ID:            "dID2",
							ParticipantID: "ID2",
							Amount:        200,
							CreatedAt:     mockedTime,
						},
					},
				},
				{
					Participant: service.Participant{
						ID:        "ID3",
						Name:      "name3",
						Phone:     "phone3",
						Note:      "note3",
						CreatedAt: mockedTime,
					},
					TotalDonation:      100,
					TotalTicketsNumber: 2,
					Donations: []service.Donation{
						{
							ID:            "dID3",
							ParticipantID: "ID3",
							Amount:        100,
							CreatedAt:     mockedTime,
						},
					},
				},
			},
		}

		s.prizeService.EXPECT().Play(s.prizeID).Return(mockedResponse, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
		s.Equal("application/json", writer.Header().Get("Content-Type"))
	})

	s.Run("error", func() {
		req, err := newRequestJSON(http.MethodGet, playPath, s.organizerID, nil)
		s.NoError(err)

		mockedErr := assert.AnError
		s.prizeService.EXPECT().Play(s.prizeID).Return(nil, mockedErr)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}
