package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

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
	Tel          string       `json:"phone"`
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

func Find(db *sql.DB, tel string) (*Listing, error) {
	l := Listing{}
	name := ""
	l.Name = &name
	stmt, err := db.Prepare("SELECT name, surname, phone, updated FROM listing WHERE phone = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Query(tel)
	if err != nil {
		return nil, err
	}

	if result.Next() {
		err = result.Scan(l.Name, &l.Surname, &l.Tel, &l.LastAccessed.Time)
		if err != nil {
			return &l, err
		}
		return &l, nil
	}
	return nil, nil
}

func All(db *sql.DB) ([]Listing, error) {
	listings := []Listing{}
	rows, err := db.Query("SELECT name, surname, phone, updated FROM listing")
	if err != nil {
		return listings, err
	}
	defer rows.Close()
	for rows.Next() {
		listing := Listing{}
		name := ""
		listing.Name = &name
		err = rows.Scan(listing.Name, &listing.Surname, &listing.Tel, &listing.LastAccessed.Time)
		if err != nil {
			return listings, err
		}
		listings = append(listings, listing)
	}
	return listings, nil
}

func (l Listing) Save(db *sql.DB) error {
	stmt, err := db.Prepare(`
	INSERT INTO listing (name, surname, phone, updated)
	    	VALUES ($1, $2, $3, $4)
		ON CONFLICT (phone) DO
			UPDATE SET name=excluded.name, surname=excluded.surname, updated=excluded.updated
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(l.Name, l.Surname, l.Tel, l.LastAccessed.Time)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("unexpected rows affected from insert")
	}
	return nil
}

func Delete(db *sql.DB, phone string) error {
	stmt, err := db.Prepare("DELETE FROM listing WHERE phone = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(phone)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrRecordNotFound
	}
	return nil
}
