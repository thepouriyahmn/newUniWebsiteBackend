package main

import (
	"api/handler"
)

func main() {
	serviceURL := "localhost:8083"
	api := handler.NewHandler(serviceURL)
	api.RunApi()

}
