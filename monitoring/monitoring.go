package monitoring

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type Monitor struct {
	intervalSec          int
	threshold            float64
	consecutiveThreshold int
	deviceSelector       string
	stop                 chan struct{}
}

func NewMonitor(intervalSec int, threshold float64, consecutiveThreshold int, deviceSelector string) *Monitor {
	return &Monitor{
		intervalSec:          intervalSec,
		threshold:            threshold,
		consecutiveThreshold: consecutiveThreshold,
		deviceSelector:       deviceSelector,
		stop:                 make(chan struct{}),
	}
}

func (m *Monitor) Start() {
	go func() {
		consecutiveCount := 0 // Count the number of times the load is below the threshold in a row
		interval := time.Second * time.Duration(m.intervalSec)

		for {
			select {
			case <-m.stop:
				// Stop signal received, return
				fmt.Println("Monitoring stopped.")
				return
			default:
				startTime := time.Now()

				load, err := getCPULoad(1)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(load)

					if load < m.threshold {
						consecutiveCount++
					} else {
						consecutiveCount = 0
					}

					if consecutiveCount >= m.consecutiveThreshold {
						shutdown()
						return
					}
				}

				waitTime := interval - time.Since(startTime)

				if waitTime > 0 {
					time.Sleep(waitTime)
				}
			}
		}
	}()
}

func (m *Monitor) Stop() {
	close(m.stop)
}

func getCPULoad(seconds int) (float64, error) {
	period := time.Second * time.Duration(seconds)
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
