//
//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package events_test

import (
	"fmt"
	"testing"

	model "github.com/edgexfoundry/edgex-go/pkg/models"

	"github.com/mteodor/edgex-app/events"

	"github.com/mteodor/edgex-app/events/postgres"
	"github.com/stretchr/testify/assert"
)

func TestEventSave(t *testing.T) {
	id := 0

	cases := []struct {
		desc  string
		Event model.Event
		err   error
	}{
		{"new Event", model.Event{Id, "pass"}, nil},
		{"duplicate Event", model.Event{Id, "pass"}, events.ErrConflict},
	}

	repo := postgres.New(db)

	for _, tc := range cases {
		err := repo.Save(tc.Event)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSingleEventRetrieval(t *testing.T) {
	id := 0

	repo := postgres.New(db)
	repo.Save(model.Event{email, "pass"})

	cases := map[string]struct {
		id  int
		err error
	}{
		"existing Event":     {id, nil},
		"non-existing Event": {-1, events.ErrNotFound},
	}

	for desc, tc := range cases {
		_, err := repo.RetrieveByID(tc.id)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}
