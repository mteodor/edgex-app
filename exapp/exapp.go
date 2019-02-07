package exapp

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	model "github.com/edgexfoundry/edgex-go/pkg/models"
	nats "github.com/nats-io/go-nats"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler() nats.MsgHandler {

	return func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
		processMsg(msg)
	}
}
func MakeHTTPHandler() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// If request is GET with correct syntax
		fmt.Fprintf(w, "%q", html.EscapeString("edgex-app iw working\n"))

	})

}
func processMsg(msg *nats.Msg) {
	data := new(model.Event)
	json.Unmarshal(msg.Data, &data)
	fmt.Printf("unmarshalled:" + data.String())

}
