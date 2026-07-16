package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewConnection(*App) fyne.CanvasObject {
	return widget.NewLabel("Connection")
}
