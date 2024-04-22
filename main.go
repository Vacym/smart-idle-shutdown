package main

import (
	"flag"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Vacym/smart-idle-shutdown/monitoring"
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
		CancelText: "Stop monitoring",
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

		newMonitor := startNewMonitoring(interval, threshold, consecutiveThreshold, deviceSelector.Selected)

		stopAction = func() {
			stopMonitoring(newMonitor)

			form.OnSubmit = startAction
			form.OnCancel = nil
			form.Refresh()
		}
		form.OnSubmit = nil
		form.OnCancel = stopAction
		form.Refresh()

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

func startNewMonitoring(intervalSec int, threshold float64, consecutiveThreshold int, deviceSelector string) monitoring.Monitor {
	fmt.Printf("Interval: %d seconds, Threshold: %.2f, Consecutive Threshold: %d, Device Selector: %s\n",
		intervalSec, threshold, consecutiveThreshold, deviceSelector)

	newMonitor := monitoring.NewMonitor(intervalSec, threshold, consecutiveThreshold, deviceSelector)
	newMonitor.Start()

	return *newMonitor
}

func stopMonitoring(monitor monitoring.Monitor) {
	monitor.Stop()
}
