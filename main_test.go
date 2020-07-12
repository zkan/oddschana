package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

type insertStub struct{}

func (insertStub) In(id, placeID int64) error {
	return nil
}

func TestCheckInHandler(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(&Check{ID: 123, PlaceID: 234})
	req := httptest.NewRequest("POST", "http://example.com/foo", payload)
	w := httptest.NewRecorder()

	CheckIn(insertStub{})(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
