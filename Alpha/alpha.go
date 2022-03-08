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
	APPID = "G469T3-267PUX7QWR"
	URI   = "http://api.wolframalpha.com/v1/result?appid=" + APPID
)

func Alpha(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if textIn, ok := t["text"].(string); ok {
			httpTextIn := url.QueryEscape(textIn)
			if textOut, err := AlphaService(httpTextIn); err == nil {
				u := map[string]interface{}{"text": textOut}
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(u); err != nil {
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
	sendUri := URI + "&i=" + text
	if req, err := http.NewRequest("GET", sendUri, nil); err == nil {
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
	r.HandleFunc("/alpha", Alpha).Methods("POST")
	http.ListenAndServe(":3001", r)
}
