build :
	go build -o cpupulse ./cmd/cpupulse

clean:
	rm -r cpupulse cpupulse_plot.png stats.csv