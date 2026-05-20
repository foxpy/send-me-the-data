package rand

import (
	"crypto/rand"
	"encoding/hex"
)

func (r *Random) SessionToken() string {
	buf := make([]byte, 32)
	// this function never fails and always fills the entire buffer
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
