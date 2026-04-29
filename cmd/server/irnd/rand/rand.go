package rand

import "github.com/foxpy/send-me-the-data/cmd/server/irnd"

type Random struct {
	alphabet []byte
}

var _ irnd.Random = &Random{}

func NewRandom() *Random {
	r := &Random{}

	for i := byte('a'); i <= byte('z'); i++ {
		r.alphabet = append(r.alphabet, i)
	}
	for i := byte('A'); i <= byte('Z'); i++ {
		r.alphabet = append(r.alphabet, i)
	}
	for i := byte('0'); i <= byte('9'); i++ {
		r.alphabet = append(r.alphabet, i)
	}

	return r
}
