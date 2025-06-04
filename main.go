package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	startTime time.Time
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cpupulse <command> [args...]")
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime = time.Now()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start process: %v\n", err)
		os.Exit(1)
	}

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		MonitorProcess(cmd.Process.Pid, done)
	}()

	select {
	case <-interruptChan:
		fmt.Println("\n[!] Interrupted. Killing process...")
		_ = cmd.Process.Kill()
		<-done
		PrintStats(time.Since(startTime))
		os.Exit(1)

	case err = <-waitProcess(cmd):
		<-done
		if err != nil {
			fmt.Printf("Process exited with error: %v\n", err)
		}
		PrintStats(time.Since(startTime))
	}
}

func waitProcess(cmd *exec.Cmd) chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- cmd.Wait()
	}()
	return ch
}
