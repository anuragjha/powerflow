package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// runs before everything else
func init() {
	// This function will be executed before everything else.
	fmt.Println("Init blockchain instance")
}

// Start handler
func Start(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok ok blockchain"))
}

// read response body
func readResponseBody(resp *http.Response) ([]byte, error) {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("cannot read request body")
	}
	defer resp.Body.Close()
	return respBody, nil
}

// read request body
func readRequestBody(r *http.Request) (string, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("cannot read request body")
	}
	defer r.Body.Close()
	return string(reqBody), nil
}
