package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	psMock := NewMockParticipantStorage(ctrl)
	pm := NewParticipantManager(psMock)

	id := "participant_id"
	p := &ParticipantRequest{
		Name:  "John Doe",
		Phone: "1234567890",
		Note:  "Test participant",
	}

	t.Run("add", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			psMock.EXPECT().Create(gomock.Any()).Return(nil)
			_, err := pm.Create(p)
			require.NoError(t, err)
		})

		t.Run("already_exists", func(t *testing.T) {
			psMock.EXPECT().Create(gomock.Any()).Return(ErrAlreadyExists)
			_, err := pm.Create(p)
			require.ErrorIs(t, err, ErrAlreadyExists)
		})
	})

	t.Run("edit", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			psMock.EXPECT().Get(gomock.Any()).Return(&Participant{}, nil)
			psMock.EXPECT().Update(gomock.Any()).Return(nil)
			err := pm.Edit(id, p)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			psMock.EXPECT().Get(gomock.Any()).Return(nil, ErrNotFound)
			err := pm.Edit(id, p)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			psMock.EXPECT().Delete(gomock.Any()).Return(nil)
			err := pm.Delete(id)
			require.NoError(t, err)
		})

		t.Run("not_found", func(t *testing.T) {
			psMock.EXPECT().Delete(gomock.Any()).Return(ErrNotFound)
			err := pm.Delete(id)
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			psMock.EXPECT().GetAll().Return([]Participant{}, nil)
			_, err := pm.List()
			require.NoError(t, err)
		})

		t.Run("error", func(t *testing.T) {
			psMock.EXPECT().GetAll().Return(nil, errors.New("test error"))
			_, err := pm.List()
			require.Error(t, err)
		})
	})
}
