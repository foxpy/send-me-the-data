package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/foxpy/send-me-the-data/src/idb/postgres"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var bcryptBenchmarkCommand = &cobra.Command{
	Use:   "bcrypt-benchmark",
	Short: "Benchmark bcrypt performance and print table with results for every cost level",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("     Cost level | Time")
		for cost := bcrypt.MinCost; cost < bcrypt.MaxCost; cost++ {
			begin := time.Now()
			_, err := bcrypt.GenerateFromPassword([]byte(`password`), cost)
			if err != nil {
				log.Fatal("bcrypt error: ", err.Error())
			}

			duration := time.Since(begin)
			if cost != bcrypt.DefaultCost {
				fmt.Printf("%15d | %s\n", cost, duration)
			} else {
				fmt.Printf("%15s | %s\n", fmt.Sprintf("(default) %d", cost), duration)
			}
		}
	},
}

var bcryptCost int
var addAdminCommand = &cobra.Command{
	Use:   "add-admin [username] [password]",
	Short: "Register new administrator",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]

		// TODO: check password strength

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

		fmt.Printf("Added admin '%s' with password hash '%s'\n", username, string(passwordHash))
	},
}

var deleteAdminCommand = &cobra.Command{
	Use:   "delete-admin [username]",
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

var postgresURL string
var rootCmd = &cobra.Command{
	Use:   "smtd-credentials",
	Short: "Manage smtd administrator accounts and passwords",
}

func init() {
	addAdminCommand.PersistentFlags().IntVar(
		&bcryptCost,
		"bcrypt-cost",
		bcrypt.DefaultCost,
		fmt.Sprintf("bcrypt cost level, a value in range [%d, %d]", bcrypt.MinCost, bcrypt.MaxCost),
	)

	rootCmd.PersistentFlags().StringVar(
		&postgresURL,
		"postgres-url",
		"postgres:///smtd?host=/run/postgresql",
		"PostgreSQL database connection URL",
	)
	rootCmd.AddCommand(bcryptBenchmarkCommand)
	rootCmd.AddCommand(addAdminCommand)
	rootCmd.AddCommand(deleteAdminCommand)
	// TODO: list admins
	// TODO: list sessions

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	goose.SetLogger(goose.NopLogger())

	rootCmd.Execute()
}
