package main

import (
	"api/handler"
	"api/service"
)

func main() {
	connService := "nats"
	var Service service.Services
	http := service.NewHttp()
	nats := service.NewNats()
	if connService == "nats" {
		Service = nats
	} else {
		Service = http
	}

	serviceURL := "localhost:8083"
	api := handler.NewHandler(serviceURL, Service)
	api.RunApi()

}
