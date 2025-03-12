package utils

import (
	"net"
)

func LocalIP() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addr {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return ""
}
