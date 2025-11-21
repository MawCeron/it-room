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

	pages := tview.NewPages()
	assetsPage := NewAssetsPage(a.db)
	pages.AddPage(assetsPage.Name(), assetsPage.View(), true, true)

	menu := tview.NewList()
	menuWidth := 20
	menu.AddItem("Assets", "", 0, func() {
		pages.SwitchToPage(assetsPage.Name())
	})
	menu.AddItem("Assignments", "", 0, nil)
	menu.AddItem("Maintenance", "", 0, nil)
	menu.AddItem("Licenses", "", 0, nil)
	menu.AddItem("Consumables", "", 0, nil)
	menu.ShowSecondaryText(false)

	frame := tview.NewFrame(menu)
	frame.SetBorder(true)
	frame.SetBorders(1, 0, 1, 1, 1, 1)
	frame.SetTitle(" IT Room ")

	flex := tview.NewFlex()
	flex.AddItem(frame, menuWidth+3, 1, true)
	flex.AddItem(pages, 0, 1, false)

	a.app.SetRoot(flex, true).EnableMouse(true)
	if err := a.app.Run(); err != nil {
		return fmt.Errorf("tview run: %w", err)
	}

	return nil
}
