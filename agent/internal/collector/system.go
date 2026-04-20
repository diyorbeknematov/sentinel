package collector

import (
	"fmt"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func StartSystemMetricsCollector(cfg *config.Config) {
	fmt.Println("System metrics collector ishga tushdi")

	for {
		// CPU %
		cpuPercent, _ := cpu.Percent(0, false)

		// RAM
		vm, _ := mem.VirtualMemory()

		// DISK
		diskStat, _ := disk.Usage("/")

		// LOG (keyin senderga beramiz)
		fmt.Printf("CPU: %.2f%% RAM: %.2f%% DISK: %.2f%%\n",
			cpuPercent[0],
			vm.UsedPercent,
			diskStat.UsedPercent,
		)

		fmt.Println(models.MetricPayload{
			CPU:     cpuPercent[0],
			RAM:     vm.UsedPercent,
			Disk:    diskStat.UsedPercent,
			LogTime: time.Now(),
		})
		// TODO: send to server

		time.Sleep(5 * time.Second)
	}
}
