package assets

import (
	"github.com/MawCeron/it-room/internal/db"
	"github.com/rivo/tview"
)

type AssetsPage struct {
	view  *tview.Flex
	db    *db.DB
	pages *tview.Pages
}

func New(db *db.DB, pages *tview.Pages) *AssetsPage {
	p := &AssetsPage{db: db, pages: pages}
	p.build()
	return p
}

func (p *AssetsPage) Name() string {
	return "Assets"
}

func (p *AssetsPage) View() tview.Primitive {
	return p.view
}

func (p *AssetsPage) build() {
	table := p.buildAssetsTable()
	statusBar := p.buildStatusBar()

	content := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	p.view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(nil, 2, 0, false).
				AddItem(content, 0, 1, true).
				AddItem(nil, 2, 0, false),
			0, 1, true).
		AddItem(nil, 1, 0, false)
}
