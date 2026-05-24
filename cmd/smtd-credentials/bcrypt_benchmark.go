package main

import (
	"fmt"
	"log"
	"time"

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
