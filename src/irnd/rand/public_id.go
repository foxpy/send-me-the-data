package rand

import "math/rand/v2"

func (r *Random) PublicID() string {
	var result [12]byte
	for i := range 12 {
		n := rand.IntN(len(r.alphabet))
		result[i] = r.alphabet[n]
	}
	return string(result[:])
}
