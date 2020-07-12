package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckInHandler(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(&Check{ID: 123, PlaceID: 234})
	req := httptest.NewRequest("POST", "http://example.com/foo", payload)
	w := httptest.NewRecorder()

	var fn InFunc = func(id, placeID int64) error {
		return nil
	}

	CheckIn(fn)(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}

func TestSealMiddleware(t *testing.T) {
	payload := bytes.NewBuffer([]byte("ewogICAgImlkIjoxLAogICAgInBsYWNlX2lkIjogMQp9"))

	req := httptest.NewRequest("POST", "http://example.com/foo", payload)
	w := httptest.NewRecorder()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(string(b))
	}

	SealMiddleware()(http.HandlerFunc(handler)).ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
