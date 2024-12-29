package util

import (
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func GetMacMachineID() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		if len(i.HardwareAddr) > 0 {
			return i.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("no MAC address found")
}

func GetOSMachineID() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("cat", "/etc/machine-id")
	case "windows":
		cmd = exec.Command("wmic", "csproduct", "get", "UUID")
	case "darwin":
		cmd = exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetMachineId() (string, error) {
	return machineid.ID()
}
