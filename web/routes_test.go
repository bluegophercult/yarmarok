package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/mocks"
	"github.com/kaznasho/yarmarok/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)

	organizerID := "organizer_id_1"

	log := logger.NewLogger(logger.LevelDebug)

	web := NewWeb(log, osMock)
	web.Routes()
	require.NotNil(t, web)

	t.Run("login_endpoint", func(t *testing.T) {
		loginPath := joinPath(ApiPath, "/login")

		req, err := newRequestWithOrigin(http.MethodPost, loginPath, emptyBody())
		require.NoError(t, err)

		req.Header.Set(GoogleUserIDHeader, organizerID)
		osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

		rw := httptest.NewRecorder()
		web.ServeHTTP(rw, req)
		require.Equal(t, http.StatusSeeOther, rw.Code)
		require.Equal(t, "/", rw.Header().Get("Location"))
	})

	t.Run("raffle_endpoint", func(t *testing.T) {
		rafflePath := joinPath(ApiPath, RafflesPath)

		t.Run("create_raffle", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				initRequest := &service.RaffleInitRequest{
					Name: "raffle_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(initRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, rafflePath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				rsMock.EXPECT().Create(initRequest).Return(&service.CreateResult{ID: "raffle_id_1"}, nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
			})

			t.Run("error", func(t *testing.T) {
				initRequest := &service.RaffleInitRequest{
					Name: "raffle_1",
					Note: "note_1",
				}

				encoded, err := json.Marshal(initRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, rafflePath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				mockedErr := assert.AnError
				rsMock.EXPECT().Create(initRequest).Return(nil, mockedErr)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})
		})

		t.Run("list_raffles", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				dummyTime := time.Now().UTC()
				expected := &service.RaffleListResponse{
					Raffles: []service.Raffle{
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
					},
				}

				req, err := newRequestWithOrigin(http.MethodGet, rafflesGroup, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				rsMock.EXPECT().List().Return(expected, nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)

				assertJSONResponse(t, expected, rw.Body)
			})

			t.Run("error", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodGet, rafflesGroup, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				mockedErr := assert.AnError
				rsMock.EXPECT().List().Return(nil, mockedErr)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})
		})

		t.Run("download_raffle_xlsx", func(t *testing.T) {
			raffleID := "raffle_id_1"
			downloadPath := joinPath(ApiPath, RafflesPath, raffleID, "/download-xlsx")

			t.Run("success", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodGet, downloadPath, nil)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				rsMock.EXPECT().Export(raffleID).Return(
					&service.RaffleExportResponse{
						FileName: "raffle.xlsx",
						Content:  []byte("content")}, nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
				require.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", rw.Header().Get("Content-Type"))
				require.Equal(t, "attachment; filename=raffle.xlsx", rw.Header().Get("Content-Disposition"))
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

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})
		})
	})

	t.Run("participant_endpoint", func(t *testing.T) {
		raffleID := "raffle_id_1"
		participantID := "participant_id_1"
		participantPath := joinPath(ApiPath, RafflesPath, raffleID, ParticipantsPath)

		t.Run("create_participant", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				participantInitRequest := &service.ParticipantRequest{
					Name:  "participant_1",
					Note:  "note_1",
					Phone: "1234567890",
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

				psMock.EXPECT().Create(participantInitRequest).Return(&service.CreateResult{ID: participantID}, nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
			})

			t.Run("error", func(t *testing.T) {
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

				psMock.EXPECT().Create(participantInitRequest).Return(nil, assert.AnError)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
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

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusBadRequest, rw.Code)
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
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
			})

			t.Run("error", func(t *testing.T) {
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

				psMock.EXPECT().Edit(participantID, participantEditRequest).Return(assert.AnError)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
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

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusBadRequest, rw.Code)
			})
		})

		t.Run("list_participants", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				dummyTime := time.Now().UTC()
				expected := &service.ParticipantListResult{
					Participants: []service.Participant{
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
					},
				}

				req, value := newRequestWithOrigin(http.MethodGet, participantPath, emptyBody())
				require.NoError(t, value)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

				psMock.EXPECT().List().Return(expected, nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
			})

			t.Run("error", func(t *testing.T) {
				req, value := newRequestWithOrigin(http.MethodGet, participantPath, nil)
				require.NoError(t, value)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockParticipantService(ctrl)
				rsMock.EXPECT().ParticipantService(raffleID).Return(psMock)

				psMock.EXPECT().List().Return(nil, assert.AnError)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
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
				web.ServeHTTP(rw, req)
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
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})
		})
	})
}

