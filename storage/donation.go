package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"

	"github.com/kaznasho/yarmarok/service"
)

type FirestoreDonationStorage struct {
	participantID   string
	firestoreClient *firestore.CollectionRef
}

func (rs *FirestoreRaffleStorage) DonationStorage(participantID string) service.DonationStorage {
	return NewFirestoreDonationStorage(rs.firestoreClient.Doc(participantID).Collection(participantCollection), participantID)
}

func NewFirestoreDonationStorage(client *firestore.CollectionRef, participantID string) *FirestoreDonationStorage {
	return &FirestoreDonationStorage{
		participantID:   participantID,
		firestoreClient: client,
	}
}

func (ds *FirestoreDonationStorage) Create(participantStorage service.ParticipantStorage, prizeStorage service.PrizeStorage, d *service.Donation) error {
	exists, err := ds.Exists(d.ID)
	if err != nil {
		return fmt.Errorf("check donation exists: %w", err)
	}

	if exists {
		return service.ErrDonationAlreadyExists
	}

	participant, err := participantStorage.Get(d.ParticipantID)
	if err != nil {
		return fmt.Errorf("check existence: %w", err)
	}
	d.ParticipantID = participant.ID

	prize, err := prizeStorage.Get(d.PrizeID)
	if err != nil {
		return fmt.Errorf("check existence: %w", err)
	}
	d.PrizeID = prize.ID

	d.TicketNumber = d.Amount / prize.TicketCost

	_, err = ds.firestoreClient.Doc(d.ID).Set(context.Background(), d)

	return err
}

func (ds *FirestoreDonationStorage) Get(id string) (*service.Donation, error) {
	doc, err := ds.firestoreClient.Doc(id).Get(context.Background())
	if err != nil {
		if isNotFound(err) {
			return nil, service.ErrDonationNotFound
		}
		return nil, fmt.Errorf("get donation: %w", err)
	}

	var d service.Donation
	if err := doc.DataTo(&d); err != nil {
		return nil, fmt.Errorf("decode donation: %w", err)
	}

	return &d, nil
}

func (ds *FirestoreDonationStorage) GetAll() ([]service.Donation, error) {
	docs, err := ds.firestoreClient.Documents(context.Background()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all donations: %w", err)
	}

	donations := make([]service.Donation, 0, len(docs))
	for _, doc := range docs {
		var d service.Donation
		if err = doc.DataTo(&d); err != nil {
			return nil, fmt.Errorf("decoce donations: %w", err)
		}

		donations = append(donations, d)
	}

	return donations, nil
}

func (ds *FirestoreDonationStorage) Update(d *service.Donation) error {
	exists, err := ds.Exists(d.ID)
	if err != nil {
		return fmt.Errorf("check donation exists: %w", err)
	}

	if !exists {
		return service.ErrDonationNotFound
	}

	if _, err := ds.firestoreClient.Doc(d.ID).Set(context.Background(), d); err != nil {
		return fmt.Errorf("update donation: %w", err)
	}

	return nil
}

func (ds *FirestoreDonationStorage) Delete(id string) error {
	exists, err := ds.Exists(id)
	if err != nil {
		return fmt.Errorf("check donation exists: %w", err)
	}

	if !exists {
		return service.ErrDonationNotFound
	}

	if _, err := ds.firestoreClient.Doc(id).Delete(context.Background()); err != nil {
		return fmt.Errorf("delete donation: %w", err)
	}

	return nil
}

func (ds *FirestoreDonationStorage) Exists(id string) (bool, error) {
	doc, err := ds.firestoreClient.Doc(id).Get(context.Background())
	if isNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return doc.Exists(), nil
}
