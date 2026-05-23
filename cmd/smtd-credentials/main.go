package main

import (
	cryptorand "crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	mathrand "math/rand"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/foxpy/send-me-the-data/src/idb/postgres"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
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
			f, err := os.OpenFile(passwordFile, os.O_RDONLY, 0)
			cobra.CheckErr(err)

			defer f.Close()
			var buf [256]byte
			// FIXME: this allows creating a newline-terminated password
			//        I don't think users will be happy about it
			n, err := f.Read(buf[:])
			cobra.CheckErr(err)

			password = string(buf[:n])
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

var (
	matchLowercase    = regexp.MustCompile(`[a-z]`)
	matchUppercase    = regexp.MustCompile(`[A-Z]`)
	matchNumeric      = regexp.MustCompile(`\d`)
	specialCharacters = `_+=%*&^$/\|.,:!(){}[]~`
)

func checkPasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !matchLowercase.MatchString(password) {
		return errors.New("password must have at least one lowercase letter")
	}

	if !matchUppercase.MatchString(password) {
		return errors.New("password must have at least one uppercase letter")
	}

	if !matchNumeric.MatchString(password) {
		return errors.New("password must have at least one numeric character")
	}

	if !strings.ContainsAny(password, specialCharacters) {
		return fmt.Errorf("password must have at least one special character from this list: %s", specialCharacters)
	}

	return nil
}

func getRandomPassword() (string, error) {
	var password []byte

	gen := func(numCharacters, alphabetSize int, alphabetBase byte) error {
		for range numCharacters {
			n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(alphabetSize)))
			if err != nil {
				return err
			}
			password = append(password, byte(n.Int64())+alphabetBase)
		}
		return nil
	}

	if gen(8, 26, 'a') != nil || gen(8, 26, 'A') != nil || gen(5, 10, '0') != nil {
		return "", nil
	}

	for range 3 {
		n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(specialCharacters))))
		if err != nil {
			return "", err
		}
		password = append(password, specialCharacters[n.Int64()])
	}

	mathrand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})
	return string(password), nil
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

func main() {
	goose.SetLogger(goose.NopLogger())

	rootCmd.Execute()
}
