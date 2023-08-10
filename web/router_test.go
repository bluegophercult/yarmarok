package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path"
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
//go:generate mockgen -destination=mocks/mock_participant.go -package=mocks github.com/kaznasho/yarmarok/service ParticipantService
//go:generate mockgen -destination=mocks/mock_prize.go -package=mocks github.com/kaznasho/yarmarok/service PrizeService
//go:generate mockgen -destination=mocks/mock_donation.go -package=mocks github.com/kaznasho/yarmarok/service DonationService

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelError))
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		loginPath := joinPath(ApiPath, "/login")

		req, err := newRequestWithOrigin(http.MethodPost, loginPath, emptyBody())
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, req)
		require.Equal(t, http.StatusSeeOther, writer.Code)
		require.Equal(t, "/", writer.Header().Get("Location"))
	})
}

func TestRaffle(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelDebug))
	require.NoError(t, err)
	require.NotNil(t, router)

	raffleID := "raffle_id_1"
	rafflePath := joinPath(ApiPath, RafflesPath)

	t.Run("create_raffle", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			raffleNew := &service.RaffleRequest{
				Name: "raffle_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(raffleNew)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPost, rafflePath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			rsMock.EXPECT().Create(raffleNew).Return(raffleID, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
			assertJSONResponse(t, CreateResponse{raffleID}, writer.Body)
		})

		t.Run("error", func(t *testing.T) {
			raffleNew := &service.RaffleRequest{
				Name: "raffle_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(raffleNew)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPost, rafflePath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			mockedErr := assert.AnError
			rsMock.EXPECT().Create(raffleNew).Return("", mockedErr)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})

		t.Run("empty_body", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodPost, rafflePath, bytes.NewBuffer([]byte{}))
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("edit_raffle", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			raffleUpd := &service.RaffleRequest{
				Name: "raffle_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(raffleUpd)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPut, joinPath(rafflePath, raffleID), body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			rsMock.EXPECT().Edit(raffleID, raffleUpd).Return(nil)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, http.StatusOK, rw.Code)
		})

		t.Run("error", func(t *testing.T) {
			upd := &service.RaffleRequest{
				Name: "raffle_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(upd)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPut, joinPath(rafflePath, raffleID), body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			rsMock.EXPECT().Edit(raffleID, upd).Return(errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})

		t.Run("empty_body", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodPut, joinPath(rafflePath, raffleID), emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("list_raffles", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
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
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			rsMock.EXPECT().List().Return(raffles, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
			assertJSONResponse(t, ListResponse[service.Raffle]{raffles}, writer.Body)
		})

		t.Run("error", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodGet, rafflePath, emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			mockedErr := assert.AnError
			rsMock.EXPECT().List().Return(nil, mockedErr)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("download_xlsx", func(t *testing.T) {
		downloadPath := joinPath(rafflePath, raffleID, "/download-xlsx")

		t.Run("success", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodGet, downloadPath, nil)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			rsMock.EXPECT().Export(raffleID).Return(
				&service.RaffleExportResult{
					FileName: "raffle.xlsx",
					Content:  []byte("content")}, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
			require.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", writer.Header().Get("Content-Type"))
			require.Equal(t, "attachment; filename=raffle.xlsx", writer.Header().Get("Content-Disposition"))
		})

		t.Run("error", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodGet, downloadPath, nil)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			mockedErr := assert.AnError
			rsMock.EXPECT().Export(raffleID).Return(nil, mockedErr)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})
}

func TestParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelDebug))
	require.NoError(t, err)
	require.NotNil(t, router)

	raffleID := "raffle_id_1"
	participantID := "participant_id_1"
	participantPath := joinPath(ApiPath, RafflesPath, raffleID, ParticipantsPath)

	t.Run("create_participant", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			participantInitRequest := &service.ParticipantRequest{
				Name: "participant_1",
				Note: "note_1",
			}

			encoded, err := json.Marshal(participantInitRequest)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPost, participantPath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().Create(participantInitRequest).Return(participantID, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
		})

		t.Run("error", func(t *testing.T) {
			prt := &service.ParticipantRequest{Name: "participant_1", Phone: "phone_1", Note: "note_1"}

			encoded, err := json.Marshal(prt)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPost, participantPath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().Create(prt).Return("", errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})

		t.Run("empty_body", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodPost, participantPath, emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("edit_participant", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			participantEditRequest := &service.ParticipantRequest{
				Name:  "participant_1",
				Note:  "note_1",
				Phone: "1234567890",
			}

			encoded, value := json.Marshal(participantEditRequest)
			require.NoError(t, value)

			body := bytes.NewReader(encoded)

			req, value := newRequestWithOrigin(http.MethodPut, joinPath(participantPath, participantID), body)
			require.NoError(t, value)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().Edit(participantID, participantEditRequest).Return(nil)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			log.Printf("%v\n", req)
			require.Equal(t, http.StatusOK, rw.Code)
		})

		t.Run("error", func(t *testing.T) {
			upd := &service.ParticipantRequest{Name: "participant_1", Phone: "phone_1", Note: "note_1"}

			encoded, err := json.Marshal(upd)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPut, joinPath(participantPath, participantID), body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().Edit(participantID, upd).Return(errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})

		t.Run("empty_body", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodPut, joinPath(participantPath, participantID), emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("list_participants", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			dummyTime := time.Now().UTC()
			expected := []service.Participant{
				{
					ID:        "participant_id_1",
					Name:      "participant_1",
					Phone:     "1323456789",
					Note:      "",
					CreatedAt: dummyTime,
				},
				{
					ID:        "participant_id_2",
					Name:      "participant_2",
					Phone:     "1323456789",
					Note:      "-",
					CreatedAt: dummyTime,
				},
				{
					ID:        "participant_id_3",
					Name:      "participant_3",
					Phone:     "1323456789",
					Note:      "bla bla bla",
					CreatedAt: dummyTime,
				},
			}

			req, err := newRequestWithOrigin(http.MethodGet, participantPath, emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().List().Return(expected, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
		})

		t.Run("error", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodGet, participantPath, nil)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().List().Return(nil, errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("delete_participant", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			req, value := newRequestWithOrigin(http.MethodDelete, joinPath(participantPath, participantID), nil)
			require.NoError(t, value)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().Delete(participantID).Return(nil)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, http.StatusOK, rw.Code)
		})

		t.Run("error", func(t *testing.T) {
			req, value := newRequestWithOrigin(http.MethodDelete, joinPath(participantPath, participantID), nil)
			require.NoError(t, value)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			psMock := mocks.NewMockParticipantService(ctrl)
			rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

			psMock.EXPECT().Delete(participantID).Return(assert.AnError)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, http.StatusInternalServerError, rw.Code)
		})
	})
}

func TestPrize(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelDebug))
	require.NoError(t, err)
	require.NotNil(t, router)

	raffleID := "raffle_id_1"
	prizeID := "prize_id_1"
	prizePath := joinPath(ApiPath, RafflesPath, raffleID, PrizesPath)

	t.Run("create_prize", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			prz := &service.PrizeRequest{
				Name:        "super_prize",
				TicketCost:  100_500,
				Description: "bla bla bla",
			}

			encoded, err := json.Marshal(prz)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPost, prizePath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Create(prz).Return(prizeID, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
		})

		t.Run("error", func(t *testing.T) {
			prz := &service.PrizeRequest{
				Name:        "super_prize",
				TicketCost:  100_500,
				Description: "bla bla bla",
			}
			encoded, err := json.Marshal(prz)

			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPost, prizePath, body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Create(prz).Return("", errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})

		t.Run("empty_body", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodPost, prizePath, emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("get_prize", func(t *testing.T) {
		prz := &service.Prize{
			ID:          prizeID,
			Name:        "super_prize",
			TicketCost:  100_500,
			Description: "bla bla bla",
			CreatedAt:   time.Now().UTC(),
		}

		t.Run("success", func(t *testing.T) {
			encoded, err := json.Marshal(prz)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodGet, joinPath(prizePath, prizeID), body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Get(prizeID).Return(prz, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
		})

		t.Run("error", func(t *testing.T) {
			encoded, err := json.Marshal(prz)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodGet, joinPath(prizePath, prizeID), body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Get(prizeID).Return(nil, errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("edit_prize", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			prz := &service.PrizeRequest{
				Name:        "prize_1",
				Description: "description_1",
				TicketCost:  100_500,
			}

			encoded, err := json.Marshal(prz)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPut, joinPath(prizePath, prizeID), body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Edit(prizeID, prz).Return(nil)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, http.StatusOK, rw.Code)
		})

		t.Run("error", func(t *testing.T) {
			upd := &service.PrizeRequest{Name: "prize_1", Description: "description_1", TicketCost: 100_500}

			encoded, err := json.Marshal(upd)
			require.NoError(t, err)

			body := bytes.NewReader(encoded)

			req, err := newRequestWithOrigin(http.MethodPut, joinPath(prizePath, prizeID), body)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Edit(prizeID, upd).Return(errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})

		t.Run("empty_body", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodPut, joinPath(prizePath, prizeID), emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("list_prizes", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
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

			req, err := newRequestWithOrigin(http.MethodGet, prizePath, emptyBody())
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().List().Return(expected, nil)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)
		})

		t.Run("error", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodGet, prizePath, nil)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().List().Return(nil, errors.New("test error"))

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

	t.Run("delete_prize", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodDelete, joinPath(prizePath, prizeID), nil)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Delete(prizeID).Return(nil)

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, http.StatusOK, rw.Code)
		})

		t.Run("error", func(t *testing.T) {
			req, err := newRequestWithOrigin(http.MethodDelete, joinPath(prizePath, prizeID), nil)
			require.NoError(t, err)

			req.Header.Set(GoogleUserIDHeader, organizerID)
			osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

			rsMock := mocks.NewMockRaffleService(ctrl)
			osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

			pzMock := mocks.NewMockPrizeService(ctrl)
			rsMock.EXPECT().PrizeService(raffleID).Return(pzMock)

			pzMock.EXPECT().Delete(prizeID).Return(errors.New("test error"))

			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			require.Equal(t, http.StatusInternalServerError, rw.Code)
		})
	})
}

