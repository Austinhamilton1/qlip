package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type navTab struct {
	Page Page
	root *fyne.Container
	background *canvas.Rectangle
	underline *canvas.Rectangle
	label *widget.Label
	icon *widget.Icon
}

type NavBar struct {
	app   *App
	shell *Shell
	bar   fyne.CanvasObject
	home *navTab
	connection *navTab
	settings *navTab
}

func NewNavBar(app *App, shell *Shell) *NavBar {
	n := &NavBar{
		app:   app,
		shell: shell,
	}

	n.home = n.makeTab(
		PageHome,
		"Home",
		theme.HomeIcon(),
	)

	n.connection = n.makeTab(
		PageConnection,
		"Connection",
		theme.ComputerIcon(),
	)

	n.settings = n.makeTab(
		PageSettings,
		"Settings",
		theme.SettingsIcon(),
	)

	n.bar = container.NewVBox(
		container.NewHBox(
			n.home.root,
			n.connection.root,
			n.settings.root,
		),
	)

	return n
}

func (n *NavBar) Navigate(page Page) {
	n.app.CurrentPage = page

	switch page {
	case PageHome:
		n.shell.SetPage(NewHome(n.app))
	case PageConnection:
		n.shell.SetPage(NewConnection(n.app))
	case PageSettings:
		n.shell.SetPage(NewSettings(n.app))
	}

	n.updateSelection()
}

func (n *NavBar) updateSelection() {
	tabs := []*navTab{
		n.home,
		n.connection,
		n.settings,
	}

	for _, tab := range tabs {
		tab.background.FillColor = color.Transparent
		tab.underline.FillColor = color.Transparent

		tab.background.Refresh()
		tab.underline.Refresh()
	}

	var selected *navTab

	switch n.app.CurrentPage {
	case PageHome:
		selected = n.home
	case PageConnection:
		selected = n.connection
	case PageSettings:
		selected = n.settings
	}

	if selected == nil {
		return
	}

	selected.background.FillColor = theme.SelectionColor()
	selected.underline.FillColor = color.NRGBA{
		R: 0,
		G: 170,
		B: 255,
		A: 255,
	}

	selected.background.Refresh()
	selected.underline.Refresh()
}

func (n *NavBar) makeTab(
	page Page,
	title string,
	icon fyne.Resource,
) *navTab {
	bg := canvas.NewRectangle(color.Transparent)

	line := canvas.NewRectangle(color.Transparent)
	line.SetMinSize(fyne.NewSize(0, 3))

	ic := widget.NewIcon(icon)

	lbl := widget.NewLabel(title)

	content := container.NewVBox (
		container.NewPadded(
			container.NewHBox(
				ic,
				lbl,
			),
		),
		line,
	)

	button := widget.NewButton("", func() {
		n.Navigate(page)
	})

	button.Importance = widget.LowImportance

	root := container.NewStack(
		bg,
		button,
		content,
	)

	return &navTab{
		Page: page,
		root: root,
		background: bg,
		underline: line,
		label: lbl,
		icon: ic,
	}
}

func (n *NavBar) Object() fyne.CanvasObject {
	return n.bar
}
