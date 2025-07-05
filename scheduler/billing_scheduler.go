package scheduler

import (
	"fmt"
	"os/exec"
	"time"

	"kubik-rental/config"
	"kubik-rental/entity"
)

func StartBillingPolling() {
	go func() {
		for {
			processExpiredBillings()
			time.Sleep(5 * time.Second) // polling interval
		}
	}()
}

func processExpiredBillings() {
	now := time.Now()
	var billings []entity.Billing

	// Ambil billing yang waktu berakhirnya lewat & masih aktif
	if err := config.DB.Preload("Device").
		Where("end_time <= ? AND status = ?", now, "active").
		Find(&billings).Error; err != nil {
		fmt.Println("Gagal ambil billing:", err)
		return
	}

	for _, billing := range billings {
		fmt.Printf("Mematikan device: %s (%s)\n", billing.Device.Name, billing.Device.IP)

		// Connect dan matikan TV
		shutdownTV(billing.Device.IP)

		// Update status billing jadi expired
		billing.Status = "expired"
		if err := config.DB.Save(&billing).Error; err != nil {
			fmt.Println("Gagal update status billing:", err)
			continue
		}

		// Set device.available = true
		if err := config.DB.Model(&entity.Device{}).
			Where("id = ?", billing.DeviceID).
			Update("available", true).Error; err != nil {
			fmt.Println("Gagal update device available:", err)
		} else {
			fmt.Printf("Device %s kini tersedia kembali.\n", billing.Device.Name)
		}
	}
}

func shutdownTV(ip string) {
	exec.Command("embed/platform-tools/adb.exe", "connect", ip).Run()
	cmd := exec.Command("embed/platform-tools/adb.exe", "-s", ip, "shell", "input", "keyevent", "26") // Tombol power
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Gagal matikan TV %s: %v\n", ip, err)
	} else {
		fmt.Printf("TV %s dimatikan. Output: %s\n", ip, string(out))
	}
}
