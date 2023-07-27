package web

import (
	"net/http"
)

const (
	ApiPath                  = "/api"
	RafflesPath              = "/raffles"
	ParticipantsPath         = "/participants"
	raffleIDParam            = "raffle_id"
	participantIDParam       = "participant_id"
	raffleIDPlaceholder      = "/{" + raffleIDParam + "}"
	participantIDPlaceholder = "/{" + participantIDParam + "}"
	rafflesGroup             = ApiPath + RafflesPath
	participantsGroup        = rafflesGroup + raffleIDPlaceholder + ParticipantsPath
)

func (w *Web) Routes() {
	w.Handle(http.MethodGet, ApiPath, "/login", w.loginHandler)

	w.Handle(http.MethodPost, rafflesGroup, "/", w.createRaffle)
	w.Handle(http.MethodGet, rafflesGroup, "/", w.listRaffles)
	w.Handle(http.MethodGet, rafflesGroup+raffleIDPlaceholder, "/download-xlsx", w.downloadRaffleXLSX, WithXLSX /*middleware*/)

	w.Handle(http.MethodPost, participantsGroup, "/", w.createParticipant)
	w.Handle(http.MethodGet, participantsGroup, "/", w.listParticipants)
	w.Handle(http.MethodPut, participantsGroup, participantIDPlaceholder, w.updateParticipant)
	w.Handle(http.MethodDelete, participantsGroup, participantIDPlaceholder, w.deleteParticipant)
}

func (w *Web) listRaffles(rw http.ResponseWriter, req *http.Request) error {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		return err
	}

	raffleService := w.svc.RaffleService(organizerID)

	return newM0(raffleService.List).Handle(rw, req)
}

func (w *Web) createRaffle(rw http.ResponseWriter, req *http.Request) error       { return nil }
func (w *Web) loginHandler(rw http.ResponseWriter, req *http.Request) error       { return nil }
func (w *Web) downloadRaffleXLSX(rw http.ResponseWriter, req *http.Request) error { return nil }
func (w *Web) createParticipant(rw http.ResponseWriter, req *http.Request) error  { return nil }
func (w *Web) updateParticipant(rw http.ResponseWriter, req *http.Request) error  { return nil }
func (w *Web) deleteParticipant(rw http.ResponseWriter, req *http.Request) error  { return nil }
func (w *Web) listParticipants(rw http.ResponseWriter, req *http.Request) error   { return nil }
