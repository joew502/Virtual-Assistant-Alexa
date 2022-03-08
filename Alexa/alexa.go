package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	URI   = "http://localhost"
	ALPHA = "3001/alpha"
	STT   = "3002/stt"
	TTS   = "3003/tts"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Alexa(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if qSpeech, ok := t["speech"].(string); ok {
			qSpeechJSON := "{\"speech\":\"" + qSpeech + "\"}"
			if qText, err := Service(qSpeechJSON, STT); err == nil {
				if aText, err := Service(qText, ALPHA); err == nil {
					fmt.Println(aText)
					if aSpeech, err := Service(aText, TTS); err == nil {
						var aJSON map[string]interface{}
						if err := json.Unmarshal([]byte(aSpeech), &aJSON); err == nil {
							u := map[string]interface{}{"speech": aJSON["speech"]}
							w.WriteHeader(http.StatusOK)
							json.NewEncoder(w).Encode(u)
						} else {
							w.WriteHeader(http.StatusInternalServerError)
						}
						//fmt.Println(aSpeech)
						//w.WriteHeader(http.StatusOK)
						//json.NewEncoder(w).Encode(aSpeech)
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

func Service(serviceString string, serviceUri string) (string, error) {
	client := &http.Client{}
	sendUri := URI + ":" + serviceUri
	req, err := http.NewRequest("POST", sendUri, bytes.NewBuffer([]byte(serviceString)))
	check(err)

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
	r := mux.NewRouter()
	r.HandleFunc("/alexa", Alexa).Methods("POST")
	http.ListenAndServe(":3000", r)
}
