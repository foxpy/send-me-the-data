package main

import (
	"bufio"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

var (
	matchLowercase      = regexp.MustCompile(`[a-z]`)
	matchUppercase      = regexp.MustCompile(`[A-Z]`)
	matchNumeric        = regexp.MustCompile(`\d`)
	specialCharacters   = `_+=%*&^$/\|.,:!(){}[]~`
	errNotPasswordFile  = errors.New("the specified file does not contain a password")
	errShortPassword    = errors.New("password must be at least 8 characters long")
	errMissingLowercase = errors.New("password must have at least one lowercase letter")
	errMissingUppercase = errors.New("password must have at least one uppercase letter")
	errMissingNumeric   = errors.New("password must have at least one numeric letter")
)

func checkPasswordStrength(password string) error {
	if len(password) < 8 {
		return errShortPassword
	}

	if !matchLowercase.MatchString(password) {
		return errMissingLowercase
	}

	if !matchUppercase.MatchString(password) {
		return errMissingUppercase
	}

	if !matchNumeric.MatchString(password) {
		return errMissingNumeric
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

func readPasswordFile(path string) (string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}

	r := bufio.NewReaderSize(f, 256)
	line, isPrefix, err := r.ReadLine()
	if err != nil {
		return "", err
	}

	if isPrefix {
		return "", errNotPasswordFile
	}

	return string(line), nil
}

const ioctlReadTermios = unix.TCGETS
const ioctlWriteTermios = unix.TCSETS

var errorPasswordsDontMatch = errors.New("sorry, passwords do not match")

// FIXME: this function ignores Ctrl-D
func readPasswordStdin() (string, error) {
	fd := int(os.Stdin.Fd())
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return "", err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)
		os.Exit(1)
	}()

	fmt.Print("Type password: ")
	pwd1, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	fmt.Print("Retype password: ")
	pwd2, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	if !slices.Equal(pwd1, pwd2) {
		return "", errorPasswordsDontMatch
	}

	return string(pwd1), nil
}
