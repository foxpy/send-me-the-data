package main

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/foxpy/send-me-the-data/src/idb/postgres"
	"github.com/lib/pq"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var bcryptCost int
var passwordFile string
var generatePassword bool
var addAdminCommand = &cobra.Command{
	Use:   "add-admin [username]",
	Short: "Register new administrator, password is read from stdin by default",
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
			cobra.CheckErr(err)

		} else if len(passwordFile) > 0 {
			password, err = readPasswordFile(passwordFile)
			cobra.CheckErr(err)

		} else {
			fmt.Print("Type password: ")
			// FIXME: this function leaves terminal in a broken state after Ctrl-C
			// FIXME: this function ignores Ctrl-D
			pwd1, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			cobra.CheckErr(err)

			fmt.Print("Retype password: ")
			pwd2, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			cobra.CheckErr(err)

			if !slices.Equal(pwd1, pwd2) {
				cobra.CheckErr(fmt.Errorf("Sorry, passwords do not match"))
			}

			password = string(pwd1)
		}

		err = checkPasswordStrength(password)
		cobra.CheckErr(err)

		if bcryptCost < bcrypt.MinCost || bcryptCost > bcrypt.MaxCost {
			cobra.CheckErr(fmt.Errorf("invalid bcrypt cost"))
		}

		db, err := postgres.NewPostgres(postgresURL)
		cobra.CheckErr(err)

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		cobra.CheckErr(err)

		err = db.CreateAdmin(username, passwordHash)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			cobra.CheckErr(fmt.Errorf("admin '%s' already exists", username))
		} else if err != nil {
			cobra.CheckErr(err)
		}

		if generatePassword {
			fmt.Printf("Added admin '%s' with password:\n%s\n", username, password)
		} else {
			fmt.Printf("Added admin '%s'\n", username)
		}
	},
}

func init() {
	addAdminCommand.PersistentFlags().IntVar(
		&bcryptCost,
		"bcrypt-cost",
		bcrypt.DefaultCost,
		fmt.Sprintf("bcrypt cost level, a value in range [%d, %d]", bcrypt.MinCost, bcrypt.MaxCost),
	)
	addAdminCommand.PersistentFlags().StringVar(
		&passwordFile,
		"password-file",
		"",
		"Read password from specified file",
	)
	addAdminCommand.PersistentFlags().BoolVar(
		&generatePassword,
		"generate-password",
		false,
		"Generate password and print to stdin",
	)
}