func TestRecoverMiddleware(t *testing.T) {
	router, err := NewRouter(nil, logger.NewLogger(logger.LevelDebug))
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("panic_recovery", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, "/", nil)
		require.NoError(t, err)

		h := func(rw http.ResponseWriter, r *http.Request) { panic("test panic") }
		rw := httptest.NewRecorder()

		router.recoverMiddleware(http.HandlerFunc(h)).ServeHTTP(rw, req)
		require.Equal(t, http.StatusInternalServerError, rw.Code)
	})

	t.Run("no_panic", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, "/", nil)
		require.NoError(t, err)

		h := func(rw http.ResponseWriter, r *http.Request) {}
		rw := httptest.NewRecorder()

		router.recoverMiddleware(http.HandlerFunc(h)).ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Code)
	})
}

func TestApplyOrganizerMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	router, err := NewRouter(osMock, logger.NewNoOpLogger())
	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		writer := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		assert.True(t, stub.Called())
	})

	t.Run("no_organizer_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		writer := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusBadRequest, writer.Code)
		assert.False(t, stub.Called())
	})

	t.Run("error", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodPost, RafflesPath, nil)
		require.NoError(t, err)

		stub := newHandlerStub()
		handler := http.HandlerFunc(stub.ServeHTTP)

		mockedErr := errors.New("mocked error")

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(mockedErr)

		writer := httptest.NewRecorder()
		router.organizerMiddleware(handler).ServeHTTP(writer, req)
		require.Equal(t, http.StatusInternalServerError, writer.Code)
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

		writer := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		require.Equal(t, defaultOrigin, writer.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("no_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)

		writer := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		require.Equal(t, "", writer.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("wrong_origin", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, RafflesPath, emptyBody())
		require.NoError(t, err)

		req.Header.Set("Origin", "wrong_origin")

		writer := httptest.NewRecorder()
		router.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(writer, req)
		require.Equal(t, http.StatusOK, writer.Code)
		require.Equal(t, "", writer.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestJoinPath(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{"path", "subpath", "subsubpath"}, "/path/subpath/subsubpath"},
		{[]string{"/", "path", "/subpath/", "/subsubpath"}, "/path/subpath/subsubpath"},
		{[]string{"path"}, "/path"},
		{[]string{"/"}, "/"},
	}

	for _, testCase := range testCases {
		result := joinPath(testCase.input...)
		if result != testCase.expected {
			t.Errorf("joinPath(%v) = %v, expected %v", testCase.input, result, testCase.expected)
		}
	}
}

func TestGeParticipantService(t *testing.T) {
	ctrl := gomock.NewController(t)

	organizerID := "organizer_id_1"
	raffleID := "raffle_id_1"

	osMock := mocks.NewMockOrganizerService(ctrl)
	rsMock := mocks.NewMockRaffleService(ctrl)
	psMock := mocks.NewMockParticipantService(ctrl)

	osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil).AnyTimes()
	osMock.EXPECT().RaffleService(organizerID).Return(rsMock).AnyTimes()

	router, err := NewRouter(osMock, logger.NewLogger(logger.LevelDebug))

	require.NoError(t, err)
	require.NotNil(t, router)

	t.Run("success", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodGet, "/api/raffles/raffle_id_1/participants", nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add(raffleIDParam, raffleID)

		rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(ctx)

		ps, err := router.getParticipantService(req)
		assert.NoError(t, err)
		assert.Equal(t, ps, psMock)
	})

	t.Run("missing_organizer_id", func(t *testing.T) {
		req, err := newRequestWithOrigin(http.MethodGet, "/api/raffle/raffle_id_1/participants", nil)
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
		req, err := newRequestWithOrigin(http.MethodGet, "/api/raffles//participants", nil)
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add(raffleIDParam, "")

		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(ctx)

		ps, err := router.getParticipantService(req)
		assert.Nil(t, ps)
		assert.ErrorIs(t, err, ErrMissingID)
	})
}

func joinPath(args ...string) string {
	return path.Clean("/" + path.Join(args...))
}

type HandlerStub struct {
	called bool
	once   sync.Once
}

func (h *HandlerStub) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
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