func TestPrize(t *testing.T) {
	ctrl := gomock.NewController(t)

	organizerID := "organizer_id_1"
	raffleID := "raffle_id_1"

	osMock := mocks.NewMockOrganizerService(ctrl)

	log := logger.NewLogger(logger.LevelDebug)

	web := NewWeb(log, osMock)
	require.NotNil(t, web)

	web.Routes()

	t.Run("prize_endpoint", func(t *testing.T) {
		prizePath := joinPath(ApiPath, RafflesPath, raffleID, PrizePath)

		t.Run("create_prize", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				prizeInitRequest := &service.PrizeRequest{
					Name:        "Updated Prize",
					TicketCost:  100_500,
					Description: "Bla bla bla",
				}

				encoded, err := json.Marshal(prizeInitRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, prizePath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				psMock.EXPECT().Create(prizeInitRequest).Return(&service.CreateResult{ID: "prize_id_1"}, nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
			})

			t.Run("error", func(t *testing.T) {
				prizeInitRequest := &service.PrizeRequest{
					Name:        "Updated Prize",
					TicketCost:  100_500,
					Description: "Bla bla bla",
				}

				encoded, err := json.Marshal(prizeInitRequest)
				require.NoError(t, err)

				body := bytes.NewReader(encoded)

				req, err := newRequestWithOrigin(http.MethodPost, prizePath, body)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				mockedErr := assert.AnError
				psMock.EXPECT().Create(prizeInitRequest).Return(nil, mockedErr)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})
		})

		t.Run("list_prizes", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				expected := &service.PrizeListResult{
					Prizes: []service.Prize{
						{
							ID:          "prize_id_1",
							Name:        "Prize 1",
							TicketCost:  10,
							Description: "This is the first prize",
							CreatedAt:   time.Now(),
						},
						{
							ID:          "prize_id_2",
							Name:        "Prize 2",
							TicketCost:  5,
							Description: "This is the second prize",
							CreatedAt:   time.Now(),
						},
						{
							ID:          "prize_id_3",
							Name:        "Prize 3",
							TicketCost:  20,
							Description: "This is the third prize",
							CreatedAt:   time.Now(),
						},
					},
				}

				req, err := newRequestWithOrigin(http.MethodGet, prizePath, emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				psMock.EXPECT().List().Return(expected, nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)

				assertJSONResponse(t, expected, rw.Body)
			})

			t.Run("error", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodGet, prizePath, nil)
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				mockedErr := assert.AnError
				psMock.EXPECT().List().Return(nil, mockedErr)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})
		})

		t.Run("update_prize", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				prizeEditRequest := &service.PrizeRequest{
					Name:        "Updated Prize",
					TicketCost:  100_500,
					Description: "Bla bla bla",
				}

				encoded, value := json.Marshal(prizeEditRequest)
				require.NoError(t, value)

				body := bytes.NewReader(encoded)

				req, value := newRequestWithOrigin(http.MethodPut, joinPath(prizePath, "prize_id_1"), body)
				require.NoError(t, value)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				psMock.EXPECT().Edit("prize_id_1", prizeEditRequest).Return(nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
			})

			t.Run("error", func(t *testing.T) {
				prizeEditRequest := &service.PrizeRequest{
					Name:        "Updated Prize",
					TicketCost:  100_500,
					Description: "Bla bla bla",
				}

				encoded, value := json.Marshal(prizeEditRequest)
				require.NoError(t, value)

				body := bytes.NewReader(encoded)

				req, value := newRequestWithOrigin(http.MethodPut, joinPath(prizePath, "prize_id_1"), body)
				require.NoError(t, value)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				mockedErr := assert.AnError
				psMock.EXPECT().Edit("prize_id_1", prizeEditRequest).Return(mockedErr)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})

			t.Run("empty_body", func(t *testing.T) {
				req, err := newRequestWithOrigin(http.MethodPut, joinPath(prizePath, "prize_id_1"), emptyBody())
				require.NoError(t, err)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusBadRequest, rw.Code)
			})
		})

		t.Run("delete_prize", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				req, value := newRequestWithOrigin(http.MethodDelete, joinPath(prizePath, "prize_id_1"), nil)
				require.NoError(t, value)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				psMock.EXPECT().Delete("prize_id_1").Return(nil)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusOK, rw.Code)
			})

			t.Run("error", func(t *testing.T) {
				req, value := newRequestWithOrigin(http.MethodDelete, joinPath(prizePath, "prize_id_1"), nil)
				require.NoError(t, value)

				req.Header.Set(GoogleUserIDHeader, organizerID)
				osMock.EXPECT().CreateOrganizerIfNotExists(organizerID).Return(nil)

				rsMock := mocks.NewMockRaffleService(ctrl)
				osMock.EXPECT().RaffleService(organizerID).Return(rsMock)

				psMock := mocks.NewMockPrizeService(ctrl)
				rsMock.EXPECT().PrizeService(raffleID).Return(psMock)

				mockedErr := assert.AnError
				psMock.EXPECT().Delete("prize_id_1").Return(mockedErr)

				rw := httptest.NewRecorder()
				web.ServeHTTP(rw, req)
				require.Equal(t, http.StatusInternalServerError, rw.Code)
			})
		})
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

func joinPath(args ...string) string {
	return path.Clean("/" + path.Join(args...))
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
