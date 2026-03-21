package view

import "fmt"

// TODO: don't use fractional decimals for bytes
func bytesToHuman(bytes int64) string {
	b := float64(bytes)
	sizes := []string{"bytes", "KiB", "MiB", "GiB", "TiB", "PiB"}
	i := 0
	for b >= 1024 && i < len(sizes)-1 {
		i++
		b /= 1024
	}
	return fmt.Sprintf("%.2f %s", b, sizes[i])
}
