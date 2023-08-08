package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kaznasho/yarmarok/service"
)

func (r *Router) getPrizeService(req *http.Request) (service.PrizeService, error) {
	raffleService, err := r.getRaffleService(req)
	if err != nil {
		return nil, err
	}

	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil {
		return nil, errors.Join(ErrMissingID, err)
	}

	return raffleService.PrizeService(raffleID), nil
}

func (r *Router) getParticipantService(req *http.Request) (service.ParticipantService, error) {
	raffleService, err := r.getRaffleService(req)
	if err != nil {
		return nil, err
	}

	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil {
		return nil, errors.Join(ErrMissingID, err)
	}

	return raffleService.ParticipantService(raffleID), nil
}

func (r *Router) getRaffleService(req *http.Request) (service.RaffleService, error) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		return nil, err
	}

	return r.organizerService.RaffleService(organizerID), nil
}

func extractOrganizerID(r *http.Request) (id string, err error) {
	defer func() {
		if localRun && err != nil {
			err = nil
			id = "dummy_test_user"
		}
	}()

	ids := r.Header.Values(GoogleUserIDHeader)

	if len(ids) != 1 {
		return "", ErrAmbiguousOrganizerIDHeader
	}

	id = ids[0]
	if id == "" {
		return "", ErrAmbiguousOrganizerIDHeader
	}

	return id, nil
}

func extractParam(req *http.Request, param string) (string, error) {
	val := chi.URLParam(req, param)
	if val == "" {
		return "", fmt.Errorf("missing param: %s", param)
	}

	return val, nil
}
