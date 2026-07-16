package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewHome(*App) fyne.CanvasObject {
	return widget.NewLabel("Home")
}
