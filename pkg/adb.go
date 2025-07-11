package pkg

import (
	"bytes"
	"fmt"
	"kubik-rental/dto"
	"os/exec"
	"runtime"
	"strings"
)

func ExecuteAdbAction(ip, action string) error {
	var cmd *exec.Cmd

	switch action {
	case dto.ActionOff:
		cmd = exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "26")
	case dto.ActionOn:
		cmd = exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "224")
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

func ConnectAdbToIP(ip string) bool {
	cmd := exec.Command("embed/platform-tools/adb.exe", "connect", ip+":5555")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ADB connect error:", err)
		return false
	}

	outStr := string(output)
	fmt.Println("ADB connect output:", outStr)

	// ❌ Jika output mengandung indikasi gagal
	if strings.Contains(outStr, "failed") || strings.Contains(outStr, "refused") {
		return false
	}

	// ✅ Jika sudah terhubung
	if strings.Contains(outStr, "connected to") || strings.Contains(outStr, "already connected") {
		return true
	}

	// Default: anggap berhasil jika tidak ada kata "failed" atau "refused"
	return true
}

func PingIP(ip string) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", ip)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return false
	}

	output := out.String()

	if runtime.GOOS == "windows" {
		// Cek TTL untuk Windows
		return strings.Contains(output, "TTL=")
	} else {
		// Cek bytes from untuk Linux/Mac
		return strings.Contains(output, "bytes from")
	}
}
