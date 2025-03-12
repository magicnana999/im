package utils

import (
	"github.com/denisbrodbeck/machineid"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func MacID() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, i := range interfaces {
		if len(i.HardwareAddr) > 0 {
			return i.HardwareAddr.String()
		}
	}
	return ""
}

func OsID() string {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("cat", "/etc/machine-id")
	case "windows":
		cmd = exec.Command("wmic", "csproduct", "get", "UUID")
	case "darwin":
		cmd = exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func MachineID() string {
	id, _ := machineid.ID()
	return id
}
