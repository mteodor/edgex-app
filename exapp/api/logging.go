package api

import (
	"fmt"
	"time"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	log "github.com/mainflux/mainflux/logger"
	"github.com/mteodor/edgex-app/exapp"
)

var _ exapp.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	svc exapp.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc exapp.Service) exapp.Service {
	return &loggingMiddleware{svc}
}
func (lm *loggingMiddleware) Logger() log.Logger {
	return lm.svc.Logger()
}

func (lm *loggingMiddleware) RegisterEvent(ev model.Event) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register event %s took %s to complete", ev.ID, time.Since(begin))
		if err != nil {
			lm.Logger().Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.Logger().Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RegisterEvent(ev)
}

func (lm *loggingMiddleware) RetrieveByID(ID string) (ev model.Event, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method retrieve event %s took %s to complete", ID, time.Since(begin))
		if err != nil {
			lm.Logger().Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.Logger().Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RetrieveByID(ID)
}
