package admin

import (
	"regexp"
	"testing"
)

func TestGenerateRandomExternalKey(t *testing.T) {
	for range 10 {
		externalKey := generateRandomExternalKey()

		matched, err := regexp.MatchString(`^[a-zA-Z0-9]{12}$`, externalKey)
		if err != nil {
			t.Fatal(err)
		}

		if !matched {
			t.Fatalf("external key %s does not match the expected format", externalKey)
		}
	}
}

func BenchmarkGenerateRandomExternalKey(b *testing.B) {
	for b.Loop() {
		generateRandomExternalKey()
	}
}

func BenchmarkGenerateRandomExternalKeyParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			generateRandomExternalKey()
		}
	})
}
