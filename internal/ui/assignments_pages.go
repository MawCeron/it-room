package ui

import (
	"github.com/MawCeron/it-room/internal/db"
	"github.com/rivo/tview"
)

type AssignmentsPage struct {
	view *tview.Flex
	db   *db.DB
}

func NewAssigmentsPage(db *db.DB) *AssignmentsPage {
	p := &AssignmentsPage{db: db}
	p.build()
	return p
}

func (p *AssignmentsPage) Name() string {
	return "Assigments"
}

func (p *AssignmentsPage) View() tview.Primitive {
	return p.view
}

func (p *AssignmentsPage) build() {
	header := tview.NewTextView().
		SetText("[::b]Assignments[::-]\nAsset assignments to employees").
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
