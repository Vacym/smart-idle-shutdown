package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	intervalPtr := flag.Int("interval", 5, "Check interval in seconds")
	thresholdPtr := flag.Float64("threshold", 30.0, "Load Threshold in Percent")
	consecutiveThresholdPtr := flag.Int("consecutive", 3, "Number of consecutive times the load should be below the threshold")

	flag.Parse()

	a := app.New()
	w := a.NewWindow("Smart Idle Shutdown")

	// Creating interface elements

	deviceSelector := widget.NewSelect([]string{"CPU"}, func(s string) {
		// Processing a change in the selected device
	})
	deviceSelector.SetSelected("CPU")

	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(fmt.Sprintf("%d", *intervalPtr))
	intervalEntry.Validator = ValidateIntervalEntry

	thresholdEntry := widget.NewEntry()
	thresholdEntry.SetText(fmt.Sprintf("%.2f", *thresholdPtr))
	thresholdEntry.Validator = ValidateThresholdEntry

	consecutiveEntry := widget.NewEntry()
	consecutiveEntry.SetText(fmt.Sprintf("%d", *consecutiveThresholdPtr))
	consecutiveEntry.Validator = ValidateConsecutiveEntry

	var startAction func()
	var stopAction func()

	// Placement of interface elements

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Device", Widget: deviceSelector, HintText: "Which device should be checked"},
			{Text: "Interval", Widget: intervalEntry, HintText: "Check interval in seconds"},
			{Text: "Threshold", Widget: thresholdEntry, HintText: "Load threshold in percents"},
			{Text: "Consecutive Threshold", Widget: consecutiveEntry, HintText: "Number of consecutive times the load should be below the threshold"},
		},
		SubmitText: "Start monitoring",
		CancelText: "Close the app to stop monitoring",
	}

	form.OnSubmit = func() {
		fmt.Println("Struct:")
		form.OnSubmit = nil
		form.OnCancel = func() {
			fmt.Println("Struct:!!!")
		}
		form.Refresh()
	}

	startAction = func() {
		// "Start/Cancel" button processing
		interval, err := strconv.Atoi(intervalEntry.Text)
		if err != nil || interval <= 0 {
			fmt.Println("Interval is incorrect")
			return
		}

		threshold, err := strconv.ParseFloat(thresholdEntry.Text, 64)
		if err != nil || threshold < 0 || threshold > 100 {
			fmt.Println("Threshold is incorrect")
			return
		}

		consecutiveThreshold, err := strconv.Atoi(consecutiveEntry.Text)
		if err != nil || consecutiveThreshold <= 0 {
			fmt.Println("Consecutive threshold is incorrect")
			return
		}

		form.OnSubmit = nil
		form.OnCancel = stopAction
		form.Refresh()

		monitoring(interval, threshold, consecutiveThreshold, deviceSelector.Selected)
	}

	stopAction = func() {
		fmt.Println("stop signal")
	}

	form.OnSubmit = startAction
	form.Refresh()

	mainContainer := container.NewVBox(
		form,
	)

	w.SetContent(mainContainer)

	// Launching the application

	w.ShowAndRun()
}

func monitoring(intervalSec int, threshold float64, consecutiveThreshold int, deviceSelector string) {
	fmt.Printf("Interval: %d seconds, Threshold: %.2f, Consecutive Threshold: %d, Device Selector: %s\n",
		intervalSec, threshold, consecutiveThreshold, deviceSelector)

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
