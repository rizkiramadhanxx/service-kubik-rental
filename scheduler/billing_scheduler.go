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
			// fmt.Println("Memeriksa billing yang expired...")
			ProcessExpiredBillings()
			time.Sleep(5 * time.Second) // polling interval
		}
	}()
}

func ProcessExpiredBillings() {
	now := time.Now()
	var billings []entity.Billing

	// Ambil billing yang waktu berakhirnya lewat & masih aktif
	if err := config.DB.Preload("Device").
		Where("end_time <= ? AND is_active = ?", now, true).
		Find(&billings).Error; err != nil {
		fmt.Println("Gagal ambil billing:", err)
		return
	}

	// ubah is_active billing jadi false
	for _, billing := range billings {
		if err := config.DB.Save(&billing).Error; err != nil {
			fmt.Println("Gagal update status billing:", err)
			continue
		}
	}

	for _, billing := range billings {

		if billing.IsLoss {
			continue
		}
		fmt.Printf("Mematikan device: %s (%s)\n ID: %d\n Status: %v\n", billing.Device.Name, billing.Device.IP, billing.ID, billing.IsActive)

		billing.IsActive = false

		// Connect dan matikan TV
		shutdownTV(billing.Device.IP)

		// Update status billing jadi expired
		if err := config.DB.Save(&billing).Error; err != nil {
			fmt.Println("Gagal update status billing:", err)
			continue
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
