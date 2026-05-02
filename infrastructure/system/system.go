package system

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }

type RandomIDGenerator struct{}

func (RandomIDGenerator) NewID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
