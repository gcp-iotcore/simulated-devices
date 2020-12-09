package httphandler

import (
	"io"
	"log"
	"net/http"
)

const (
	baseURL string = "https://cloudiotdevice.googleapis.com/v1/"
)

var (
	client HTTPClient = &http.Client{}
)

//HTTPClient: interface to be used later on for mocking Do for testing
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func MakeHTTPCall(method string, apiURL string, JWTToken string, reqBody io.Reader) (*http.Response, error) {
	log.Println("making http calls")
	req, err := http.NewRequest(method, (baseURL + apiURL), reqBody)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", ("Bearer " + JWTToken))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return response, err
}
