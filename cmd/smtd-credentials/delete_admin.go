package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/foxpy/send-me-the-data/src/idb/postgres"
	"github.com/spf13/cobra"
)

var deleteAdminCmd = &cobra.Command{
	Use:   "delete-admin <username>",
	Short: "Delete administrator",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		db, err := postgres.NewPostgres(postgresURL)
		cobra.CheckErr(err)

		err = db.DeleteAdmin(username)
		if errors.Is(err, sql.ErrNoRows) {
			cobra.CheckErr(fmt.Errorf("Admin '%s' doesn't exist", username))
		} else {
			cobra.CheckErr(err)
		}

		fmt.Printf("Deleted admin '%s'\n", username)
	},
}
