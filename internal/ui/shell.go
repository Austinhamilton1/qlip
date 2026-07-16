package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Shell struct {
	app           *App
	root          *fyne.Container
	pageContainer *fyne.Container
	navbar        *NavBar
}

func NewShell(app *App) *Shell {
	s := &Shell{
		app: app,
	}

	s.pageContainer = container.NewStack()

	s.navbar = NewNavBar(app, s)

	s.root = container.NewBorder(
		s.navbar.Object(),
		nil,
		nil,
		nil,
		container.NewPadded(s.pageContainer),
	)

	return s
}

func (s *Shell) Root() fyne.CanvasObject {
	return s.root
}

func (s *Shell) SetPage(obj fyne.CanvasObject) {
	s.pageContainer.Objects = []fyne.CanvasObject{
		obj,
	}

	s.pageContainer.Refresh()
}
