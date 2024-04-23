package main

import (
	"flag"

	"fyne.io/fyne/v2/app"
	"github.com/Vacym/smart-idle-shutdown/monitoring"
)

func main() {
	intervalPtr := flag.Int("interval", 5, "Check interval in seconds")
	thresholdPtr := flag.Float64("threshold", 30.0, "Load Threshold in Percent")
	consecutiveThresholdPtr := flag.Int("consecutive", 3, "Number of consecutive times the load should be below the threshold")

	flag.Parse()

	startSettings := monitoring.Settings{
		Interval:             *intervalPtr,
		Threshold:            *thresholdPtr,
		ConsecutiveThreshold: *consecutiveThresholdPtr,
	}

	a := app.New()
	w := a.NewWindow("Smart Idle Shutdown")

	// Initialize the GUI
	gui := newGUI(startSettings)

	// Pass the window to the GUI for content setup
	gui.setupWindow(w)

	// Launch the application
	w.ShowAndRun()
}
