/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benaheilman/phonebook/data"
	"github.com/benaheilman/phonebook/db"
	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Inserts a listing into the phone book",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, err := cmd.Parent().Flags().GetString("database")
		if err != nil {
			log.Fatal(err)
		}
		pb := db.LoadDatabase(dbPath)

		surname, err := cmd.Flags().GetString("surname")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		phone, err := cmd.Flags().GetString("phone-number")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		l := data.Listing{
			Surname:      surname,
			Tel:          phone,
			LastAccessed: data.NullableTime{Time: time.Now().UTC()},
		}
		if cmd.Flag("name").Changed {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			l.Name = &name
		}
		pb.Listings = append(pb.Listings, l)

		if err := db.SaveDatabase(pb, dbPath); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Listing saved...")
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	insertCmd.Flags().StringP("name", "n", "", "Person's name")
	insertCmd.Flags().StringP("surname", "s", "", "Person's surname")
	insertCmd.MarkFlagRequired("surname")
	insertCmd.Flags().StringP("phone-number", "p", "", "Person's phone number (must be unique)")
	insertCmd.MarkFlagRequired("phone-number")
}
