package schedulerconfig

import (
	"net"
	"os"
)

func GetServerAddress() string {
	mode := os.Getenv("APP_MODE") // "debug" or "prod"
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	if mode == "debug" {
		return "localhost:" + port
	}

	// production / container / k8s
	return ":" + port // binds to all interfaces
}

func GetAdvertisedLeaderAddress() string {
	mode := os.Getenv("APP_MODE")
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	if mode == "debug" {
		return "http://localhost:" + port
	}

	// In production you usually want the machine IP or pod IP
	ip := os.Getenv("POD_IP") // Kubernetes injects this
	if ip == "" {
		ip = getLocalIP() // fallback
	}

	return "http://" + ip + ":" + port
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		// We only want IPv4, non-loopback
		if ip == nil || ip.IsLoopback() {
			continue
		}

		ip = ip.To4()
		if ip == nil {
			continue // not IPv4
		}

		return ip.String()
	}

	// Fallback if nothing found
	return "127.0.0.1"
}
