package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"go-chat-cli/client"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("🔑 Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	for {
		defaultIP := getLocalIPv4()
		ip := promptIPWithDefault(defaultIP, reader)

		wsURL := fmt.Sprintf("ws://%s:8080/ws?username=%s", ip, username)

		fmt.Println("🔌 Connecting to", wsURL)

		// Thử kết nối WebSocket (timeout nhanh)
		success := testConnection(wsURL)
		if success {
			client.StartClientUI(wsURL, username)
			break
		} else {
			fmt.Println("❌ Cannot connect to server. Please check IP and try again.\n")
		}
	}
}

// Gợi ý IP mặc định → user có thể sửa
func promptIPWithDefault(defaultIP string, reader *bufio.Reader) string {
	fmt.Printf("🌐 Enter server IP [%s]: ", defaultIP)
	ipInput, _ := reader.ReadString('\n')
	ipInput = strings.TrimSpace(ipInput)
	if ipInput == "" {
		return defaultIP
	}
	return ipInput
}

// Get first non-loopback IPv4 address
func getLocalIPv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}

	for _, iface := range ifaces {
		// Skip loopback, non-active, or non-real network adapters
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		name := strings.ToLower(iface.Name)
		if strings.Contains(name, "veth") || strings.Contains(name, "loopback") || strings.Contains(name, "virtual") || strings.Contains(name, "docker") || strings.Contains(name, "default switch") {
			continue
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // IPv6 → skip
			}

			// ✅ This is a valid IPv4 address on a real network adapter
			return ip.String()
		}
	}
	return "127.0.0.1"
}


// Test connection to WebSocket server quickly
func testConnection(wsURL string) bool {
	conn, _, err := client.DialWebSocket(wsURL, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
