package pkg

import (
	"fmt"
	"kubik-rental/dto"
	"os/exec"
	"strings"
)

func ExecuteAdbAction(ip, action string) error {
	var cmd *exec.Cmd

	switch action {
	case dto.ActionOff, dto.ActionOn:
		cmd = exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "26")
	case dto.ActionVolumeUp:
		cmd = exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "24")
	case dto.ActionVolumeDown:
		cmd = exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "25")
	case dto.ActionMute:
		cmd = exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "164")
	case dto.ActionHome:
		cmd = exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "3")
	default:
		return fmt.Errorf("unsupported action")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ADB error: %s - %s", err.Error(), string(output))
	}
	return nil
}

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
