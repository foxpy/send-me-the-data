package main

import (
	"bufio"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"os"
	"regexp"
	"strings"
)

var (
	matchLowercase       = regexp.MustCompile(`[a-z]`)
	matchUppercase       = regexp.MustCompile(`[A-Z]`)
	matchNumeric         = regexp.MustCompile(`\d`)
	specialCharacters    = `_+=%*&^$/\|.,:!(){}[]~`
	errorNotPasswordFile = errors.New("the specified file does not contain a password")
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
		return "", errorNotPasswordFile
	}

	return string(line), nil
}
