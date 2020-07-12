package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zkan/oddschana/logger"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	viper.SetDefault("port", "8000")
	viper.SetDefault("db.conn", "thaichana.db")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	l, _ := zap.NewDevelopment()
	defer l.Sync() // flushes buffer, if any

	hostname, _ := os.Hostname()
	l = l.With(zap.String("hostname", hostname))
	zap.ReplaceGlobals(l)

	r := mux.NewRouter()
	r.Use(logger.Middleware(l))

	db, err := sql.Open("sqlite3", viper.GetString("db.conn"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// This will serve files under http://localhost:8000/static/<filename>
	r.HandleFunc("/recently", Recently).Methods(http.MethodPost)
	r.HandleFunc("/checkin", CheckIn(NewInsertCheckIn(db))).Methods(http.MethodPost)
	r.HandleFunc("/checkout", CheckOut).Methods(http.MethodPost)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + viper.GetString("port"),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	zap.L().Info("start...", zap.String("port", viper.GetString("port")))
	log.Fatal(srv.ListenAndServe())
}

type Check struct {
	ID      int64 `json:"id"`
	PlaceID int64 `json:"place_id"`
}

type Location struct {
	Lat  float64
	Long float64
}

// Recently returns currently visited
func Recently(w http.ResponseWriter, r *http.Request) {

}

func NewInsertCheckIn(db *sql.DB) InFunc {
	return func(id, placeID int64) error {
		_, err := db.Exec("INSERT INTO visits VALUES(?, ?);", id, placeID)
		if err != nil {
			return err
		}
		return nil
	}
}

type InFunc func(id, placeID int64) error

func (fn InFunc) In(id, placeID int64) error {
	return fn(id, placeID)
}

type Iner interface {
	In(id, placeID int64) error
}

// CheckIn check-in to place, returns density (ok, too much)
func CheckIn(check Iner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var chk Check

		logger.L(r.Context()).Info("check-in")

		if err := json.NewDecoder(r.Body).Decode(&chk); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
			return
		}
		defer r.Body.Close()

		if err := check.In(chk.ID, chk.PlaceID); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
			return
		}
	}
}

// CheckOut check-out from place
func CheckOut(w http.ResponseWriter, r *http.Request) {

}
