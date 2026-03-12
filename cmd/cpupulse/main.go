package main

import (
	"flag"
	"fmt"
	// "net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
	"cpuPulse/internal/monitor"
	"cpuPulse/internal/logger"
	"cpuPulse/internal/plotter"

	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

// var (
//     cpuGauge = prometheus.NewGauge(prometheus.GaugeOpts{
//         Name: "cpupulse_cpu_percent",
//         Help: "CPU usage percent of monitored process",
//     })
//     memGauge = prometheus.NewGauge(prometheus.GaugeOpts{
//         Name: "cpupulse_memory_mb",
//         Help: "Memory RSS in MB of monitored process",
//     })
// )

var (
	startTime    time.Time
	logFileName  string
	enablePlot   bool
	// samples      []float64
	// peakMemoryMB float64 = 0
	// memSamples   []float64
)

func init() {
	flag.StringVar(&logFileName, "log", "", "Log CPU and memory usage to a file (CSV format)")
	flag.BoolVar(&enablePlot, "plot", false, "Generate a CPU usage vs time plot")
	//flag.Parse()
	// prometheus.MustRegister(cpuGauge)
    // prometheus.MustRegister(memGauge)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: cpupulse [--log filename] [--plot] <command> [args...]")
		os.Exit(1)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime = time.Now()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start process: %v\n", err)
		os.Exit(1)
	}

	if cmd.Process == nil {
		fmt.Println("Error: could not retrieve process PID")
		os.Exit(1)
	}

	pid := cmd.Process.Pid

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	sysMonitor := monitor.NewMonitor(pid)


	go func() {
		// monitor.MonitorProcess(pid, done)
		sysMonitor.Start(done)
	}()

	// go func() {
	// 	http.Handle("/metrics", promhttp.Handler())
	// 	fmt.Println("Prometheus metrics at :2112/metrics")
	// 	http.ListenAndServe(":2112", nil)
	// }()


	select {
	case <-interruptChan:
		fmt.Println("\n[!] Interrupted. Killing process...")
		_ = cmd.Process.Kill()
		<-done
		// monitor.PrintStats(time.Since(startTime))
		// monitor.LogStats(logFileName)
		// monitor.PlotStats(enablePlot)
		sysMonitor.Print(time.Since(startTime))
        // sysMonitor.Log(logFileName)
		if logFileName != "" {
			if err := logger.WriteCSV(logFileName, sysMonitor.CPUSamples, sysMonitor.MemSamples); err != nil {
				fmt.Printf("Logging error: %v\n", err)
			}
		}
        // sysMonitor.Plot(enablePlot)
		if enablePlot {
			if err := plotter.GeneratePNG("cpupulse_plot.png", sysMonitor.CPUSamples, sysMonitor.MemSamples); err != nil {
				fmt.Printf("Plotting error: %v\n", err)
			}
		}
		os.Exit(1)

	case err = <-waitProcess(cmd):
		<-done
		if err != nil {
			fmt.Printf("Process exited with error: %v\n", err)
		}
		// monitor.PrintStats(time.Since(startTime))
		// monitor.LogStats(logFileName)
		// monitor.PlotStats(enablePlot)
		sysMonitor.Print(time.Since(startTime))
        // sysMonitor.Log(logFileName)
		if logFileName != "" {
			if err := logger.WriteCSV(logFileName, sysMonitor.CPUSamples, sysMonitor.MemSamples); err != nil {
				fmt.Printf("Logging error: %v\n", err)
			}
		}
        // sysMonitor.Plot(enablePlot)
		if enablePlot {
			if err := plotter.GeneratePNG("cpupulse_plot.png", sysMonitor.CPUSamples, sysMonitor.MemSamples); err != nil {
				fmt.Printf("Plotting error: %v\n", err)
			}
		}
	}
}

func waitProcess(cmd *exec.Cmd) chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- cmd.Wait()
	}()
	return ch
}
