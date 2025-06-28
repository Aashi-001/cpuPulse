//var (
//	samples      []float64
//	peakMemoryMB float64 = 0
//)

//func MonitorProcess(pid int, done chan struct{}) {
//	p, err := process.NewProcess(int32(pid))
//	if err != nil {
//		fmt.Println("Error: could not attach to process:", err)
//		close(done)
//		return
//	}
//
//	for {
//		exists, err := p.IsRunning()
//		if err != nil || !exists {
//			break
//		}
//
//		cpuPercent, _ := p.CPUPercent()
//		samples = append(samples, cpuPercent)
//
//		memInfo, _ := p.MemoryInfo()
//		currMemMB := float64(memInfo.RSS) / 1024.0 / 1024.0
//		if currMemMB > peakMemoryMB {
//			peakMemoryMB = currMemMB
//		}
//
//		time.Sleep(100 * time.Millisecond)
//	}
//	close(done)
//}
//
//func PrintStats(duration time.Duration) {
//	var total float64 = 0
//	var maxCPU float64 = 0
//	for _, s := range samples {
//		total += s
//		if s > maxCPU {
//			maxCPU = s
//		}
//	}
//
//	avg := 0.0
//	if len(samples) > 0 {
//		avg = total / float64(len(samples))
//	}
//
//	fmt.Println("========== CPUPulse Report ==========")
//	fmt.Printf("Duration        : %.2fs\n", duration.Seconds())
//	fmt.Printf("Avg CPU Usage   : %.2f%%\n", avg)
//	fmt.Printf("Peak CPU Usage  : %.2f%%\n", maxCPU)
//	fmt.Printf("Peak Memory RSS : %.2f MB\n", peakMemoryMB)
//}

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"os"
	"time"
)

func MonitorProcess(pid int, done chan struct{}) {
	if pid <= 0 {
		fmt.Println("Invalid PID provided to monitor.")
		close(done)
		return
	}

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
		memSamples = append(memSamples, currMemMB)
		if currMemMB > peakMemoryMB {
			peakMemoryMB = currMemMB
		}

		time.Sleep(10 * time.Millisecond)

		// cpuGauge.Set(cpuPercent)
		// memGauge.Set(currMemMB)
	}
	close(done)
}

func PrintStats(duration time.Duration) {
	if len(samples) == 0 {
		fmt.Println("No samples recorded.")
		return
	}

	var total float64 = 0
	var maxCPU float64 = 0
	for _, s := range samples {
		total += s
		if s > maxCPU {
			maxCPU = s
		}
	}

	avg := total / float64(len(samples))

	fmt.Println("========== CPUPulse Report ==========")
	fmt.Printf("Duration        : %.2fs\n", duration.Seconds())
	fmt.Printf("Avg CPU Usage   : %.2f%%\n", avg)
	fmt.Printf("Peak CPU Usage  : %.2f%%\n", maxCPU)
	fmt.Printf("Peak Memory RSS : %.2f MB\n", peakMemoryMB)
}

func LogStats() {
	if logFileName == "" || len(samples) == 0 {
		return
	}

	file, err := os.Create(logFileName)
	if err != nil {
		fmt.Println("Failed to create log file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"SampleIndex", "CPUPercent", "MemoryMB"})
	for i := range samples {
		mem := 0.0
		if i < len(memSamples) {
			mem = memSamples[i]
		}
		writer.Write([]string{
			fmt.Sprint(i),
			fmt.Sprintf("%.2f", samples[i]),
			fmt.Sprintf("%.2f", mem),
		})
	}
	fmt.Printf("Logged data to %s\n", logFileName)
}

func PlotStats() {
	if !enablePlot || len(samples) == 0 || len(memSamples) == 0 {
		return
	}

	p := plot.New()
	//p.Title.Text = "CPU Usage Over Time"
	//p.X.Label.Text = "Sample"
	//p.Y.Label.Text = "CPU (%)"
	//
	//pts := make(plotter.XYs, len(samples))
	//for i := range samples {
	//	pts[i].X = float64(i)
	//	pts[i].Y = samples[i]
	//}
	//
	//line, err := plotter.NewLine(pts)
	//if err != nil {
	//	fmt.Println("Failed to create plot line:", err)
	//	return
	//}
	//
	//p.Add(line)
	//if err := p.Save(6*vg.Inch, 4*vg.Inch, "cpupulse_plot.png"); err != nil {
	//	fmt.Println("Failed to save plot:", err)
	//	return
	//}
	//fmt.Println("Saved plot to cpupulse_plot.png")
	cpuPts := make(plotter.XYs, len(samples))
	memPts := make(plotter.XYs, len(memSamples))

	for i := range samples {
		cpuPts[i].X = float64(i)
		cpuPts[i].Y = samples[i]
	}
	for i := range memSamples {
		memPts[i].X = float64(i)
		memPts[i].Y = memSamples[i]
	}

	cpuLine, err := plotter.NewLine(cpuPts)
	if err != nil {
		fmt.Println("Failed to create CPU line:", err)
		return
	}
	cpuLine.Color = plotutil.Color(0)
	cpuLine.LineStyle.Width = vg.Points(1)
	cpuLine.LineStyle.Dashes = []vg.Length{vg.Points(3), vg.Points(3)}

	memLine, err := plotter.NewLine(memPts)
	if err != nil {
		fmt.Println("Failed to create Memory line:", err)
		return
	}
	memLine.Color = plotutil.Color(1)
	memLine.LineStyle.Width = vg.Points(1)

	p.Add(cpuLine, memLine)
	p.Legend.Add("CPU (%)", cpuLine)
	p.Legend.Add("Memory (MB)", memLine)
	p.Legend.Top = true
	p.Legend.Left = false
	//p.Legend.XOff = 0
	//p.Legend.YOff = 0

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "cpupulse_plot.png"); err != nil {
		fmt.Println("Failed to save plot:", err)
		return
	}
	fmt.Println("Saved plot to cpupulse_plot.png")
}
