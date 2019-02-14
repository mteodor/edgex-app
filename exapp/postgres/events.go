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

	"github.com/mainflux/mainflux/logger"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/lib/pq"
	events "github.com/mteodor/edgex-app/exapp/events"
)

var _ events.EventsRepository = (*eventRepository)(nil)

const errDuplicate = "unique_violation"

type eventRepository struct {
	logger logger.Logger
	db     *sql.DB
}

// New instantiates a PostgreSQL implementation of user
// repository.
func New(db *sql.DB, l logger.Logger) events.EventsRepository {
	return &eventRepository{l, db}
}

func (rr eventRepository) Save(event model.Event) error {

	rr.logger.Info("saving event")

	q := `INSERT INTO events (id, pushed, device, created, modified, origin, event ) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if _, err := rr.db.Exec(q, event.ID, event.Pushed, event.Device, event.Created, event.Modified, event.Origin, event.Event); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && errDuplicate == pqErr.Code.Name() {
			return events.ErrConflict
		}
		return err
	}

	q = `INSERT INTO readings (eid, rid, pushed,  created, origin, modified, device, value  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for _, r := range event.Readings {

		if _, err := rr.db.Exec(q, event.ID, r.Id, r.Pushed, r.Created, r.Origin, r.Modified, r.Device, r.Name, r.Value); err != nil {
			if pqErr, ok := err.(*pq.Error); ok && errDuplicate == pqErr.Code.Name() {
				return events.ErrConflict
			}
			return err
		}
	}

	return nil
}

func (rr eventRepository) RetrieveByID(id string) (model.Event, error) {
	q := `SELECT id, pushed, device, created, modified, origin, event id FROM events WHERE id = $1`

	event := model.Event{}
	if err := rr.db.QueryRow(q, id).Scan(&event.ID, &event.Pushed, &event.Device, &event.Created, &event.Modified, &event.Origin, &event.Event); err != nil {
		if err == sql.ErrNoRows {
			return event, events.ErrNotFound
		}
		return event, err
	}

	q = `SELECT rid, pushed, created , origin, modified, device, name, value from readings where eid = $1`

	items := []model.Reading{}

	rows, err := rr.db.Query(q, id)
	if err != nil {
		rr.logger.Error(fmt.Sprintf("Failed to retrieve things due to %s", err))
		return event, err
	}
	defer rows.Close()

	for rows.Next() {
		r := model.Reading{}
		if err = rows.Scan(&r.Id, &r.Pushed, &r.Created, &r.Origin, &r.Modified, &r.Device, &r.Name, &r.Value); err != nil {

			rr.logger.Error(fmt.Sprintf("Failed to read retrieved thing due to %s", err))
		}
		items = append(items, r)
	}

	event.ID = id

	return event, nil
}
