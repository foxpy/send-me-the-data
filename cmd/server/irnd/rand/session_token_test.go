package rand

import (
	"regexp"
	"testing"
)

func TestSessionToken(t *testing.T) {
	r := NewRandom()

	for range 10 {
		token := r.SessionToken()

		matched, err := regexp.MatchString(`^[0-9a-f]{64}$`, token)
		if err != nil {
			t.Fatal(err)
		}

		if !matched {
			t.Fatalf("session token %s does not match the expected format", token)
		}
	}
}
