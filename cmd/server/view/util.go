package view

import "fmt"

func bytesToHuman(bytes uint64) string {
	b := float64(bytes)
	sizes := []struct {
		name, format string
	}{
		{"bytes", "%.0f %s"},
		{"KiB", "%.2f %s"}, // TODO: I would like to omit .00 if the value is an integer
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
