package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/kaznasho/yarmarok/service"
)

var (
	// ErrAmbiguousOrganizerIDHeader is returned when the organizer id header is not set or is ambiguous.
	ErrAmbiguousOrganizerIDHeader = errors.New("ambiguous organizer id format")

	// ErrMissingID is returned when id is missing.
	ErrMissingID = errors.New("missing id")
)

// localRun is true if app is build for local run
var localRun = false

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

func (r *Router) getDonationService(req *http.Request) (service.DonationService, error) {
	prizeService, err := r.getPrizeService(req)
	if err != nil {
		return nil, err
	}

	prizeID, err := extractParam(req, prizeIDParam)
	if err != nil {
		return nil, errors.Join(ErrMissingID, err)
	}

	return prizeService.DonationService(prizeID), nil
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
