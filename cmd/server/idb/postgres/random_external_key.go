package postgres

import "math/rand/v2"

var (
	alphabet []byte
)

func init() {
	for i := byte('a'); i <= byte('z'); i++ {
		alphabet = append(alphabet, i)
	}
	for i := byte('A'); i <= byte('Z'); i++ {
		alphabet = append(alphabet, i)
	}
	for i := byte('0'); i <= byte('9'); i++ {
		alphabet = append(alphabet, i)
	}
}

func (d *Postgres) GenerateRandomExternalKey() string {
	var result [12]byte
	for i := range 12 {
		n := rand.IntN(len(alphabet))
		result[i] = alphabet[n]
	}
	return string(result[:])
}
