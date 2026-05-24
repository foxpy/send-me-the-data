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
	rootCmd.AddCommand(bcryptBenchmarkCmd)
	rootCmd.AddCommand(addAdminCmd)
	rootCmd.AddCommand(deleteAdminCmd)
	rootCmd.AddCommand(changeAdminPasswordCmd)
	// TODO: list admins
	// TODO: list sessions

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
