package random

import (
	"github.com/seehuhn/mt19937"
	"math/rand"
	"sync"
	"time"
)

var (
	rng        *rand.Rand
	randomOnce sync.Once
)

func Random() *rand.Rand {

	randomOnce.Do(func() {
		rng = rand.New(mt19937.New())
		rng.Seed(time.Now().UnixNano())
	})

	return rng
}

func RandomUint64() uint64 {
	r := Random()
	return r.Uint64()
}

func RandomUint32() uint32 {
	r := Random()
	return r.Uint32()
}

func RandomUint16() uint16 {
	r := Random()
	return uint16(r.Int31n(65536))
}

func RandomUint8() uint8 {
	r := Random()
	return uint8(r.Int31n(256))
}
