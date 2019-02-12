package http

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/mteodor/edgex-app/events"
)

func getStatusEndpoint(svc events.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		fmt.Println("get status")
		req := request.(statusRequest)
		var greeting string
		greeting = fmt.Sprintf("Hello, %s, I'm working fine", req.Name)
		res := statusResponse{
			Greeting: greeting,
			Err:      "",
		}
		return res, nil

	}
}

func getEventsEndpoint(svc events.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		fmt.Println("get event")
		req := request.(eventRequest)

		res, err := svc.RetrieveByID(req.Id)

		return res, err

	}
}
