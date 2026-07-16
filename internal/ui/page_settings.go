package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewSettings(*App) fyne.CanvasObject {
	return widget.NewLabel("Settings")
}
