package ui

import (
	"fmt"

	"github.com/MawCeron/it-room/internal/db"
	"github.com/rivo/tview"
)

type App struct {
	app *tview.Application
	db  *db.DB
}

func NewApp(d *db.DB) *App {
	a := tview.NewApplication()
	app := &App{app: a, db: d}
	return app
}

func (a *App) Run() error {
	list := tview.NewList().ShowSecondaryText(false)
	list.AddItem("Load Assets...", "", 0, nil)
	list.AddItem("Quit", "", 'q', func() { a.app.Stop() })

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(tview.NewBox().SetBorder(true).SetTitle("IT Room"), 3, 0, false)
	flex.AddItem(list, 0, 1, true)

	a.app.SetRoot(flex, true)
	if err := a.app.Run(); err != nil {
		return fmt.Errorf("tview run: %w", err)
	}

	return nil
}
