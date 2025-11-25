package assets

import (
	"fmt"

	"github.com/MawCeron/it-room/internal/models"
	"github.com/MawCeron/it-room/internal/repo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (p *AssetsPage) buildAssetsTable() *tview.Flex {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	p.addTableHeaders(table)

	assets := p.loadAssets()
	if assets == nil {
		return tview.NewFlex().
			AddItem(tview.NewTextView().
				SetText("Error loading assets"), 0, 1, false)
	}

	p.fillTableRows(table, assets)
	p.bindTableEvents(table, assets)

	box := tview.NewFlex().AddItem(table, 0, 1, true)
	box.SetBorder(true).
		SetTitle(" [::b]Assets[::-] - IT equipment inventory management ")

	return box
}

func (p *AssetsPage) addTableHeaders(t *tview.Table) {
	headers := []string{"Internal Code", "Type", "Model", "Serial Number", "Status"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)
		t.SetCell(0, col, cell)
	}
}

func (p *AssetsPage) loadAssets() []*models.Asset {
	assetRepo := repo.NewAssetRepo(p.db.Conn)
	assets, err := assetRepo.List()
	if err != nil {
		return nil
	}
	return assets
}

func (p *AssetsPage) fillTableRows(t *tview.Table, assets []*models.Asset) {
	for row, asset := range assets {
		r := row + 1
		t.SetCell(r, 0, tview.NewTableCell(asset.AssetCode))
		t.SetCell(r, 1, tview.NewTableCell(fmt.Sprintf("%d", asset.TypeID)))
		t.SetCell(r, 2, tview.NewTableCell(fmt.Sprintf("%s %s", asset.Make, asset.Model)))
		t.SetCell(r, 3, tview.NewTableCell(asset.SerialNumber))
		t.SetCell(r, 4, p.statusCell(asset.StatusID))
	}
}

func (p *AssetsPage) bindTableEvents(t *tview.Table, assets []*models.Asset) {
	t.SetSelectedFunc(func(row, _ int) {
		if row == 0 {
			return
		}
		asset := assets[row-1]
		p.showAssetModal(asset)
	})

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := t.GetSelection()

		switch event.Rune() {
		case 'n', 'N':
			p.showNewAssetForm()
			return nil

		case 'e', 'E':
			if row > 0 && row <= len(assets) {
				p.showEditAssetForm(assets[row-1])
			}
			return nil
		}

		return event
	})
}

func (p *AssetsPage) buildStatusBar() *tview.TextView {
	return tview.NewTextView().
		SetText(" [yellow]↑↓[white] Navigate  [yellow]Enter[white] View details [yellow]f[white] Filters  [yellow]n[white] New Asset  [yellow]a[white] Change Assignation  [red]r[white] Retire Asset  [yellow]?[white] Help").
		SetDynamicColors(true)
}
