package web

import (
	"bytes"
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

type RaffleSuite struct {
	suite.Suite
	organizerService service.OrganizerService
	raffleService    *mocks.MockRaffleService
	router           *Router
	organizerID      string
}

func (s *RaffleSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.organizerService = mocks.NewMockOrganizerService(ctrl)
	s.raffleService = mocks.NewMockRaffleService(ctrl)
	s.organizerID = "organizer_id_1"

	s.organizerService.(*mocks.MockOrganizerService).EXPECT().CreateOrganizerIfNotExists(s.organizerID).Return(nil).AnyTimes()
	s.organizerService.(*mocks.MockOrganizerService).EXPECT().RaffleService(s.organizerID).Return(s.raffleService).AnyTimes()

	var err error
	s.router, err = NewRouter(s.organizerService, logger.NewLogger(logger.LevelDebug))
	s.Require().NoError(err)
}

func TestRaffle(t *testing.T) {
	suite.Run(t, &RaffleSuite{})
}

func (s *RaffleSuite) TestCreate() {

	raffleID := "raffle_id_1"
	rafflePath := joinPath(ApiPath, RafflesPath)

	s.Run("success", func() {
		raffleNew := &service.RaffleRequest{
			Name: "raffle_1",
			Note: "note_1",
		}

		req, err := newRequestJSON(http.MethodPost, rafflePath, s.organizerID, raffleNew)
		s.Require().NoError(err)

		s.raffleService.EXPECT().Create(raffleNew).Return(raffleID, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
		assertJSONResponse(s.T(), CreateResponse{raffleID}, writer.Body)
	})

	s.Run("error", func() {
		raffleNew := &service.RaffleRequest{
			Name: "raffle_1",
			Note: "note_1",
		}

		req, err := newRequestJSON(http.MethodPost, rafflePath, s.organizerID, raffleNew)
		s.Require().NoError(err)

		mockedErr := assert.AnError
		s.raffleService.EXPECT().Create(raffleNew).Return("", mockedErr)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPost, rafflePath, bytes.NewBuffer([]byte{}))
		s.Require().NoError(err)
		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *RaffleSuite) TestEdit() {
	raffleID := "raffle_id_1"
	rafflePath := joinPath(ApiPath, RafflesPath, raffleID)

	s.Run("success", func() {
		raffleUpd := &service.RaffleRequest{
			Name: "raffle_1",
			Note: "note_1",
		}

		req, err := newRequestJSON(http.MethodPut, rafflePath, s.organizerID, raffleUpd)
		s.Require().NoError(err)

		s.raffleService.EXPECT().Edit(raffleID, raffleUpd).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		upd := &service.RaffleRequest{
			Name: "raffle_1",
			Note: "note_1",
		}

		req, err := newRequestJSON(http.MethodPut, rafflePath, s.organizerID, upd)
		s.Require().NoError(err)

		mockedErr := assert.AnError
		s.raffleService.EXPECT().Edit(raffleID, upd).Return(mockedErr)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})

	s.Run("empty_body", func() {
		req, err := newRequestWithOrigin(http.MethodPut, rafflePath, bytes.NewBuffer([]byte{}))
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *RaffleSuite) TestDelete() {
	raffleID := "raffle_id_1"
	rafflePath := joinPath(ApiPath, RafflesPath, raffleID)

	s.Run("success", func() {
		req, err := newRequestWithOrigin(http.MethodDelete, rafflePath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		s.raffleService.EXPECT().Delete(raffleID).Return(nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
	})

	s.Run("error", func() {
		req, err := newRequestWithOrigin(http.MethodDelete, rafflePath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		mockedErr := assert.AnError
		s.raffleService.EXPECT().Delete(raffleID).Return(mockedErr)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *RaffleSuite) TestList() {
	rafflePath := joinPath(ApiPath, RafflesPath)

	s.Run("success", func() {
		dummyTime := time.Now().UTC()
		raffles := []service.Raffle{
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
		}

		req, err := newRequestWithOrigin(http.MethodGet, rafflePath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		s.raffleService.EXPECT().List().Return(raffles, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
		assertJSONResponse(s.T(), ListResponse[service.Raffle]{raffles}, writer.Body)
	})

	s.Run("error", func() {
		req, err := newRequestWithOrigin(http.MethodGet, rafflePath, emptyBody())
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		mockedErr := assert.AnError
		s.raffleService.EXPECT().List().Return(nil, mockedErr)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}

func (s *RaffleSuite) TestDownloadXLSX() {
	raffleID := "raffle_id_1"
	downloadPath := joinPath(ApiPath, RafflesPath, raffleID, "/download-xlsx")

	s.Run("success", func() {
		req, err := newRequestWithOrigin(http.MethodGet, downloadPath, nil)
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		s.raffleService.EXPECT().Export(raffleID).Return(
			&service.RaffleExportResult{
				FileName: "raffle.xlsx",
				Content:  []byte("content")}, nil)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusOK, writer.Code)
		s.Equal("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", writer.Header().Get("Content-Type"))
		s.Equal("attachment; filename=raffle.xlsx", writer.Header().Get("Content-Disposition"))
	})

	s.Run("error", func() {
		req, err := newRequestWithOrigin(http.MethodGet, downloadPath, nil)
		s.Require().NoError(err)

		req.Header.Set(GoogleUserIDHeader, s.organizerID)

		mockedErr := assert.AnError
		s.raffleService.EXPECT().Export(raffleID).Return(nil, mockedErr)

		writer := httptest.NewRecorder()
		s.router.ServeHTTP(writer, req)
		s.Equal(http.StatusInternalServerError, writer.Code)
	})
}
