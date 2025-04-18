package random

import (
	"github.com/seehuhn/mt19937"
	"math/rand"
	"time"
)

func Random() *rand.Rand {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	return rng
}
