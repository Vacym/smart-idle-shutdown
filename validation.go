package main

import (
	"fyne.io/fyne/v2/data/validation"
)

var ValidateIntervalEntry = validation.NewRegexp(`^[0-9]+$`, "not a valid int number")

var ValidateThresholdEntry = validation.NewRegexp(`^[0-9]+(\.[0-9]+)?$`, "not a valid float number")

var ValidateConsecutiveEntry = validation.NewRegexp(`^[0-9]+$`, "not a valid int number")
