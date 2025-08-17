package main

import (
	"api/handler"
)

func main() {
	url := "localhost:8083"
	api := handler.NewHandler(url)
	api.RunApi()

}
