package exapp

import (
	"encoding/json"
	"fmt"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/mteodor/edgex-app/events"
	nats "github.com/nats-io/go-nats"
)

// NatsMSGHandler process incomming message
func NatsMSGHandler(svc events.Service) nats.MsgHandler {

	return func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
		processMsg(svc, msg)
	}
}

func processMsg(svc events.Service, msg *nats.Msg) {
	data := model.Event{}
	json.Unmarshal(msg.Data, &data)

	err := svc.RegisterEvent(data)
	if err != nil {
		fmt.Printf("failed to save data")
	}
	fmt.Printf("unmarshalled:" + data.String())

}
