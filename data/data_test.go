package data

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadPhonebookEmpty(t *testing.T) {
	phonebook, err := LoadPhonebook(strings.NewReader(``))
	assert.Nil(t, err)
	assert.Equal(t, []Listing{}, phonebook.Listings)
}

func TestLoadPhonebookSingle(t *testing.T) {
	input := `
	[{
		"name": "first",
		"surname": "last",
		"telephone_number": "3125551212",
		"last_accessed": "2006-01-02T15:04:05Z"
	}]
	`
	phonebook, err := LoadPhonebook(strings.NewReader(input))
	assert.Nil(t, err)
	listing := phonebook.Listings[0]

	assert.Equal(t, "first", *listing.Name)
	assert.Equal(t, "last", listing.Surname)
	assert.Equal(t, "3125551212", listing.Tel)
	assert.Equal(t, NullableTime{time.Time(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))}, listing.LastAccessed)
}

func TestLoadPhonebookNullName(t *testing.T) {
	input := `
	[{
		"name": null,
		"surname": "last",
		"telephone_number": "3125551212"
	}]
	`
	phonebook, err := LoadPhonebook(strings.NewReader(input))
	assert.Nil(t, err)
	listing := phonebook.Listings[0]

	assert.Nil(t, listing.Name)
}

func TestLoadPhonebookNoName(t *testing.T) {
	input := `
	[{
		"surname": "last",
		"telephone_number": "3125551212"
	}]
	`
	phonebook, err := LoadPhonebook(strings.NewReader(input))
	assert.Nil(t, err)
	listing := phonebook.Listings[0]

	assert.Nil(t, listing.Name)
}

func TestLoadPhonebookNullTime(t *testing.T) {
	input := `
	[{
		"surname": "last",
		"telephone_number": "3125551212",
		"last_accessed": null
	}]
	`
	phonebook, err := LoadPhonebook(strings.NewReader(input))
	assert.Nil(t, err)
	entity := phonebook.Listings[0]

	assert.True(t, entity.LastAccessed.IsZero())
}

func TestLoadPhonebookNoTime(t *testing.T) {
	input := `
	[{
		"surname": "last",
		"telephone_number": "3125551212"
	}]
	`
	phonebook, err := LoadPhonebook(strings.NewReader(input))
	assert.Nil(t, err)
	entity := phonebook.Listings[0]

	assert.True(t, entity.LastAccessed.IsZero())
}

func TestPhonebookSaveEmpty(t *testing.T) {
	w := bytes.Buffer{}
	ph := Phonebook{}
	assert.Nil(t, ph.Save(&w))
	assert.Equal(t, "[]", strings.TrimSpace(w.String()))
}

func ref[t any](i t) *t {
	return &i
}

func TestPhonebookSaveSingle(t *testing.T) {
	w := bytes.Buffer{}
	ph := Phonebook{Listings: []Listing{
		{
			Name:         ref("first"),
			Surname:      "last",
			Tel:          "3125551212",
			LastAccessed: NullableTime{time.Time(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))}},
	}}
	assert.Nil(t, ph.Save(&w))
	expected := `[{"name":"first","surname":"last","telephone_number":"3125551212","last_accessed":"2006-01-02T15:04:05Z"}]`
	assert.Equal(t, expected, strings.TrimSpace(w.String()))
}

func TestPhonebookSaveNilName(t *testing.T) {
	w := bytes.Buffer{}
	ph := Phonebook{Listings: []Listing{
		{
			Name:         nil,
			Surname:      "last",
			Tel:          "3125551212",
			LastAccessed: NullableTime{time.Time(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))},
		},
	}}
	assert.Nil(t, ph.Save(&w))
	expected := `[{"surname":"last","telephone_number":"3125551212","last_accessed":"2006-01-02T15:04:05Z"}]`
	assert.Equal(t, expected, strings.TrimSpace(w.String()))
}

func TestPhonebookSaveZeroTime(t *testing.T) {
	w := bytes.Buffer{}
	ph := Phonebook{Listings: []Listing{
		{
			Name:         ref("first"),
			Surname:      "last",
			Tel:          "3125551212",
			LastAccessed: NullableTime{},
		},
	}}
	assert.Nil(t, ph.Save(&w))
	expected := `[{"name":"first","surname":"last","telephone_number":"3125551212","last_accessed":null}]`
	assert.Equal(t, expected, strings.TrimSpace(w.String()))
}
