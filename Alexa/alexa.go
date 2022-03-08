package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	URI   = "http://localhost"
	ALPHA = "3001/alpha"
	STT   = "3002/stt"
	TTS   = "3003/tts"
)

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
							if err := json.NewEncoder(w).Encode(u); err != nil {
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

func Service(serviceString string, serviceUri string) (string, error) {
	client := &http.Client{}
	sendUri := URI + ":" + serviceUri
	if req, err := http.NewRequest("POST", sendUri, bytes.NewBuffer([]byte(serviceString))); err == nil {
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
	r.HandleFunc("/alexa", Alexa).Methods("POST")
	http.ListenAndServe(":3000", r)
}
