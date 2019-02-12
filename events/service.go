//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package events

import (
	model "github.com/edgexfoundry/edgex-go/pkg/models"
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Register adds new evetn
	RegisterEvent(model.Event) error
<<<<<<< HEAD
	RetrieveByID(string) (model.Event, error)
=======
>>>>>>> 8b2d9977cbbd8ff37891a6bc6f99d4ffa1abe5d9
}

var _ Service = (*eventsService)(nil)

type eventsService struct {
	events EventsRepository
	//	hasher Hasher
	//	idp    IdentityProvider
}

// New instantiates the events service implementation.
func New(events EventsRepository) Service {
	return &eventsService{events: events}
}

func (svc eventsService) RegisterEvent(event model.Event) error {

	return svc.events.Save(event)
}
<<<<<<< HEAD

func (svc eventsService) RetrieveByID(id string) (model.Event, error) {
	return svc.events.RetrieveByID(id)

}
=======
>>>>>>> 8b2d9977cbbd8ff37891a6bc6f99d4ffa1abe5d9
