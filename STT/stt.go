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
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
	KEY = "19c1cb3c0aa848608fed5a5a8a23d640"
)

type TextJSON struct {
	RecognitionStatus string
	DisplayText       string
	Offset            int
	Duration          int
}

func SpeechToText(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if speechEncoded, ok := t["speech"].(string); ok {
			if speech, err := base64.StdEncoding.DecodeString(speechEncoded); err == nil {
				if text, err := SttService(speech); err == nil {
					var textJSON TextJSON
					if err := json.Unmarshal([]byte(text), &textJSON); err == nil {
						u := map[string]interface{}{"text": textJSON.DisplayText}
						w.WriteHeader(http.StatusOK)
						err := json.NewEncoder(w).Encode(u)
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

func SttService(speech []byte) (string, error) {
	client := &http.Client{}
	if req, err := http.NewRequest("POST", URI, bytes.NewReader(speech)); err == nil {
		req.Header.Set("Content-Type",
			"audio/wav;codecs=audio/pcm;samplerate=16000")
		req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
		if rsp, err := client.Do(req); err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(rsp.Body)
			if rsp.StatusCode == http.StatusOK {
				if body, err := ioutil.ReadAll(rsp.Body); err == nil {
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stt", SpeechToText).Methods("POST")
	http.ListenAndServe(":3002", r)
}
