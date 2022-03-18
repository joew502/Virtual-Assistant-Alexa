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
	REGION = "uksouth"                                                              // Region of server to use for request
	URI    = "https://" + REGION + ".tts.speech.microsoft.com/cognitiveservices/v1" // URI for Microsoft TTS service
	KEY    = "19c1cb3c0aa848608fed5a5a8a23d640"                                     // Key for Microsoft TTS service
)

// TextToSpeech - This function handles the incoming data,
//prepares it to be sent to the Microsoft TTS Service and handles the returned speech data.
func TextToSpeech(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil { // Decodes the incoming JSON data
		if text, ok := t["text"].(string); ok { // Checks the data sent in is appropriate
			mainText := []byte("<speak version=\"1.0\" xml:lang=\"en-US\"><voice xml:lang=\"en-US\" " +
				"name=\"en-US-JennyNeural\">" + text + "</voice></speak>") // Puts together the XML to be sent to the
			// Microsoft TTS Service
			if speech, err := TtsService(mainText); err == nil { // Passes the speech to the Microsoft TTS service
				speechEncoded := base64.StdEncoding.EncodeToString(speech) // Encodes the returned speech to base64
				u := map[string]interface{}{"speech": speechEncoded}
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(u) // Responds to the http request with the speech
				if err != nil {
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

// TtsService - Handles outgoing http requests the Microsoft TTS service
func TtsService(text []byte) ([]byte, error) {
	client := &http.Client{}
	if req, err := http.NewRequest("POST", URI, bytes.NewBuffer(text)); err == nil { // Initiates the http request
		req.Header.Set("Content-Type", "application/ssml+xml") // Informs Microsoft TTS service of the request format
		req.Header.Set("Ocp-Apim-Subscription-Key", KEY)       // Informs Microsoft TTS service of the API Key
		req.Header.Set("X-Microsoft-OutputFormat",
			"riff-16khz-16bit-mono-pcm") // Informs Microsoft TTS service of the format to respond in
		if rsp, err := client.Do(req); err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(rsp.Body)
			if rsp.StatusCode == http.StatusOK { // Checks http request was handled without error
				if body, err := ioutil.ReadAll(rsp.Body); err == nil { // Handles and returns the response from the http request
					return body, nil
				} else {
					return nil, err
				}
			} else {
				return nil, errors.New(string(rune(rsp.StatusCode)))
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

}

// main - uses mux to handle all incoming http requests and routes them accordingly
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tts", TextToSpeech).Methods("POST")
	http.ListenAndServe(":3003", r)
}
