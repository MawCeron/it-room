package ui

import (
	"github.com/MawCeron/it-room/internal/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AssetsPage struct {
	view *tview.Flex
	db   *db.DB
}

func NewAssetsPage(db *db.DB) *AssetsPage {
	p := &AssetsPage{db: db}
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
	searchInput := tview.NewInputField().
		SetLabel("Search:").
		SetFieldWidth(25)

	btnNew := tview.NewButton("New Asset").
		SetStyle(tcell.StyleDefault.Background(tcell.ColorGreen))

	topBar := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(searchInput, 50, 0, true).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(btnNew, 16, 0, false)

	header := tview.NewTextView().
		SetText("[::b]Assets[::-]\nIT equipment inventory").
		SetDynamicColors(true)

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	headers := []string{"Internal Code", "Type", "Model", "Serial Number", "Status", "Location"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)
		table.SetCell(0, col, cell)
	}

	assets := []struct {
		code, aType, model, serial, status, location string
	}{
		{"LAP-001", "Laptop", "HP EliteBook 840 G8", "5CD1234ABC", "Asignado", "Oficina Central"},
		{"MON-015", "Monitor", "Dell P2422H", "CN-0ABC123", "Disponible", "Almacén TI"},
		{"IMP-003", "Impresora", "Canon imageRUNNER 2625i", "KBN12345", "En Mantenimiento", "Piso 2"},
	}

	for row, asset := range assets {
		r := row + 1
		table.SetCell(r, 0, tview.NewTableCell(asset.code))
		table.SetCell(r, 1, tview.NewTableCell(asset.aType))
		table.SetCell(r, 2, tview.NewTableCell(asset.model))
		table.SetCell(r, 3, tview.NewTableCell(asset.serial))

		statusCell := tview.NewTableCell(asset.status)
		switch asset.status {
		case "Disponible":
			statusCell.SetTextColor(tcell.ColorGreen)
		case "Asignado":
			statusCell.SetTextColor(tcell.ColorDodgerBlue)
		case "En Mantenimiento":
			statusCell.SetTextColor(tcell.ColorOrange)
		}
		table.SetCell(r, 4, statusCell)

		table.SetCell(r, 5, tview.NewTableCell(asset.location))
		table.SetCell(r, 6, tview.NewTableCell("[Ver]").SetTextColor(tcell.ColorAqua))
	}

	tableBox := tview.NewFlex().AddItem(table, 0, 1, true)
	tableBox.SetBorder(true).SetTitle(" Inventory ")

	statusBar := tview.NewTextView().
		SetText(" [yellow]↑↓[white] Navigate  [yellow]Enter[white] View details  [yellow]n[white] New  [yellow]f[white] Filters  [yellow]d[white] Delete  [yellow]?[white] Help").
		SetDynamicColors(true)

	content := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 2, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(topBar, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tableBox, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

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
