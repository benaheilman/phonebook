/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/benaheilman/phonebook/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List entries in the Phone Book",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, err := cmd.Parent().Flags().GetString("database")
		if err != nil {
			log.Fatal(err)
		}
		pb := db.LoadDatabase(dbPath)
		fmt.Println("EMPLOYEE PHONE BOOK")
		fmt.Println("===================")
		for _, listing := range pb.Listings {
			fmt.Print(listing.String())
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
