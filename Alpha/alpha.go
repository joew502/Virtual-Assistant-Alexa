package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	APPID = "G469T3-267PUX7QWR"
	URI   = "http://api.wolframalpha.com/v1/result?appid=" + APPID
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Alpha(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if textIn, ok := t["text"].(string); ok {
			httpTextIn := url.QueryEscape(textIn)
			if textOut, err := AlphaService(httpTextIn); err == nil {
				u := map[string]interface{}{"text": textOut}
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

func AlphaService(text string) (string, error) {
	client := &http.Client{}
	sendUri := URI + "&i=" + text
	req, err := http.NewRequest("GET", sendUri, nil)
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
	r.HandleFunc("/alpha", Alpha).Methods("POST")
	http.ListenAndServe(":3001", r)
}
