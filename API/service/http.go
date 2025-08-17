package service

import (
	"fmt"
	"io"
	"net/http"
)

type Http struct {
	serviceURL string
}

func NewHttp() Http {
	return Http{}
}
func (h Http) SignUp(body io.Reader) error {

	resp, err := http.Post("http://"+h.serviceURL+"/signUp", "application/json", body)
	if err != nil {
		// fmt.Println("err is: ", err)
		// http.Error(w, "business service unavailable", http.StatusServiceUnavailable)
		return err
	}

	fmt.Println("done")
	defer resp.Body.Close()
	return nil
}
