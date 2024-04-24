package monitoring

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// Settings represents the configuration settings for the monitor.
type Settings struct {
	Interval             int
	Threshold            float64
	ConsecutiveThreshold int
	Device               string
}

// Monitor is responsible for monitoring the CPU load and initiating a shutdown if the load is below the threshold for a specified number of consecutive times.
type Monitor struct {
	interval             time.Duration
	threshold            float64
	consecutiveThreshold int
	device               string
	stop                 chan struct{}
}

// NewMonitor creates a new instance of the Monitor with the provided settings.
func NewMonitor(settings Settings) *Monitor {
	return &Monitor{
		interval:             time.Duration(settings.Interval) * time.Second,
		threshold:            settings.Threshold,
		consecutiveThreshold: settings.ConsecutiveThreshold,
		device:               settings.Device,
		stop:                 make(chan struct{}),
	}
}

// Start initiates the monitoring process in a separate goroutine.
func (m *Monitor) Start() {
	go m.monitorLoop()
}

// Stop sends a signal to stop the monitoring process.
func (m *Monitor) Stop() {
	close(m.stop)
}

func (m *Monitor) monitorLoop() {
	consecutiveCount := 0

	for {
		select {
		case <-m.stop:
			fmt.Println("Monitoring stopped.")
			return
		default:
			startTime := time.Now()

			load, err := getCPULoad(1)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("Load: %.2f%%\n", load)

			if load < m.threshold {
				consecutiveCount++
			} else {
				consecutiveCount = 0
			}

			if consecutiveCount >= m.consecutiveThreshold {
				initiateShutdown()
				return
			}

			waitTime := m.interval - time.Since(startTime)
			if waitTime > 0 {
				time.Sleep(waitTime)
			}
		}
	}
}

func getCPULoad(seconds int) (float64, error) {
	period := time.Duration(seconds) * time.Second
	percentages, err := cpu.Percent(period, false)
	if err != nil {
		return 0, err
	}
	return percentages[0], nil
}

func initiateShutdown() {
	fmt.Println("Initiating shutdown...")
	time.Sleep(10 * time.Second) // Give the user time to see the message
	cmd := exec.Command("shutdown", "/s", "/t", "1")
	cmd.Run()
	os.Exit(0)
}
