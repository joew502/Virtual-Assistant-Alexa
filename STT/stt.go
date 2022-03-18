package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	REGION = "uksouth" // Region of server to use for request
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US" // URI for Microsoft STT service
	KEY = "19c1cb3c0aa848608fed5a5a8a23d640" // Key for Microsoft STT service
)

type TextJSON struct { // Struct for use in decoding JSON returned from Microsoft STT Service
	RecognitionStatus string
	DisplayText       string
	Offset            int
	Duration          int
}

// SpeechToText - This function handles the incoming data,
//prepares it to be sent to the Microsoft STT Service and handles the returned text data.
func SpeechToText(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil { // Decodes the incoming JSON data
		if speechEncoded, ok := t["speech"].(string); ok { // Checks the data sent in is appropriate
			if speech, err := base64.StdEncoding.DecodeString(
				speechEncoded); err == nil { // Decodes the speech from base64
				if text, err := SttService(speech); err == nil { // Passes the speech to the Microsoft STT service
					var textJSON TextJSON
					if err := json.Unmarshal([]byte(text), &textJSON); err == nil { // Decodes the returned JSON
						u := map[string]interface{}{"text": textJSON.DisplayText}
						w.WriteHeader(http.StatusOK)
						err := json.NewEncoder(w).Encode(u) // Responds to the http request with the text
						if err != nil {
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

// SttService - Handles outgoing http requests the Microsoft STT service
func SttService(speech []byte) (string, error) {
	client := &http.Client{}
	if req, err := http.NewRequest("POST", URI, bytes.NewReader(speech)); err == nil { // Initiates the http request
		req.Header.Set("Content-Type",
			"audio/wav;codecs=audio/pcm;samplerate=16000") // Informs Microsoft STT service of the speech format
		req.Header.Set("Ocp-Apim-Subscription-Key", KEY) // Informs Microsoft STT service of the API Key
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
	r.HandleFunc("/stt", SpeechToText).Methods("POST")
	http.ListenAndServe(":3002", r)
}
