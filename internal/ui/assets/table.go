package assets

import (
	"fmt"

	"github.com/MawCeron/it-room/internal/models"
	"github.com/MawCeron/it-room/internal/repo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// buildAssetsTable creates the main assets table with headers, data, and event bindings
// Returns a flex container with the bordered table inside
func (p *AssetsPage) buildAssetsTable() *tview.Flex {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	p.addTableHeaders(table)

	assets := p.loadAssets()
	if len(assets) > 0 {
		p.fillTableRows(table, assets)
	}

	// Always bind events, even if assets is nil or empty
	p.bindTableEvents(table, assets)

	box := tview.NewFlex().AddItem(table, 0, 1, true)
	box.SetBorder(true).
		SetTitle(" [::b]Assets[::-] - IT equipment inventory management ")

	return box
}

// bindTableEvents attaches event handlers for table interactions
// Handles row selection (Enter) and keyboard shortcuts (n=new, e=edit)
func (p *AssetsPage) bindTableEvents(t *tview.Table, assets []*models.Asset) {
	// Handle row selection (Enter key)
	t.SetSelectedFunc(func(row, _ int) {
		if row == 0 || assets == nil || row > len(assets) {
			return
		}
		asset := assets[row-1]
		p.showAssetModal(asset)
	})

	// Handle keyboard shortcuts
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := t.GetSelection()

		switch event.Rune() {
		case 'n', 'N':
			p.showNewAssetForm()
			return nil
		case 'e', 'E':
			if assets != nil && row > 0 && row <= len(assets) {
				p.showEditAssetForm(assets[row-1])
			}
			return nil
		}
		return event
	})
}

// addTableHeaders sets up the column headers for the assets table
// Headers are displayed in yellow and are not selectable
func (p *AssetsPage) addTableHeaders(t *tview.Table) {
	headers := []string{"Asset Tag", "Type", "Model", "Serial Number", "Status"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)
		t.SetCell(0, col, cell)
	}
}

// loadAssets retrieves all assets from the database
// Returns nil if there's an error loading assets
func (p *AssetsPage) loadAssets() []*models.Asset {
	assetRepo := repo.NewAssetRepo(p.db.Conn)
	assets, err := assetRepo.List()
	if err != nil {
		return nil
	}
	return assets
}

// fillTableRows populates the table with asset data
// Each row displays: asset tag, type ID, maker+model, serial number, and colored status
func (p *AssetsPage) fillTableRows(t *tview.Table, assets []*models.Asset) {
	for row, asset := range assets {
		r := row + 1
		t.SetCell(r, 0, tview.NewTableCell(asset.AssetTag))
		t.SetCell(r, 1, tview.NewTableCell(fmt.Sprintf("%d", asset.TypeID)))
		t.SetCell(r, 2, tview.NewTableCell(fmt.Sprintf("%s %s", asset.Maker, asset.Model)))
		t.SetCell(r, 3, tview.NewTableCell(asset.SerialNumber))
		t.SetCell(r, 4, p.statusCell(asset.StatusID))
	}
}

// buildStatusBar creates the bottom status bar showing available keyboard shortcuts
// Displays navigation keys and action shortcuts with color formatting
func (p *AssetsPage) buildStatusBar() *tview.TextView {
	return tview.NewTextView().
		SetText(" [yellow]↑↓[white] Navigate  [yellow]Enter[white] View details [yellow]f[white] Filters  [yellow]n[white] New Asset  [yellow]a[white] Change Assignation  [red]r[white] Retire Asset  [yellow]?[white] Help").
		SetDynamicColors(true)
}
