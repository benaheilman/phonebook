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
	defer fh.Close()

	pb, err := data.LoadPhonebook(fh)
	if err != nil {
		log.Fatalf("Database corruption: %s", err)
	}
	return pb
}

func SaveDatabase(pb data.Phonebook, path string) error {
	fh, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	if err := pb.Save(fh); err != nil {
		return err
	}
	return nil
}
