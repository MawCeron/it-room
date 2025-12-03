package assets

import (
	"github.com/MawCeron/it-room/internal/db"
	"github.com/rivo/tview"
)

// AssetsPage represents the main assets management page
// It contains the view layout, database connection, and page manager
type AssetsPage struct {
	view  *tview.Flex
	db    *db.DB
	pages *tview.Pages
}

// New creates and initializes a new AssetsPage instance
// It builds the page layout and returns the configured page
func New(db *db.DB, pages *tview.Pages) *AssetsPage {
	p := &AssetsPage{db: db, pages: pages}
	p.build()
	return p
}

// Name returns the display name of this page
func (p *AssetsPage) Name() string {
	return "Assets"
}

// View returns the root primitive for this page
// This is used by the page manager to display the page
func (p *AssetsPage) View() tview.Primitive {
	return p.view
}

// build constructs the complete page layout
// It creates the assets table, status bar, and applies padding
func (p *AssetsPage) build() {
	table := p.buildAssetsTable()
	statusBar := p.buildStatusBar()

	// Main content area with table and status bar
	content := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	// Add horizontal and vertical padding
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
