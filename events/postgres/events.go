//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package postgres

import (
	"database/sql"
	"fmt"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/lib/pq"
	events "github.com/mteodor/edgex-app/events"
)

var _ events.EventsRepository = (*eventRepository)(nil)

const errDuplicate = "unique_violation"

type eventRepository struct {
	db *sql.DB
}

// New instantiates a PostgreSQL implementation of user
// repository.
func New(db *sql.DB) events.EventsRepository {
	return &eventRepository{db}
}

func (rr eventRepository) Save(event model.Event) error {

	fmt.Println("saving event")
	q := `INSERT INTO events (id, pushed, device, created, modifed, origin, event ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	if _, err := rr.db.Exec(q, event.ID, event.Pushed, event.Created, event.Origin, event.Modified, event.Device, event.Origin, event.Event); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && errDuplicate == pqErr.Code.Name() {
			return events.ErrConflict
		}
		return err
	}

	return nil
}

func (rr eventRepository) RetrieveByID(Id string) (model.Event, error) {
	q := `SELECT id, pushed, device, created, modified, origin, event id FROM events WHERE id = $1`

	event := model.Event{}
	if err := rr.db.QueryRow(q, Id).Scan(&event.ID, &event.Pushed, &event.Device, &event.Created, &event.Modified, &event.Origin, &event.Event); err != nil {
		if err == sql.ErrNoRows {
			return event, events.ErrNotFound
		}
		return event, err
	}

	event.ID = Id

	return event, nil
}
