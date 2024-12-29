package util

import (
	"fmt"
	"testing"
)

func TestGetGlobalRand(t *testing.T) {
	r, _ := getRnd()
	fmt.Println(r.Int63())
	fmt.Println(r.Int31())
	fmt.Println(r.Int())
}

func TestGetGlobalRandomUint64(t *testing.T) {
	r, _ := GetGlobalRandomUint64()
	fmt.Println(r)
}
func TestGetGlobalRandomUint32(t *testing.T) {
	r, _ := GetGlobalRandomUint32()
	fmt.Println(r)
}

func TestGetGlobalRandomUint16(t *testing.T) {
	r, _ := GetGlobalRandomUint16()
	fmt.Println(r)
}

func TestGetGlobalRandomUint8(t *testing.T) {
	r, _ := GetGlobalRandomUint8()
	fmt.Println(r)
}

func TestGetRandoUint64(t *testing.T) {
	r, _ := GetRandomUint64()
	fmt.Println(r)
}

func TestGetRandomUint32(t *testing.T) {
	r, _ := GetRandomUint32()
	fmt.Println(r)
}

func TestGetRandomUint16(t *testing.T) {
	r, _ := GetRandomUint16()
	fmt.Println(r)
}

func TestGetRandomUint8(t *testing.T) {
	r, _ := GetRandomUint8()
	fmt.Println(r)
}
