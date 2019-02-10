package events

import (
	"errors"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
)

var (
	// ErrConflict indicates usage of the existing event
	ErrConflict = errors.New("event already registred")
	ErrNotFound = errors.New("event not found")
)

// eventsRepository specifies an account persistence API.
type EventsRepository interface {
	// Save persists the events. A non-nil error is returned to indicate
	// operation failure.
	Save(model.Event) error

	// RetrieveByID retrieves event by its unique identifier (i.e. email).
	RetrieveByID(string) (model.Event, error)
}
