package util

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"sync"
	"time"
)

var (
	rng *rand.Rand
	mu  sync.Mutex
)

func getRnd() (*rand.Rand, error) {

	if rng == nil {
		mu.Lock()
		defer mu.Unlock()
		rng = rand.New(mt19937.New())
		rng.Seed(time.Now().UnixNano())
		return rng, nil
	}

	if rng == nil {
		return nil, fmt.Errorf("rand could not be initialized")
	}

	return rng, nil
}

func GetGlobalRand() (*rand.Rand, error) {
	return getRnd()
}

func GetGlobalRandomUint64() (uint64, error) {
	r, err := getRnd()
	if err != nil {
		return 0, fmt.Errorf("could not get rand instance")
	}
	return r.Uint64(), nil
}

func GetGlobalRandomUint32() (uint32, error) {
	r, err := getRnd()
	if err != nil {
		return 0, fmt.Errorf("could not get rand instance")
	}
	return r.Uint32(), nil
}
func GetGlobalRandomUint16() (uint16, error) {
	r, err := getRnd()
	if err != nil {
		return 0, fmt.Errorf("could not get rand instance")
	}
	return uint16(r.Int31n(65536)), nil
}

func GetGlobalRandomUint8() (uint8, error) {
	r, err := getRnd()
	if err != nil {
		return 0, fmt.Errorf("could not get rand instance")
	}
	return uint8(r.Int31n(256)), nil
}
func GetRandomUint64() (uint64, error) {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return r.Uint64(), nil
}

func GetRandomUint32() (uint32, error) {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return r.Uint32(), nil
}

func GetRandomUint16() (uint16, error) {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return uint16(r.Int63n(65536)), nil
}

func GetRandomUint8() (uint8, error) {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return uint8(r.Int63n(256)), nil
}
