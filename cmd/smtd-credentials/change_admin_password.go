package main

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/src/idb/postgres"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var changeAdminPasswordCmd = &cobra.Command{
	Use:   "change-admin-password <username>",
	Short: "Change administrator password, password is read from stdin by default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(passwordFile) > 0 && generatePassword {
			cobra.CheckErr(fmt.Errorf("conflicting options --password-file and --generate-password"))
		}

		username := args[0]
		var password string
		var err error
		if generatePassword {
			password, err = getRandomPassword()
		} else if len(passwordFile) > 0 {
			password, err = readPasswordFile(passwordFile)
		} else {
			password, err = readPasswordStdin()
		}
		cobra.CheckErr(err)

		err = checkPasswordStrength(password)
		cobra.CheckErr(err)

		if bcryptCost < bcrypt.MinCost || bcryptCost > bcrypt.MaxCost {
			cobra.CheckErr(fmt.Errorf("invalid bcrypt cost"))
		}

		db, err := postgres.NewPostgres(postgresURL)
		cobra.CheckErr(err)

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		cobra.CheckErr(err)

		err = db.UpdateAdmin(username, passwordHash)
		cobra.CheckErr(err)

		if generatePassword {
			fmt.Printf("Updated admin '%s' password to:\n%s\n", username, password)
		} else {
			fmt.Printf("Updated admin '%s' password\n", username)
		}
	},
}

func init() {
	changeAdminPasswordCmd.PersistentFlags().IntVar(
		&bcryptCost,
		"bcrypt-cost",
		bcrypt.DefaultCost,
		fmt.Sprintf("bcrypt cost level, a value in range [%d, %d]", bcrypt.MinCost, bcrypt.MaxCost),
	)
	changeAdminPasswordCmd.PersistentFlags().StringVar(
		&passwordFile,
		"password-file",
		"",
		"Read password from specified file",
	)
	changeAdminPasswordCmd.PersistentFlags().BoolVar(
		&generatePassword,
		"generate-password",
		false,
		"Generate password and print to stdin",
	)
}
