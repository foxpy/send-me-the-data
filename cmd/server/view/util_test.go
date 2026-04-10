package view

import (
	"testing"
	"time"
)

var mockTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
var mockTimeMilli = uint64(mockTime.UTC().UnixMilli())

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
		{
			name:  "fractional",
			bytes: 1024 + 512,
			human: "1.50 KiB",
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
