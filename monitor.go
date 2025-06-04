package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"time"
)

var (
	samples      []float64
	peakMemoryMB float64 = 0
)

func MonitorProcess(pid int, done chan struct{}) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("Error: could not attach to process:", err)
		close(done)
		return
	}

	for {
		exists, err := p.IsRunning()
		if err != nil || !exists {
			break
		}

		cpuPercent, _ := p.CPUPercent()
		samples = append(samples, cpuPercent)

		memInfo, _ := p.MemoryInfo()
		currMemMB := float64(memInfo.RSS) / 1024.0 / 1024.0
		if currMemMB > peakMemoryMB {
			peakMemoryMB = currMemMB
		}

		time.Sleep(100 * time.Millisecond)
	}
	close(done)
}

func PrintStats(duration time.Duration) {
	var total float64 = 0
	var maxCPU float64 = 0
	for _, s := range samples {
		total += s
		if s > maxCPU {
			maxCPU = s
		}
	}

	avg := 0.0
	if len(samples) > 0 {
		avg = total / float64(len(samples))
	}

	fmt.Println("========== CPUPulse Report ==========")
	fmt.Printf("Duration        : %.2fs\n", duration.Seconds())
	fmt.Printf("Avg CPU Usage   : %.2f%%\n", avg)
	fmt.Printf("Peak CPU Usage  : %.2f%%\n", maxCPU)
	fmt.Printf("Peak Memory RSS : %.2f MB\n", peakMemoryMB)
}
