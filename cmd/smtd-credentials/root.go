package main

import "github.com/spf13/cobra"

var postgresURL string
var rootCmd = &cobra.Command{
	Use:   "smtd-credentials",
	Short: "Manage smtd administrator accounts and passwords",
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&postgresURL,
		"postgres-url",
		"postgres:///smtd?host=/run/postgresql",
		"PostgreSQL database connection URL",
	)
	rootCmd.AddCommand(bcryptBenchmarkCommand)
	rootCmd.AddCommand(addAdminCommand)
	rootCmd.AddCommand(deleteAdminCommand)
	// TODO: change admin password
	// TODO: list admins
	// TODO: list sessions

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
