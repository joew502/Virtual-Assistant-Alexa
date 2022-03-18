package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	APPID = "G469T3-267PUX7QWR"                                    // Key for Wolfram Alpha API
	URI   = "http://api.wolframalpha.com/v1/result?appid=" + APPID // URI for Wolfram Alpha API
)

// Alpha - This function handles the incoming data,
//prepares it to be sent to the Wolfram Alpha API and handles the returned data.
func Alpha(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil { // Decodes the incoming JSON data
		if textIn, ok := t["text"].(string); ok { // Checks the data sent in is appropriate
			httpTextIn := url.QueryEscape(textIn)                     // Puts the question in the correct URL format
			if textOut, err := AlphaService(httpTextIn); err == nil { // Passes the speech to the Wolfram Alpha API
				u := map[string]interface{}{"text": textOut}
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(u); err != nil { // Responds to the http request with the text answer
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func AlphaService(text string) (string, error) {
	client := &http.Client{}
	sendUri := URI + "&i=" + text                                     // Formats the URI for the http request
	if req, err := http.NewRequest("GET", sendUri, nil); err == nil { // Initiates the http request
		if rsp, err := client.Do(req); err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(rsp.Body)
			if rsp.StatusCode == http.StatusOK { // Checks http request was handled without error
				if body, err := ioutil.ReadAll(rsp.Body); err == nil { // Handles and returns the response from the http request
					return string(body), nil
				} else {
					return "", err
				}
			} else {
				return "", errors.New(string(rune(rsp.StatusCode)))
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

// main - uses mux to handle all incoming http requests and routes them accordingly
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alpha", Alpha).Methods("POST")
	http.ListenAndServe(":3001", r)
}
