package util

import (
	"fmt"
	"testing"
)

func TestGetMacMachineID(t *testing.T) {
	fmt.Println(GetMacMachineID())
}

func TestGetOSMachineID(t *testing.T) {
	fmt.Println(GetOSMachineID())
}

func TestGetMachineID(t *testing.T) {
	fmt.Println(GetMachineId())
}
