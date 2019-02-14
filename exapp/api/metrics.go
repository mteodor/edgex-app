package api

import (
	"time"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/go-kit/kit/metrics"
	log "github.com/mainflux/mainflux/logger"
	"github.com/mteodor/edgex-app/exapp"
)

var _ exapp.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     exapp.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc exapp.Service, counter metrics.Counter, latency metrics.Histogram) exapp.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) Logger() log.Logger {
	return ms.svc.Logger()
}

func (ms *metricsMiddleware) RegisterEvent(ev model.Event) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "register_event").Add(1)
		ms.latency.With("method", "register_event").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RegisterEvent(ev)
}

func (ms *metricsMiddleware) RetrieveByID(ID string) (model.Event, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "retrieve_event").Add(1)
		ms.latency.With("method", "retrieve_event").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RetrieveByID(ID)
}
