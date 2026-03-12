package monitor

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"time"
)


type SystemStats struct {
    PID          int
    CPUSamples   []float64
    MemSamples   []float64
    PeakMemoryMB float64
}

func NewMonitor(pid int) *SystemStats {
	return &SystemStats{
		PID:          pid,
		CPUSamples:   make([]float64, 0),
		MemSamples:   make([]float64, 0),
		PeakMemoryMB: 0,
	}
}

func (s *SystemStats) Start(done chan struct{}) {
	defer close(done)

	if s.PID <= 0 {
		fmt.Println("Invalid PID provided to monitor.")
		return
	}

	p, err := process.NewProcess(int32(s.PID))
	if err != nil {
		fmt.Println("Error: could not attach to process:", err)
		return
	}

	for {
		exists, err := p.IsRunning()
		if err != nil || !exists {
			break
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			break 
		}
		s.CPUSamples = append(s.CPUSamples, cpuPercent)

		memInfo, err := p.MemoryInfo()
		if err != nil {
			break
		}
		
		currMemMB := float64(memInfo.RSS) / 1024.0 / 1024.0
		s.MemSamples = append(s.MemSamples, currMemMB)
		if currMemMB > s.PeakMemoryMB {
			s.PeakMemoryMB = currMemMB
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (s *SystemStats) Print(duration time.Duration) {
	if len(s.CPUSamples) == 0 {
		fmt.Println("No samples recorded.")
		return
	}

	var total float64 = 0
	var maxCPU float64 = 0
	for _, cpu := range s.CPUSamples {
		total += cpu
		if cpu > maxCPU {
			maxCPU = cpu
		}
	}

	avg := total / float64(len(s.CPUSamples))

	fmt.Println("\n========== CPUPulse Report ==========")
	fmt.Printf("Duration        : %.2fs\n", duration.Seconds())
	fmt.Printf("Avg CPU Usage   : %.2f%%\n", avg)
	fmt.Printf("Peak CPU Usage  : %.2f%%\n", maxCPU)
	fmt.Printf("Peak Memory RSS : %.2f MB\n", s.PeakMemoryMB)
}