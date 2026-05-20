package view

import (
	"testing"
)

func TestBytesToHuman(t *testing.T) {
	for _, tc := range []struct {
		name  string
		bytes uint64
		human string
	}{

		{
			name:  "zero",
			bytes: 0,
			human: "0 bytes",
		},
		{
			name:  "bytes",
			bytes: 100,
			human: "100 bytes",
		},
		{
			name:  "kibibytes",
			bytes: 2048,
			human: "2 KiB",
		},
		{
			name:  "kibibytes fractional",
			bytes: 2048 + 256,
			human: "2.25 KiB",
		},
		{
			name:  "tebibytes",
			bytes: 5 * (1 << 40),
			human: "5 TiB",
		},
		{
			name:  "pebibytes",
			bytes: 7_200 * (1 << 50),
			human: "7200 PiB",
		},
		{
			name:  "pebibytes fractional",
			bytes: 7_200*(1<<50) + 512*(1<<40),
			human: "7200.50 PiB",
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
