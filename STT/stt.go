package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"

	KEY = "19c1cb3c0aa848608fed5a5a8a23d640"
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

func SttService(speech []byte) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", URI, bytes.NewReader(speech))
	check(err)

	req.Header.Set("Content-Type",
		"audio/wav;codecs=audio/pcm;samplerate=16000")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEY)

	rsp, err2 := client.Do(req)
	check(err2)

	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)
		return string(body), nil
	} else {
		return "", errors.New("cannot convert to speech to text")
	}
}

func main() {
	speech, err1 := ioutil.ReadFile("speech.wav")
	check(err1)
	text, err2 := SpeechToText(speech)
	check(err2)
	fmt.Println(text)
	r := mux.NewRouter()
	r.HandleFunc("/tts", TextToSpeech).Methods("POST")
	http.ListenAndServe(":3002", r)
}
