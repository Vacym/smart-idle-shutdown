package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func main() {
	intervalPtr := flag.Int("interval", 1, "Check interval in seconds")
	thresholdPtr := flag.Float64("threshold", 30.0, "Load Threshold in Percent")
	consecutiveThresholdPtr := flag.Int("consecutive", 1, "Number of consecutive times the load should be below the threshold")

	flag.Parse()

	monitoring(*intervalPtr, *thresholdPtr, *consecutiveThresholdPtr)
}

func monitoring(intervalSec int, threshold float64, consecutiveThreshold int) {
	consecutiveCount := 0 // Count the number of times the load is below the threshold in a row
	interval := time.Second * time.Duration(intervalSec)

	for {
		startTime := time.Now()

		load, err := getCPULoad(1)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Load: %.2f%%\n", load)

		if load < threshold {
			consecutiveCount++
		} else {
			consecutiveCount = 0
		}

		if consecutiveCount >= consecutiveThreshold {
			shutdown()
			return
		}

		waitTime := interval - time.Since(startTime)

		if waitTime > 0 {
			time.Sleep(waitTime)
		}
	}
}

func getCPULoad(periodSec int) (float64, error) {
	period := time.Second * time.Duration(periodSec)
	percentages, err := cpu.Percent(period, false)
	if err != nil {
		return 0, err
	}
	return percentages[0], nil
}

func shutdown() {
	fmt.Println("Initiating shutdown...")
	time.Sleep(10 * time.Second) // Give the user time to see the message
	cmd := exec.Command("shutdown", "/s", "/t", "1")
	cmd.Run()
	os.Exit(0)
}
