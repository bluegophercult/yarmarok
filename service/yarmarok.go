package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// stringUUID is a plumbing function for generating UUIDs.
// It is overridden in tests.
var stringUUID = func() string {
	return uuid.New().String()
}

// timeNow is a plumbing function for getting the current time.
// It is overridden in tests.
var timeNow = func() time.Time {
	return time.Now()
}

// Yarmarok represents a yarmarok.
type Yarmarok struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
	Note      string
}

// YarmarokService is a service for yarmaroks.
type YarmarokService interface {
	Init(*YarmarokInitRequest) (*InitResult, error)
	Get(id string) (*Yarmarok, error)
	List() (*YarmarokListResponse, error)
}

// YarmarokStorage is a storage for yarmaroks.
type YarmarokStorage interface {
	Create(*Yarmarok) error
	Get(id string) (*Yarmarok, error)
	GetAll() ([]Yarmarok, error)
}

// YarmarokManager is an implementation of YarmarokService.
type YarmarokManager struct {
	yarmarokStorage YarmarokStorage
}

// NewYarmarokManager creates a new YarmarokManager.
func NewYarmarokManager(ys YarmarokStorage) *YarmarokManager {
	return &YarmarokManager{
		yarmarokStorage: ys,
	}
}

// Init initializes a yarmarok.
func (ym *YarmarokManager) Init(y *YarmarokInitRequest) (*InitResult, error) {
	yarmarok := Yarmarok{
		ID:        stringUUID(),
		Name:      y.Name,
		Note:      y.Note,
		CreatedAt: timeNow(),
	}

	err := ym.yarmarokStorage.Create(&yarmarok)
	if err != nil {
		return nil, err
	}

	return &InitResult{
		ID: yarmarok.ID,
	}, nil
}

// Get returns a yarmarok by id.
func (ym *YarmarokManager) Get(id string) (*Yarmarok, error) {
	return ym.yarmarokStorage.Get(id)
}

// List lists yarmaroks in user's scope.
func (ym *YarmarokManager) List() (*YarmarokListResponse, error) {
	yarmaroks, err := ym.yarmarokStorage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all yarmaroks: %w", err)
	}

	return &YarmarokListResponse{
		Yarmaroks: yarmaroks,
	}, nil
}

// YarmarokInitRequest is a request for initializing a yarmarok.
type YarmarokInitRequest struct {
	Name string
	Note string
}

// InitResult is a generic result of entity initialization.
type InitResult struct {
	ID string
}

// YarmarokListResponse is a response for listing yarmaroks.
type YarmarokListResponse struct {
	Yarmaroks []Yarmarok
}
