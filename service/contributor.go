package service

import (
	"errors"
	"time"
)

var (
	ErrContributorAlreadyExists = errors.New("contributor already exists")
	ErrContributorNotFound      = errors.New("contributor not found")
)

// Contributor represents a contributor of the application.
type Contributor struct {
	ID        string
	Name      string
	Phone     string
	Note      string
	CreatedAt time.Time
}

// ContributorAddRequest is a request for creating a new contributor.
type ContributorAddRequest struct {
	Name  string
	Phone string
	Note  string
}

// ContributorEditRequest is a request for updating a contributor.
type ContributorEditRequest Contributor

// ContributorListResult is a response for listing contributors.
type ContributorListResult struct {
	Contributors []Contributor
}

// ContributorService is a service for contributors.
type ContributorService interface {
	Add(*ContributorAddRequest) (*InitResult, error)
	Edit(*ContributorEditRequest) (*Result, error)
	List() (*ContributorListResult, error)
}

// ContributorStorage is a storage for contributors.
type ContributorStorage interface {
	Create(*Contributor) error
	Get(id string) (*Contributor, error)
	Update(*Contributor) error
	GetAll() ([]Contributor, error)
	Delete(id string) error
}

// ContributorManager is an implementation of ContributorService.
type ContributorManager struct {
	contributorStorage ContributorStorage
}

// NewContributorManager creates a new ContributorManager.
func NewContributorManager(s ContributorStorage) *ContributorManager {
	return &ContributorManager{contributorStorage: s}
}

// Add creates a new contributor
func (m *ContributorManager) Add(ctb *ContributorAddRequest) (*InitResult, error) {
	contributor := toContributor(ctb)
	if err := m.contributorStorage.Create(contributor); err != nil {
		return nil, err
	}

	return &InitResult{ID: contributor.ID}, nil
}

// Edit updates a contributor
func (m *ContributorManager) Edit(ctb *ContributorEditRequest) (*Result, error) {
	contributor, err := m.contributorStorage.Get(ctb.ID)
	if err != nil {
		return &Result{StatusError}, err
	}

	if err := m.contributorStorage.Update(contributor); err != nil {
		return &Result{StatusError}, err
	}

	return &Result{StatusSuccess}, nil
}

// List returns all contributors.
func (m *ContributorManager) List() (*ContributorListResult, error) {
	contributors, err := m.contributorStorage.GetAll()
	if err != nil {
		return nil, err
	}

	return &ContributorListResult{Contributors: contributors}, nil
}

func toContributor(ctb *ContributorAddRequest) *Contributor {
	return &Contributor{
		ID:        stringUUID(),
		Name:      ctb.Name,
		Phone:     ctb.Phone,
		Note:      ctb.Note,
		CreatedAt: timeNow(),
	}
}
