package main

import (
	"github.com/pressly/goose/v3"
)

func main() {
	goose.SetLogger(goose.NopLogger())
	rootCmd.Execute()
}
