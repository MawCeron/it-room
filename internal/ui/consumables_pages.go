package ui

import (
	"github.com/MawCeron/it-room/internal/db"
	"github.com/rivo/tview"
)

type ConsumablesPage struct {
	view *tview.Flex
	db   *db.DB
}

func NewConsumablesPage(db *db.DB) *ConsumablesPage {
	p := &ConsumablesPage{db: db}
	p.build()
	return p
}

func (p *ConsumablesPage) Name() string {
	return "Consumables"
}

func (p *ConsumablesPage) View() tview.Primitive {
	return p.view
}

func (p *ConsumablesPage) build() {
	header := tview.NewTextView().
		SetText("[::b]Consumables[::-]\nToner, drum, and other consumables tracking").
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
