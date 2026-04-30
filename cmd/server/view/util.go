package view

import (
	"fmt"
	"math"
)

func bytesToHuman(bytes uint64) string {
	b := float64(bytes)
	names := []string{
		"bytes", "KiB", "MiB", "GiB", "TiB", "PiB",
	}

	i := 0
	for b >= 1024 && i < len(names)-1 {
		i++
		b /= 1024
	}

	format := "%.2f %s"
	if math.Floor(b) == b {
		format = "%.0f %s"
	}

	return fmt.Sprintf(format, b, names[i])
}
