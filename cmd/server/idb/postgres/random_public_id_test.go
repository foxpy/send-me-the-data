package postgres

import (
	"regexp"
	"testing"
)

func TestGenerateRandomPublicID(t *testing.T) {
	p := Postgres{}

	for range 10 {
		id := p.GenerateRandomPublicID()

		matched, err := regexp.MatchString(`^[a-zA-Z0-9]{12}$`, id)
		if err != nil {
			t.Fatal(err)
		}

		if !matched {
			t.Fatalf("public ID %s does not match the expected format", id)
		}
	}
}

func BenchmarkGenerateRandomPublicID(b *testing.B) {
	p := Postgres{}

	for b.Loop() {
		p.GenerateRandomPublicID()
	}
}

func BenchmarkGenerateRandomPublicIDParallel(b *testing.B) {
	p := Postgres{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.GenerateRandomPublicID()
		}
	})
}
