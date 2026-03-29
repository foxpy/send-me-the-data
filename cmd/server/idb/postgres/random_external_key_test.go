package postgres

import (
	"regexp"
	"testing"
)

func TestGenerateRandomExternalKey(t *testing.T) {
	p := Postgres{}

	for range 10 {
		externalKey := p.GenerateRandomExternalKey()

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
	p := Postgres{}

	for b.Loop() {
		p.GenerateRandomExternalKey()
	}
}

func BenchmarkGenerateRandomExternalKeyParallel(b *testing.B) {
	p := Postgres{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.GenerateRandomExternalKey()
		}
	})
}
