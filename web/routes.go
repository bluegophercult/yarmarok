package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kaznasho/yarmarok/service"
)

const (
	ApiPath                  = "/api"
	RafflesPath              = "/raffles"
	ParticipantsPath         = "/participants"
	PrizePath                = "/prizes"
	raffleIDParam            = "raffle_id"
	participantIDParam       = "participant_id"
	raffleIDPlaceholder      = "/{" + raffleIDParam + "}"
	participantIDPlaceholder = "/{" + participantIDParam + "}"
	prizeIDPlaceholder       = "/{" + raffleIDParam + "}"
	rafflesGroup             = ApiPath + RafflesPath
	participantsGroup        = rafflesGroup + raffleIDPlaceholder + ParticipantsPath
	prizeGroup               = rafflesGroup + raffleIDPlaceholder + PrizePath
)

// localRun is true if app is build for local run
var localRun = false

var (
	// ErrAmbiguousOrganizerIDHeader is returned when organizer id header is not set or is ambiguous.
	ErrAmbiguousOrganizerIDHeader = errors.New("ambiguous organizer id format")
	// ErrMissingID is returned when id is missing.
	ErrMissingID = errors.New("missing id")
)

func (w *Web) Routes() {
	w.Handle(http.MethodPost, ApiPath, "/login", w.loginHandler)

	w.Handle(http.MethodPost, rafflesGroup, "/", w.createRaffle)
	w.Handle(http.MethodGet, rafflesGroup, "/", w.listRaffles)
	w.Handle(http.MethodGet, rafflesGroup+raffleIDPlaceholder, "/download-xlsx", w.downloadRaffleXLSX, WithXLSX /*middleware*/)

	w.Handle(http.MethodPost, participantsGroup, "/", w.createParticipant)
	w.Handle(http.MethodGet, participantsGroup, "/", w.listParticipants)
	w.Handle(http.MethodPut, participantsGroup, participantIDPlaceholder, w.updateParticipant)
	w.Handle(http.MethodDelete, participantsGroup, participantIDPlaceholder, w.deleteParticipant)

	w.Handle(http.MethodPost, prizeGroup, "/", w.createPrize)
	w.Handle(http.MethodGet, prizeGroup, "/", w.listPrizes)
	w.Handle(http.MethodPut, prizeGroup, prizeIDPlaceholder, w.updatePrize)
	w.Handle(http.MethodDelete, prizeGroup, prizeIDPlaceholder, w.deletePrize)
}

func (w *Web) loginHandler(rw http.ResponseWriter, req *http.Request) error {
	http.RedirectHandler("/", http.StatusSeeOther).ServeHTTP(rw, req)
	return nil
}

/*Raffle */

func (w *Web) createRaffle(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getRaffleService(req)
	if err != nil {
		return err
	}

	return newCreate(svc.Create).Handle(rw, req)
}

func (w *Web) listRaffles(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getRaffleService(req)
	if err != nil {
		return err
	}

	return newList(svc.List).Handle(rw, req)
}

func (w *Web) downloadRaffleXLSX(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getRaffleService(req)
	if err != nil {
		return err
	}

	id, err := extractParam(req, raffleIDParam)
	if err != nil {
		return errors.Join(err, ErrMissingID)
	}

	resp, err := svc.Export(id)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Disposition", "attachment; filename="+resp.FileName)

	if err := Respond(rw, resp.Content); err != nil {
		return fmt.Errorf("writing xlsx: %w", err)
	}

	return nil
}

/* Participant */

func (w *Web) createParticipant(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getParticipantService(req)
	if err != nil {
		return err
	}

	return newCreate(svc.Create).Handle(rw, req)
}

func (w *Web) updateParticipant(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getParticipantService(req)
	if err != nil {
		return err
	}

	return newUpdate(svc.Edit).Handle(rw, req)
}

func (w *Web) deleteParticipant(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getParticipantService(req)
	if err != nil {
		return err
	}

	return newDelete(svc.Delete).Handle(rw, req)
}

func (w *Web) listParticipants(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getParticipantService(req)
	if err != nil {
		return err
	}

	return newList(svc.List).Handle(rw, req)
}

/* Prize */

func (w *Web) createPrize(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getPrizeService(req)
	if err != nil {
		return err
	}
	return newCreate(svc.Create).Handle(rw, req)
}

func (w *Web) updatePrize(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getPrizeService(req)
	if err != nil {
		return err
	}
	return newUpdate(svc.Edit).Handle(rw, req)
}

func (w *Web) deletePrize(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getPrizeService(req)
	if err != nil {
		return err
	}
	return newDelete(svc.Delete).Handle(rw, req)
}
func (w *Web) listPrizes(rw http.ResponseWriter, req *http.Request) error {
	svc, err := w.getPrizeService(req)
	if err != nil {
		return err
	}
	return newList(svc.List).Handle(rw, req)
}

/* Tooling */

func (w *Web) getPrizeService(req *http.Request) (service.PrizeService, error) {
	raffleService, err := w.getRaffleService(req)
	if err != nil {
		return nil, err
	}

	raffleID, err := extractRaffleID(req)
	if err != nil {
		return nil, err
	}

	return raffleService.PrizeService(raffleID), nil
}

func (w *Web) getParticipantService(req *http.Request) (service.ParticipantService, error) {
	raffleService, err := w.getRaffleService(req)
	if err != nil {
		return nil, err
	}

	raffleID, err := extractRaffleID(req)
	if err != nil {
		return nil, err
	}

	return raffleService.ParticipantService(raffleID), nil
}

func (w *Web) getRaffleService(req *http.Request) (service.RaffleService, error) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		return nil, NewError(err, http.StatusBadRequest)
	}

	return w.svc.RaffleService(organizerID), nil
}

func extractRaffleID(req *http.Request) (id string, err error) {
	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil || raffleID == "" {
		return "", NewError(ErrMissingID, http.StatusBadRequest)
	}
	return raffleID, nil
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
	if param == "" {
		return "", fmt.Errorf("missing param: %s", param)
	}

	return val, nil
}
