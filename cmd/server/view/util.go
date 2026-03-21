package view

import "fmt"

// TODO: maybe I will switch to CSR for timestamps, showing time in user's local time zone
const DateTimeFormat = "Jan 2 15:04:05 MST 2006"

func bytesToHuman(bytes int64) string {
	b := float64(bytes)
	sizes := []struct {
		name, format string
	}{
		{"bytes", "%.0f %s"},
		{"KiB", "%.2f %s"},
		{"MiB", "%.2f %s"},
		{"GiB", "%.2f %s"},
		{"TiB", "%.2f %s"},
		{"PiB", "%.2f %s"},
	}

	i := 0
	for b >= 1024 && i < len(sizes)-1 {
		i++
		b /= 1024
	}
	return fmt.Sprintf(sizes[i].format, b, sizes[i].name)
}
