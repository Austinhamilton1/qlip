package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type NavBar struct {
	app   *App
	shell *Shell
	bar   fyne.CanvasObject
}

func NewNavBar(app *App, shell *Shell) *NavBar {
	n := &NavBar{
		app:   app,
		shell: shell,
	}

	home := widget.NewButtonWithIcon(
		"Home",
		theme.HomeIcon(),
		func() {
			shell.SetPage(
				NewHome(app),
			)
		},
	)

	connection := widget.NewButtonWithIcon(
		"Connections",
		theme.ComputerIcon(),
		func() {
			shell.SetPage(
				NewConnection(app),
			)
		},
	)

	settings := widget.NewButtonWithIcon(
		"Settings",
		theme.SettingsIcon(),
		func() {
			shell.SetPage(
				NewSettings(app),
			)
		},
	)

	n.bar = container.NewVBox(
		container.NewHBox(
			home,
			connection,
			settings,
		),
		widget.NewSeparator(),
	)

	shell.SetPage(
		NewHome(app),
	)

	return n
}

func (n *NavBar) Object() fyne.CanvasObject {
	return n.bar
}
