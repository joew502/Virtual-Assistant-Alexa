package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	URI   = "http://localhost" // URI that the other microservices are hosted on
	ALPHA = "3001/alpha"       // Port and route for Alpha microservice
	STT   = "3002/stt"         // Port and route for STT microservice
	TTS   = "3003/tts"         // Port and route for TTS microservice
)

// Alexa - This function handles the incoming data,
//prepares it for use in the other microservices and passes it between said microservices.
//The final speech data is then formatted and sent as a response to the post request.
func Alexa(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil { // Decodes the incoming JSON data
		if qSpeech, ok := t["speech"].(string); ok { // Checks the data sent in is appropriate
			qSpeechJSON := "{\"speech\":\"" + qSpeech + "\"}"        // Formats the speech data to be sent to STT
			if qText, err := Service(qSpeechJSON, STT); err == nil { // Passes the question speech data to STT
				if aText, err := Service(qText, ALPHA); err == nil { // Passes the question text data to Alpha
					if aSpeech, err := Service(aText, TTS); err == nil { // Passes the answer text data to TTS
						var aJSON map[string]interface{}
						if err := json.Unmarshal([]byte(aSpeech), &aJSON); err == nil { // Handles returned speech data
							u := map[string]interface{}{"speech": aJSON["speech"]}
							w.WriteHeader(http.StatusOK)
							if err := json.NewEncoder(w).Encode(u); err != nil { // Responds to the request with speech
								w.WriteHeader(http.StatusInternalServerError)
							}
						} else {
							w.WriteHeader(http.StatusInternalServerError)
						}
					} else {
						w.WriteHeader(http.StatusInternalServerError)
					}
				} else {
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

// Service - Handles all outgoing http requests to STT, TTS and Alpha
func Service(serviceString string, serviceUri string) (string, error) {
	client := &http.Client{}
	sendUri := URI + ":" + serviceUri // Puts the correct URI together for the http request
	if req, err := http.NewRequest("POST", sendUri,
		bytes.NewBuffer([]byte(serviceString))); err == nil { // Initiates the http request
		if rsp, err := client.Do(req); err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(rsp.Body)
			if rsp.StatusCode == http.StatusOK { // Checks http request was handled without error
				if body, err := ioutil.ReadAll(rsp.Body); err == nil { // Handles and returns the response from the
					// http request
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
	r.HandleFunc("/alexa", Alexa).Methods("POST")
	http.ListenAndServe(":3000", r)
}
