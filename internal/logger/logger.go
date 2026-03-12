package logger

import (
	"encoding/csv"
	"fmt"
	"os"
)

// WriteCSV takes the sample data and writes it to the specified file path.
func WriteCSV(fileName string, cpuSamples []float64, memSamples []float64) error {
	if fileName == "" || len(cpuSamples) == 0 {
		return nil // Nothing to log
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"SampleIndex", "CPUPercent", "MemoryMB"}); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	for i := range cpuSamples {
		mem := 0.0
		if i < len(memSamples) {
			mem = memSamples[i]
		}

		err := writer.Write([]string{
			fmt.Sprint(i),
			fmt.Sprintf("%.2f", cpuSamples[i]),
			fmt.Sprintf("%.2f", mem),
		})

		if err != nil {
			return fmt.Errorf("error writing CSV row: %w", err)
		}
	}

	fmt.Printf("Logged data to %s\n", fileName)
	return nil
}