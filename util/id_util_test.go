package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateXId(t *testing.T) {
	id := GenerateXId()
	fmt.Println(strings.ToUpper(id))
}

func TestGenerateSonyFlakeId(t *testing.T) {
	r, _ := GenerateSonyFlakeId()
	fmt.Println(r)
}
