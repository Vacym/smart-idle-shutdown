package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Vacym/smart-idle-shutdown/monitoring"
)

type gui struct {
	window           fyne.Window
	form             *widget.Form
	deviceSelector   *widget.Select
	intervalEntry    *widget.Entry
	thresholdEntry   *widget.Entry
	consecutiveEntry *widget.Entry
	monitor          *monitoring.Monitor
}

func newGUI(values monitoring.Settings) *gui {
	g := &gui{
		deviceSelector:   widget.NewSelect([]string{"CPU"}, func(s string) {}),
		intervalEntry:    widget.NewEntry(),
		thresholdEntry:   widget.NewEntry(),
		consecutiveEntry: widget.NewEntry(),
	}

	g.deviceSelector.SetSelected("CPU")
	g.intervalEntry.SetText(fmt.Sprintf("%d", values.Interval))
	g.thresholdEntry.SetText(fmt.Sprintf("%.2f", values.Threshold))
	g.consecutiveEntry.SetText(fmt.Sprintf("%d", values.ConsecutiveThreshold))

	g.intervalEntry.Validator = ValidateIntervalEntry
	g.thresholdEntry.Validator = ValidateThresholdEntry
	g.consecutiveEntry.Validator = ValidateConsecutiveEntry

	return g
}

func (g *gui) setupWindow(w fyne.Window) {
	g.window = w

	g.form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Device", Widget: g.deviceSelector, HintText: "Which device should be checked"},
			{Text: "Interval", Widget: g.intervalEntry, HintText: "Check interval in seconds"},
			{Text: "Threshold", Widget: g.thresholdEntry, HintText: "Load threshold in percents"},
			{Text: "Consecutive Threshold", Widget: g.consecutiveEntry, HintText: "Number of consecutive times the load should be below the threshold"},
		},
		SubmitText: "Start Monitoring",
		CancelText: "Stop monitoring",
		OnSubmit:   g.startMonitoring,
	}

	mainContainer := container.NewVBox(
		g.form,
	)

	w.SetContent(mainContainer)
}

func (g *gui) startMonitoring() {
	interval, err := strconv.Atoi(g.intervalEntry.Text)
	if err != nil || interval <= 0 {
		fmt.Println("Interval is incorrect")
		return
	}

	threshold, err := strconv.ParseFloat(g.thresholdEntry.Text, 64)
	if err != nil || threshold < 0 || threshold > 100 {
		fmt.Println("Threshold is incorrect")
		return
	}

	consecutiveThreshold, err := strconv.Atoi(g.consecutiveEntry.Text)
	if err != nil || consecutiveThreshold <= 0 {
		fmt.Println("Consecutive threshold is incorrect")
		return
	}

	if g.monitor != nil {
		fmt.Println("Monitoring already started")
		return
	}

	g.monitor = monitoring.NewMonitor(monitoring.Settings{
		Interval:             interval,
		Threshold:            threshold,
		ConsecutiveThreshold: consecutiveThreshold,
		Device:               g.deviceSelector.Selected,
	})
	g.monitor.Start()

	g.form.OnSubmit = nil
	g.form.OnCancel = g.stopMonitoring
	g.form.Refresh()
}

func (g *gui) stopMonitoring() {
	if g.monitor == nil {
		fmt.Println("No active monitoring to stop")
		return
	}
	g.monitor.Stop()
	g.monitor = nil

	g.form.OnSubmit = g.startMonitoring
	g.form.OnCancel = nil
	g.form.Refresh()
}
