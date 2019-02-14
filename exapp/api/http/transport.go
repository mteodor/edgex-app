package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/things"
	"github.com/mteodor/edgex-app/exapp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType = "application/json"
	offset      = "offset"
	limit       = "limit"

	defOffset = 0
	defLimit  = 10
)

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	errInvalidQueryParams     = errors.New("invalid query params")
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc exapp.Service, logger logger.Logger) http.Handler {
	if svc == nil {
		return nil
	}
	logger.Info("making handler")
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := bone.New()

	r.Post("/status", kithttp.NewServer(
		getStatusEndpoint(svc),
		decodeStatusRequest,
		encodeResponse,
		opts...,
	))

	r.Post("/events/:id", kithttp.NewServer(
		getEventsEndpoint(svc),
		decodeEventsRequest,
		encodeResponse,
		opts...,
	))

	r.GetFunc("/version", exapp.Version("exapp"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func decodeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req statusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		return nil, err
	}

	return req, nil
}

func decodeEventsRequest(_ context.Context, r *http.Request) (interface{}, error) {

	var req eventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case things.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case things.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case things.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case errInvalidQueryParams:
		w.WriteHeader(http.StatusBadRequest)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	default:
		switch err.(type) {
		case *json.SyntaxError:
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
