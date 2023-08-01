package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kaznasho/yarmarok/service"
)

const (
	/*	ApiPath                  = "/api"
		RafflesPath              = "/raffles"
		ParticipantsPath         = "/participants"
		raffleIDParam            = "raffle_id"
		participantIDParam       = "participant_id"
		raffleIDPlaceholder      = "/{" + raffleIDParam + "}" */
	participantIDPlaceholder = "/{" + participantIDParam + "}"
	rafflesGroup             = ApiPath + RafflesPath
	participantsGroup        = rafflesGroup + raffleIDPlaceholder + ParticipantsPath
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
}

func (w *Web) loginHandler(rw http.ResponseWriter, req *http.Request) error {
	http.RedirectHandler("/", http.StatusSeeOther).ServeHTTP(rw, req)
	return nil
}

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

	if _, err := rw.Write(resp.Content); err != nil {
		return fmt.Errorf("writing xlsx: %w", err)
	}

	return nil
}

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

/* TOOLING */

func (w *Web) getParticipantService(req *http.Request) (service.ParticipantService, error) {
	raffleService, err := w.getRaffleService(req)
	if err != nil {
		return nil, err
	}

	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil || raffleID == "" {
		return nil, NewError(ErrMissingID, http.StatusBadRequest)
	}

	return raffleService.ParticipantService(raffleID), nil
}

func (w *Web) getRaffleService(req *http.Request) (service.RaffleService, error) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		return nil, err
	}

	return w.svc.RaffleService(organizerID), nil
}
