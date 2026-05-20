package rand

import (
	"regexp"
	"testing"
)

func TestPublicID(t *testing.T) {
	r := NewRandom()

	for range 10 {
		id := r.PublicID()

		matched, err := regexp.MatchString(`^[a-zA-Z0-9]{12}$`, id)
		if err != nil {
			t.Fatal(err)
		}

		if !matched {
			t.Fatalf("public ID %s does not match the expected format", id)
		}
	}
}

func BenchmarkPublicID(b *testing.B) {
	r := NewRandom()

	for b.Loop() {
		r.PublicID()
	}
}

func BenchmarkParallelPublicID(b *testing.B) {
	r := NewRandom()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.PublicID()
		}
	})
}
