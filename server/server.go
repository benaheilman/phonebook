package server

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/benaheilman/phonebook/data"
	"github.com/benaheilman/phonebook/db"
	"github.com/gorilla/mux"
)

var database = db.OpenDatabase("phonebook.sqlite")

func listingsHandler(w http.ResponseWriter, r *http.Request) {
	listings, err := data.All(database)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	e.Encode(listings)
}

func listingGetHandler(w http.ResponseWriter, r *http.Request) {
	phone := mux.Vars(r)["phone"]
	listing, err := data.Find(database, phone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if listing == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 listing not found"))
		return
	}
	o, err := json.MarshalIndent(listing, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(o))
}

func listingsPostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	listing := data.Listing{}
	listing.LastAccessed = data.NullableTime{Time: time.Now()}
	err = json.Unmarshal(b, &listing)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = listing.Save(database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("200 OK"))
}

func listingPutHandler(w http.ResponseWriter, r *http.Request) {
	phone := mux.Vars(r)["phone"]

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	proposed := data.Listing{}
	proposed.LastAccessed = data.NullableTime{Time: time.Now()}
	err = json.Unmarshal(b, &proposed)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if proposed.Tel != phone {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("cannot change phone number"))
		return
	}
	err = proposed.Save(database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func listingDeleteHandler(w http.ResponseWriter, r *http.Request) {
	phone := mux.Vars(r)["phone"]
	err := data.Delete(database, phone)
	switch err {
	case data.ErrRecordNotFound:
		w.WriteHeader(http.StatusNoContent)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func Serve() {
	r := mux.NewRouter()
	r.HandleFunc("/listing", listingsHandler).Methods("GET")
	r.HandleFunc("/listing", listingsPostHandler).Methods("POST")
	r.HandleFunc("/listing/{phone}", listingGetHandler).Methods("GET")
	r.HandleFunc("/listing/{phone}", listingPutHandler).Methods("PUT")
	r.HandleFunc("/listing/{phone}", listingDeleteHandler).Methods("DELETE")

	laddr, err := net.ResolveTCPAddr("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	listen, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(listen, r)
}
