package service

import (
	"fmt"
	"io"
	"log"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Nc *nats.Conn
}

func NewNats() Nats {
	nc, _ := nats.Connect(nats.DefaultURL)
	return Nats{
		Nc: nc,
	}
}
func (n Nats) SignUp(body io.Reader) error {

	reqBodyBytes, err := io.ReadAll(body)
	if err != nil {
		fmt.Printf("reading error: %v", err)
		return err
	}
	// bodyStr := reqBodyBytes
	// data, _ := json.Marshal(reqBodyBytes)
	_, err = n.Nc.Request("signup.request", reqBodyBytes, nats.DefaultTimeout)
	if err != nil {
		log.Fatal("Error:", err)
	}
	return nil

}
