package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.HandleFunc("/recently", Recently).Methods(http.MethodPost)
	r.HandleFunc("/checkin", CheckIn).Methods(http.MethodPost)
	r.HandleFunc("/checkout", CheckOut).Methods(http.MethodPost)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type Check struct {
	ID      int64
	PlaceID int64
}

type Location struct {
	Lat  float64
	Long float64
}

// Recently returns currently visited
func Recently(w http.ResponseWriter, r *http.Request) {

}

// CheckIn check-in to place, returns density (ok, too much)
func CheckIn(w http.ResponseWriter, r *http.Request) {

}

// CheckOut check-out from place
func CheckOut(w http.ResponseWriter, r *http.Request) {

}
