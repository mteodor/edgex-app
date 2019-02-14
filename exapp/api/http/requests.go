package http

type statusRequest struct {
	Name string `json:"name"`
}

type statusResponse struct {
	Greeting string `json:"greeting"`
	Err      string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type eventRequest struct {
	Id string `json:"id"`
}

type eventResponse interface{}
