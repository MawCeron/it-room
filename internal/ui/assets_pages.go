package ui

import (
	"fmt"

	"github.com/MawCeron/it-room/internal/db"
	"github.com/MawCeron/it-room/internal/models"
	"github.com/MawCeron/it-room/internal/repo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AssetsPage struct {
	view  *tview.Flex
	db    *db.DB
	pages *tview.Pages
}

func NewAssetsPage(db *db.DB, pages *tview.Pages) *AssetsPage {
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

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	headers := []string{"Internal Code", "Type", "Model", "Serial Number", "Status"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)
		table.SetCell(0, col, cell)
	}

	assetRepo := repo.NewAssetRepo(p.db.Conn)

	assets, err := assetRepo.List()
	if err != nil {
		p.view = tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(
				tview.NewTextView().SetText(fmt.Sprintf("Error loading assets:\n%v[::-]", err)),
				0, 1, false)
		return
	}

	for row, asset := range assets {
		r := row + 1
		table.SetCell(r, 0, tview.NewTableCell(asset.AssetCode))
		table.SetCell(r, 1, tview.NewTableCell(fmt.Sprintf("%d", asset.TypeID)))
		table.SetCell(r, 2, tview.NewTableCell(fmt.Sprintf("%s %s", asset.Make, asset.Model)))
		table.SetCell(r, 3, tview.NewTableCell(asset.SerialNumber))

		var statusCell *tview.TableCell
		switch asset.StatusID {
		case 2:
			statusCell = tview.NewTableCell("Available").SetTextColor(tcell.ColorGreen)
		case 1:
			statusCell = tview.NewTableCell("Assigned").SetTextColor(tcell.ColorDodgerBlue)
		case 3:
			statusCell = tview.NewTableCell("Under Maintenance").SetTextColor(tcell.ColorOrange)
		case 4:
			statusCell = tview.NewTableCell("Retired").SetTextColor(tcell.ColorRed)
		default:
			statusCell = tview.NewTableCell("Unknown").SetTextColor(tcell.ColorGray)
		}
		table.SetCell(r, 4, statusCell)

	}

	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		asset := assets[row-1]
		p.showAssetModal(asset)
	})

	tableBox := tview.NewFlex().AddItem(table, 0, 1, true)
	tableBox.SetBorder(true).SetTitle(" [::b]Assets[::-] - IT equipment inventory management ")

	statusBar := tview.NewTextView().
		SetText(" [yellow]↑↓[white] Navigate  [yellow]Enter[white] View details  [yellow]n[white] New  [yellow]f[white] Filters  [yellow]d[white] Delete  [yellow]?[white] Help").
		SetDynamicColors(true)

	content := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tableBox, 0, 1, true).
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

func (p *AssetsPage) showAssetModal(asset *models.Asset) {
	modal := tview.NewModal().
		SetText("Asset Details:\n" +
			"Code: " + asset.AssetCode + "\n" +
			"Model: " + asset.Model).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(idx int, label string) {
			p.pages.RemovePage("assetModal")
		})

	p.pages.AddPage("assetModal", modal, true, true)
}
