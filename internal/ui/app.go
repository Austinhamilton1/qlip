package ui

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
)

type App struct {
	app    fyne.App
	window fyne.Window
	state  *State
	shell  *Shell
}

func New() *App {
	a := &App{}
	a.app = fyneapp.New()
	a.window = a.app.NewWindow("Qlip")
	a.window.Resize(fyne.NewSize(1000, 700))
	a.state = &State{}
	a.shell = NewShell(a)
	a.window.SetContent(a.shell.Root())
	return a
}

func (a *App) Run() {
	a.window.ShowAndRun()
}
