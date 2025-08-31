package messagebrokers

import (
	"encoding/json"
)

//	type natsRequest struct {
//		Token string          `json:"token,omitempty"`
//		Query string          `json:"query,omitempty"`
//		Body  json.RawMessage `json:"body,omitempty"`
//	}
type IMessageBroker interface {
	// Subscribe(topic string, cb func()) error
	Publish(subject string, req json.RawMessage) error
	Run()
}
