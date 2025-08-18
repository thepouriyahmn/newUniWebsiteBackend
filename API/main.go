package main

import (
	"api/handler"
	"api/service"
)

func main() {
	connService := "nats"
	var Service service.Services
	serviceURL := "localhost:8083"
	httpAdapter := service.NewHttp(serviceURL)
	natsAdapter := service.NewNats()
	if connService == "nats" {
		Service = natsAdapter
	} else {
		Service = httpAdapter
	}

	api := handler.NewHandler(serviceURL, Service)
	api.RunApi()

}
