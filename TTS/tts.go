package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".tts.speech.microsoft.com/cognitiveservices/v1"
	KEY    = "19c1cb3c0aa848608fed5a5a8a23d640"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TextToSpeech(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if text, ok := t["text"].(string); ok {
			main_text := []byte("<speak version=\"1.0\" xml:lang=\"en-US\"><voice xml:lang=\"en-US\" " +
				"name=\"en-US-JennyNeural\">" + text + "</voice></speak>")
			if speech, err := TtsService(main_text); err == nil {
				speech_encoded := base64.StdEncoding.EncodeToString([]byte(speech))
				u := map[string]interface{}{"Speech": speech_encoded}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(u)
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
	req, err := http.NewRequest("POST", URI, bytes.NewBuffer(text))
	check(err)

	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
	req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")

	rsp, err2 := client.Do(req)
	check(err2)

	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)
		return body, nil
	} else {
		return nil, errors.New("cannot convert text to speech")
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tts", TextToSpeech).Methods("POST")
	http.ListenAndServe(":3003", r)
}
