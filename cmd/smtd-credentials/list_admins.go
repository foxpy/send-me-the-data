package main

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/src/idb/postgres"
	"github.com/spf13/cobra"
)

var listAdminsCmd = &cobra.Command{
	Use:   "list-admins",
	Short: "Print all registered administrator usernames",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := postgres.NewPostgres(postgresURL)
		cobra.CheckErr(err)

		admins, err := db.GetAllAdmins()
		cobra.CheckErr(err)

		for admin := range admins {
			fmt.Println(admin)
		}
	},
}
