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
	p.view = tview.NewFlex().
		AddItem(tview.NewTextView().SetText("Assignments Under Construction"), 0, 0, false)
}
