package utils

import (
	"net"
)

// GetAllLocalIPs returns all IPv4 addresses of non-loopback interfaces
func GetAllLocalIPs() []string {
	var ips []string
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips
}
