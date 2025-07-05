package pkg

import (
	"os/exec"
	"strings"
)

func CheckAdbConnected(ip string) bool {
	cmd := exec.Command("embed/platform-tools/adb.exe", "devices")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Cari IP:5555 di hasil adb devices
	return strings.Contains(string(output), ip+":5555")
}

func PingIP(ip string) bool {
	cmd := exec.Command("ping", "-n", "1", ip) // Untuk Windows
	// Untuk Linux/Mac pakai: exec.Command("ping", "-c", "1", ip)

	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
