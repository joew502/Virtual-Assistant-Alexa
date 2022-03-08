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
	REGION = "uksouth"
	URI    = "https://" + REGION + ".tts.speech.microsoft.com/cognitiveservices/v1"
	KEY    = "19c1cb3c0aa848608fed5a5a8a23d640"
)

func TextToSpeech(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if text, ok := t["text"].(string); ok {
			mainText := []byte("<speak version=\"1.0\" xml:lang=\"en-US\"><voice xml:lang=\"en-US\" " +
				"name=\"en-US-JennyNeural\">" + text + "</voice></speak>")
			if speech, err := TtsService(mainText); err == nil {
				speechEncoded := base64.StdEncoding.EncodeToString(speech)
				u := map[string]interface{}{"speech": speechEncoded}
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(u)
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

func TtsService(text []byte) ([]byte, error) {
	client := &http.Client{}
	if req, err := http.NewRequest("POST", URI, bytes.NewBuffer(text)); err == nil {
		req.Header.Set("Content-Type", "application/ssml+xml")
		req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
		req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
		if rsp, err := client.Do(req); err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(rsp.Body)
			if rsp.StatusCode == http.StatusOK {
				if body, err := ioutil.ReadAll(rsp.Body); err == nil {
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tts", TextToSpeech).Methods("POST")
	http.ListenAndServe(":3003", r)
}
