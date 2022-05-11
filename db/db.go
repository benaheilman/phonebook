package db

import (
	"log"
	"os"

	"github.com/benaheilman/phonebook/data"
)

func LoadDatabase(path string) data.Phonebook {
	fh, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	pb, err := data.LoadPhonebook(fh)
	if err != nil {
		log.Fatalf("Database corruption: %s", err)
	}
	return pb
}
