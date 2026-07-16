package ui

import (
	"time"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
)

type App struct {
	app    fyne.App
	window fyne.Window
	state  *State
	CurrentPage Page
	shell  *Shell
}

func New() *App {
	a := &App{}
	a.app = fyneapp.NewWithID("com.github.austinhamilton.qlip")
	a.window = a.app.NewWindow("Qlip")
	a.state = &State{}
	a.shell = NewShell(a)
	a.window.SetContent(a.shell.Root())
	return a
}

func (a *App) Run() {
	a.window.Show()

	go func() {
		time.Sleep(50 * time.Millisecond)
		fyne.Do(func() {
			a.window.Resize(fyne.NewSize(1000, 700))
		})
	}()

	a.app.Run()
}
