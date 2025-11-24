package ui

import (
	"github.com/MawCeron/it-room/internal/db"
	"github.com/rivo/tview"
)

type LicensesPage struct {
	view *tview.Flex
	db   *db.DB
}

func NewLicensesPage(db *db.DB) *LicensesPage {
	p := &LicensesPage{db: db}
	p.build()
	return p
}

func (p *LicensesPage) Name() string {
	return "Licenses"
}

func (p *LicensesPage) View() tview.Primitive {
	return p.view
}

func (p *LicensesPage) build() {
	header := tview.NewTextView().
		SetText("[::b]Licenses[::-]\nLicense management and assignments").
		SetDynamicColors(true)

	content := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 2, 0, false).
		AddItem(nil, 1, 0, false)

	p.view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(nil, 2, 0, false).
				AddItem(content, 0, 1, true).
				AddItem(nil, 2, 0, false),
			0, 1, true).
		AddItem(nil, 1, 0, false)
}
