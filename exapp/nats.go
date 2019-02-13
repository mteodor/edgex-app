package exapp

import (
	"encoding/json"
	"fmt"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	nats "github.com/nats-io/go-nats"
)

// NatsMSGHandler process incomming message
func NatsMSGHandler(svc Service) nats.MsgHandler {

	return func(msg *nats.Msg) {
		(*svc.GetLogger()).Info((fmt.Sprintf("Received a message: %s\n", string(msg.Data))))
		processMsg(svc, msg)
	}
}

func processMsg(svc Service, msg *nats.Msg) {
	data := model.Event{}
	json.Unmarshal(msg.Data, &data)

	err := svc.RegisterEvent(data)
	if err != nil {
		(*svc.GetLogger()).Error("failed to save data")
	}

}
