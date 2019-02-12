package exapp

import (
	"encoding/json"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	nats "github.com/nats-io/go-nats"
)

// NatsMSGHandler process incomming message
func NatsMSGHandler(svc Service) nats.MsgHandler {

	return func(msg *nats.Msg) {
		//logger.Info(fmt.Sprintf("Received a message: %s\n", string(msg.Data)))
		processMsg(svc, msg)
	}
}

func processMsg(svc Service, msg *nats.Msg) {
	data := model.Event{}
	json.Unmarshal(msg.Data, &data)

	err := svc.RegisterEvent(data)
	if err != nil {
		//logger.Error("failed to save data")
	}

}
