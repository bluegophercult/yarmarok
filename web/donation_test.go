package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/web/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type DonationSuite struct {
	suite.Suite
	organizerService *mocks.MockOrganizerService
	raffleService    *mocks.MockRaffleService
	prizeService     *mocks.MockPrizeService
	donationService  *mocks.MockDonationService
	router           *Router
	organizerID      string
	raffleID         string
	prizeID          string
	donationID       string
}

func TestDonation(t *testing.T) {
	suite.Run(t, &DonationSuite{})
}

func (s *DonationSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.organizerService = mocks.NewMockOrganizerService(ctrl)
	s.raffleService = mocks.NewMockRaffleService(ctrl)
	s.prizeService = mocks.NewMockPrizeService(ctrl)
	s.donationService = mocks.NewMockDonationService(ctrl)
	s.organizerID = "organizer_id_1"
	s.raffleID = "raffle_id_1"
	s.prizeID = "participant_id_1"
	s.donationID = "donation_id_1"

	s.organizerService.EXPECT().CreateOrganizerIfNotExists(s.organizerID).Return(nil).AnyTimes()
	s.organizerService.EXPECT().RaffleService(s.organizerID).Return(s.raffleService).AnyTimes()
	s.raffleService.EXPECT().PrizeService(s.raffleID).Return(s.prizeService).AnyTimes()
	s.prizeService.EXPECT().DonationService(s.prizeID).Return(s.donationService, nil).AnyTimes()

	var err error
	s.router, err = NewRouter(s.organizerService, logger.NewLogger(logger.LevelDebug))
	s.Require().NoError(err)
}

func (s *DonationSuite) TestCreate() {

	donationPath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID, DonationsPath)

	s.Run("success", func() {
		donationNew := &service.DonationRequest{
			Amount:        100,
			ParticipantID: "participant_id_1",
		}

		req, err := newRequestJSON(http.MethodPost, donationPath, s.organizerID, donationNew)
		s.Require().NoError(err)

		s.donationService.EXPECT().Create(donationNew).Return(s.donationID, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		donationNew := &service.DonationRequest{
			Amount:        100,
			ParticipantID: "participant_id_1",
		}

		req, err := newRequestJSON(http.MethodPost, donationPath, s.organizerID, donationNew)
		s.Require().NoError(err)

		s.donationService.EXPECT().Create(donationNew).Return("", service.ErrDonationAlreadyExists)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPost, donationPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *DonationSuite) TestEdit() {
	donationPath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID, DonationsPath, s.donationID)

	s.Run("success", func() {
		donationEdit := &service.DonationRequest{
			Amount:        100,
			ParticipantID: "participant_id_1",
		}

		req, err := newRequestJSON(http.MethodPut, donationPath, s.organizerID, donationEdit)
		s.Require().NoError(err)

		s.donationService.EXPECT().Edit(s.donationID, donationEdit).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		donationEdit := &service.DonationRequest{
			Amount:        100,
			ParticipantID: "participant_id_1",
		}

		req, err := newRequestJSON(http.MethodPut, donationPath, s.organizerID, donationEdit)
		s.Require().NoError(err)

		s.donationService.EXPECT().Edit(s.donationID, donationEdit).Return(service.ErrDonationNotFound)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPut, donationPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *DonationSuite) TestDelete() {
	donationPath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID, DonationsPath, s.donationID)

	s.Run("success", func() {
		req, err := newRequestJSON(http.MethodDelete, donationPath, s.organizerID, nil)
		s.Require().NoError(err)

		s.donationService.EXPECT().Delete(s.donationID).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestJSON(http.MethodDelete, donationPath, s.organizerID, nil)
		s.Require().NoError(err)

		s.donationService.EXPECT().Delete(s.donationID).Return(service.ErrDonationNotFound)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *DonationSuite) TestList() {
	donationPath := joinPath(ApiPath, RafflesPath, s.raffleID, PrizesPath, s.prizeID, DonationsPath)

	s.Run("success", func() {
		req, err := newRequestWithOrigin(http.MethodGet, donationPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		donations := []service.Donation{
			{
				ID:     "donation_id_1",
				Amount: 100,
			},
			{
				ID:     "donation_id_2",
				Amount: 200,
			},
		}

		s.donationService.EXPECT().List().Return(donations, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestWithOrigin(http.MethodGet, donationPath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		s.donationService.EXPECT().List().Return(nil, service.ErrDonationNotFound)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)

		s.Require().Equal(http.StatusInternalServerError, writer.Code)
	})
}
