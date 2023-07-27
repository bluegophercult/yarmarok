package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/web/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeb(t *testing.T) {
	ctrl := gomock.NewController(t)

	osMock := mocks.NewMockOrganizerService(ctrl)
	organizerID := "organizer_id_1"

	log := logger.NewLogger(logger.LevelDebug)

	web := NewWeb(log, osMock)
	web.Routes()
	require.NotNil(t, web)

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

			writer := httptest.NewRecorder()
			web.ServeHTTP(writer, req)
			require.Equal(t, http.StatusOK, writer.Code)

			assertJSONResponse(t, expected, writer.Body)
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

			writer := httptest.NewRecorder()
			web.ServeHTTP(writer, req)
			require.Equal(t, http.StatusInternalServerError, writer.Code)
		})
	})

}
