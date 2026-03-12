package plotter

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func GeneratePNG(outputFileName string, cpuSamples []float64, memSamples []float64) error {
	if len(cpuSamples) == 0 || len(memSamples) == 0 {
		return fmt.Errorf("insufficient data to plot")
	}

	p := plot.New()

	cpuPts := make(plotter.XYs, len(cpuSamples))
	memPts := make(plotter.XYs, len(memSamples))

	for i := range cpuSamples {
		cpuPts[i].X = float64(i)
		cpuPts[i].Y = cpuSamples[i]
	}
	for i := range memSamples {
		memPts[i].X = float64(i)
		memPts[i].Y = memSamples[i]
	}

	cpuLine, err := plotter.NewLine(cpuPts)
	if err != nil {
		return fmt.Errorf("failed to create CPU line: %w", err)
	}
	cpuLine.Color = plotutil.Color(0)
	cpuLine.LineStyle.Width = vg.Points(1)
	cpuLine.LineStyle.Dashes = []vg.Length{vg.Points(3), vg.Points(3)}

	memLine, err := plotter.NewLine(memPts)
	if err != nil {
		return fmt.Errorf("failed to create Memory line: %w", err)
	}
	memLine.Color = plotutil.Color(1)
	memLine.LineStyle.Width = vg.Points(1)

	p.Add(cpuLine, memLine)
	p.Legend.Add("CPU (%)", cpuLine)
	p.Legend.Add("Memory (MB)", memLine)
	p.Legend.Top = true
	p.Legend.Left = false

	if err := p.Save(8*vg.Inch, 4*vg.Inch, outputFileName); err != nil {
		return fmt.Errorf("failed to save plot: %w", err)
	}

	fmt.Printf("Saved plot to %s\n", outputFileName)
	return nil
}