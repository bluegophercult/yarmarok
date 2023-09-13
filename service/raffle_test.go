package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRaffle(t *testing.T) {
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true") // To test in Windows remove before commit

	ctrl := gomock.NewController(t)

	rsMock := NewMockRaffleStorage(ctrl)
	rm := NewRaffleManager(rsMock)

	raf := RaffleRequest{
		Name: "raffle_name_1",
		Note: "raffle_note_1",
	}

	mockedRaffleID := "raffle_id_1"
	mockedTime := time.Now().UTC()
	mockedErr := assert.AnError

	mockedRaffle := Raffle{
		ID:        mockedRaffleID,
		Name:      raf.Name,
		Note:      raf.Note,
		CreatedAt: mockedTime,
	}

	t.Run("create", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			rsMock.EXPECT().Create(gomock.Any()).Return(mockedErr)

			res, err := rm.Create(&raf)
			require.ErrorIs(t, err, mockedErr)
			require.Empty(t, res)
		})

		t.Run("success", func(t *testing.T) {
			setUUIDMock(mockedRaffleID)
			setTimeNowMock(mockedTime)

			rsMock.EXPECT().Create(&mockedRaffle).Return(nil)

			resID, err := rm.Create(&raf)
			require.NoError(t, err)
			require.Equal(t, mockedRaffleID, resID)
		})
	})

	t.Run("get", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			rsMock.EXPECT().Get(mockedRaffleID).Return(nil, mockedErr)

			res, err := rm.Get(mockedRaffleID)
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			rsMock.EXPECT().Get(mockedRaffleID).Return(&mockedRaffle, nil)

			raf, err := rm.Get(mockedRaffleID)
			require.NoError(t, err)
			require.Equal(t, &mockedRaffle, raf)
		})
	})

	t.Run("edit", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			rsMock.EXPECT().Get(mockedRaffleID).Return(&mockedRaffle, nil)
			rsMock.EXPECT().Update(&mockedRaffle).Return(nil)
			err := rm.Edit(mockedRaffleID, &raf)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			rsMock.EXPECT().Get(mockedRaffleID).Return(nil, ErrNotFound)
			err := rm.Edit(mockedRaffleID, &raf)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			rsMock.EXPECT().Delete(mockedRaffleID).Return(nil)
			err := rm.Delete(mockedRaffleID)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			rsMock.EXPECT().Delete(mockedRaffleID).Return(ErrNotFound)
			err := rm.Delete(mockedRaffleID)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("error", func(t *testing.T) {
			rsMock.EXPECT().GetAll().Return(nil, mockedErr)

			res, err := rm.List()
			require.ErrorIs(t, err, mockedErr)
			require.Nil(t, res)
		})

		t.Run("success", func(t *testing.T) {
			raffles := []Raffle{mockedRaffle, mockedRaffle, mockedRaffle}

			rsMock.EXPECT().GetAll().Return(raffles, nil)

			res, err := rm.List()
			require.NoError(t, err)
			require.Equal(t, raffles, res)
		})
	})

	t.Run("Export non-empty collection s", func(t *testing.T) {
		raf := &Raffle{ID: mockedRaffleID, Name: "Raffle Test"}
		prts := []Participant{
			{ID: "p1", Name: "Participant 1"},
			{ID: "p2", Name: "Participant 2"},
		}
		przs := []Prize{
			{ID: "pr1", Name: "Prize 1"},
			{ID: "pr2", Name: "Prize 2"},
		}

		rsMock.EXPECT().Get(mockedRaffleID).Return(raf, nil)

		psMock := NewMockParticipantStorage(ctrl)
		rsMock.EXPECT().ParticipantStorage(mockedRaffleID).Return(psMock)
		psMock.EXPECT().GetAll().Return(prts, nil)

		pzMock := NewMockPrizeStorage(ctrl)
		rsMock.EXPECT().PrizeStorage(mockedRaffleID).Return(pzMock)
		pzMock.EXPECT().GetAll().Return(przs, nil)

		res, err := rm.Export(mockedRaffleID)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "yarmarok_"+mockedRaffleID+".xlsx", res.FileName)
		require.NotEmpty(t, res.Content)
	})

	t.Run("PlayPrize", func(t *testing.T) {
		prts := []Participant{
			{ID: "p1", Name: "Participant 1"},
			{ID: "p2", Name: "Participant 2"},
			{ID: "p3", Name: "Participant 3"},
		}
		przs := []Prize{
			{ID: "pr1", Name: "Prize 1", TicketCost: 10},
			{ID: "pr2", Name: "Prize 2", TicketCost: 20},
		}

		dnt := []Donation{
			{ID: "dn1", PrizeID: "pr1", ParticipantID: "p1", Amount: 100},
			{ID: "dn1", PrizeID: "pr1", ParticipantID: "p1", Amount: 100},
			{ID: "dn1", PrizeID: "pr1", ParticipantID: "p2", Amount: 200},
			{ID: "dn1", PrizeID: "pr1", ParticipantID: "p2", Amount: 200},
			{ID: "dn1", PrizeID: "pr1", ParticipantID: "p3", Amount: 300},
		}

		mockedPrizeID := "pz1"
		mockedDonation := "dn1"

		psMock := NewMockParticipantStorage(ctrl)
		rsMock.EXPECT().ParticipantStorage(mockedRaffleID).Return(psMock)
		psMock.EXPECT().GetAll().Return(prts, nil)

		pzMock := NewMockPrizeStorage(ctrl)
		rsMock.EXPECT().PrizeStorage(mockedRaffleID).Return(pzMock)
		pzMock.EXPECT().Get(mockedPrizeID).Return(&przs[0], nil)

		dnMock := NewMockDonationStorage(ctrl)
		pzMock.EXPECT().DonationStorage(mockedPrizeID).Return(dnMock)
		dnMock.EXPECT().GetAll().Return(dnt, nil)

		dnMock.EXPECT().Get(mockedDonation).Return(&dnt[0], nil)

		res, err := rm.PlayPrize(mockedRaffleID, mockedPrizeID)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEmpty(t, res.Winners)
		require.NotEmpty(t, res.PlayParticipants)
	})

	t.Run("PlayPrizeAgain", func(t *testing.T) {
		przs := []Prize{
			{ID: "pr1", Name: "Prize 1", TicketCost: 10},
			{ID: "pr2", Name: "Prize 2", TicketCost: 20},
		}

		mockedPreviousResult := &PrizePlayResult{
			Winners: []PlayParticipant{
				{
					Participant: Participant{
						ID:        "ID1",
						Name:      "name1",
						Phone:     "phone1",
						Note:      "note1",
						CreatedAt: mockedTime,
					},
					TotalDonation:      300,
					TotalTicketsNumber: 10,
					Donations: []Donation{
						{
							ID:            "dID1",
							PrizeID:       "prID1",
							ParticipantID: "id1",
							Amount:        300,
							TicketsNumber: 10,
							CreatedAt:     time.Time{},
						},
					},
				},
			},
			PlayParticipants: []PlayParticipant{
				{
					Participant: Participant{
						ID:        "ID2",
						Name:      "name2",
						Phone:     "phone2",
						Note:      "note2",
						CreatedAt: mockedTime,
					},
					TotalDonation:      200,
					TotalTicketsNumber: 5,
					Donations: []Donation{
						{
							ID:            "dID2",
							PrizeID:       "prID1",
							ParticipantID: "ID2",
							Amount:        200,
							TicketsNumber: 5,
							CreatedAt:     mockedTime,
						},
					},
				},
				{
					Participant: Participant{
						ID:        "ID3",
						Name:      "name3",
						Phone:     "phone3",
						Note:      "note3",
						CreatedAt: mockedTime,
					},
					TotalDonation:      100,
					TotalTicketsNumber: 2,
					Donations: []Donation{
						{
							ID:            "dID3",
							PrizeID:       "prID1",
							ParticipantID: "ID3",
							Amount:        1000,
							TicketsNumber: 2,
							CreatedAt:     mockedTime,
						},
					},
				},
			},
		}

		mockedPrizeID := "pz1"

		pzMock := NewMockPrizeStorage(ctrl)
		rsMock.EXPECT().PrizeStorage(mockedRaffleID).Return(pzMock)
		pzMock.EXPECT().Get(mockedPrizeID).Return(&przs[0], nil)

		res, err := rm.PlayPrizeAgain(mockedRaffleID, mockedPrizeID, mockedPreviousResult)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEmpty(t, res.Winners)
		require.NotEmpty(t, res.PlayParticipants)
	})
}

func setUUIDMock(uuid string) {
	stringUUID = func() string {
		return uuid
	}
}

func setTimeNowMock(t time.Time) {
	timeNow = func() time.Time {
		return t
	}
}
