package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Encapsulating structure to marshal zero time as null
type NullableTime struct {
	time.Time
}

func (t NullableTime) MarshalJSON() ([]byte, error) {
	switch t.IsZero() {
	case true:
		return []byte("null"), nil
	default:
		return t.Time.MarshalJSON()
	}
}

// Listing in the phone book
type Listing struct {
	Name         *string      `json:"name,omitempty"`
	Surname      string       `json:"surname"`
	Tel          string       `json:"telephone_number"`
	LastAccessed NullableTime `json:"last_accessed,omitempty"`
}

func (l *Listing) String() string {
	r := ""
	if l.Name != nil {
		r += fmt.Sprintf("Name:             %s\n", *l.Name)
	}
	r += fmt.Sprintf("Surname:          %s\n", l.Surname)
	r += fmt.Sprintf("Telephone Number: %s\n", l.Tel)
	return r
}

// The collection of telephone contacts
type Phonebook struct {
	Listings []Listing
}

func LoadPhonebook(r io.Reader) (Phonebook, error) {
	decoder := json.NewDecoder(r)
	entities := []Listing{}
	if err := decoder.Decode(&entities); err != nil && err != io.EOF {
		return Phonebook{Listings: []Listing{}}, err
	}
	return Phonebook{Listings: entities}, nil
}

func (pb Phonebook) Save(w io.Writer) error {
	if pb.Listings == nil {
		pb.Listings = []Listing{}
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(pb.Listings); err != nil {
		return err
	}
	return nil
}
