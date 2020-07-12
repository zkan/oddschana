package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.HandleFunc("/recently", Recently).Methods(http.MethodPost)
	r.Handle("/checkin", &CheckIn{InsertCheckIn: insertCheckIn}).Methods(http.MethodPost)
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
	ID      int64 `json:"id"`
	PlaceID int64 `json:"people_id"`
}

type Location struct {
	Lat  float64
	Long float64
}

// Recently returns currently visited
func Recently(w http.ResponseWriter, r *http.Request) {

}

func insertCheckIn(id, placeID int64) error {
	db, err := sql.Open("sqlite3", "thaichana.db")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO visits VALUES(?, ?);", id, placeID)
	if err != nil {
		return err
	}
	return nil
}

type InsertCheckInFunc func(id, placeID int64) error

type CheckIn struct {
	InsertCheckIn InsertCheckInFunc
}

// CheckIn check-in to place, returns density (ok, too much)
func (c *CheckIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var chk Check
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&chk); err != nil {
		fmt.Println(chk.ID)
		fmt.Println(chk.PlaceID)
		fmt.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer r.Body.Close()

	if err := c.InsertCheckIn(chk.ID, chk.PlaceID); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
}

// CheckOut check-out from place
func CheckOut(w http.ResponseWriter, r *http.Request) {

}
