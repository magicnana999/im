package id

import (
	"fmt"
	"testing"
)

func TestGenerateXId(t *testing.T) {
	fmt.Println(GenerateXId())
	fmt.Println(GenerateXId())
}

func TestSnowflakeID(t *testing.T) {
	fmt.Println(SnowflakeID())
	fmt.Println(SnowflakeID())
	fmt.Println(SnowflakeID())
}
