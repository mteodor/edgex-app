package http

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/mteodor/edgex-app/exapp"
)

func getStatusEndpoint(svc exapp.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		req := request.(statusRequest)
		svc.Logger().Info(fmt.Sprintf("get status for %s" + req.Name))
		var greeting string
		greeting = fmt.Sprintf("Hello, %s, I'm working fine", req.Name)
		res := statusResponse{
			Greeting: greeting,
			Err:      "",
		}
		return res, nil

	}
}

func getEventsEndpoint(svc exapp.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		req := request.(eventRequest)
		svc.Logger().Info(fmt.Sprintf("get events for %s" + req.Id))
		res, err := svc.RetrieveByID(req.Id)

		return res, err

	}
}
