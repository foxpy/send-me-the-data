package view

import (
	"testing"

	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
)

type linkFiles struct {
	name  string
	files []ifs.File
}

func TestBytesToHuman(t *testing.T) {
	for _, tc := range []struct {
		name  string
		bytes int64
		human string
	}{

		{
			name:  "zero",
			bytes: 0,
			human: "0.00 bytes",
		},
		{
			name:  "bytes",
			bytes: 100,
			human: "100.00 bytes",
		},
		{
			name:  "kibibytes",
			bytes: 2048,
			human: "2.00 KiB",
		},
		{
			name:  "tebibytes",
			bytes: 5 * (1 << 40),
			human: "5.00 TiB",
		},
		{
			name:  "pebibytes",
			bytes: 7_200 * (1 << 50),
			human: "7200.00 PiB",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			human := bytesToHuman(tc.bytes)
			if human != tc.human {
				t.Fatalf("expected %s, got %s", tc.human, human)
			}
		})
	}
}
