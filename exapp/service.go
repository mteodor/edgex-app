//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package exapp

import (
	model "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/mainflux/mainflux/logger"
	"github.com/mteodor/edgex-app/exapp/events"
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Register adds new evetn
	RegisterEvent(model.Event) error

	RetrieveByID(string) (model.Event, error)

	Logger() logger.Logger
}

var _ Service = (*eventsService)(nil)

type eventsService struct {
	events events.EventsRepository
	logger logger.Logger
	//	hasher Hasher
	//	idp    IdentityProvider
}

// New instantiates the events service implementation.
func New(events events.EventsRepository, l logger.Logger) Service {
	return &eventsService{events, l}
}

// GetLogger - retrieves logger that can be used
func (svc eventsService) Logger() logger.Logger {
	return svc.logger
}

func (svc eventsService) RegisterEvent(event model.Event) error {

	return svc.events.Save(event)
}

func (svc eventsService) RetrieveByID(id string) (model.Event, error) {
	return svc.events.RetrieveByID(id)

}
