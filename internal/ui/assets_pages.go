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
		SetLabel("🔍 ").
		SetPlaceholder("Buscar por código, marca, modelo, serie...").
		SetFieldWidth(40)

	btnFilters := tview.NewButton("Filtros").
		SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkSlateGray))

	btnNew := tview.NewButton("+ Nuevo Activo").
		SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkGreen))

	topBar := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(searchInput, 50, 0, true).
		AddItem(nil, 1, 0, false). // Espaciador
		AddItem(btnFilters, 10, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(btnNew, 16, 0, false).
		AddItem(nil, 0, 1, false) // Empuja a la izquierda

	// 2. Encabezado de sección
	header := tview.NewTextView().
		SetText("[::b]Activos[::-]\nGestión de equipos e inventario de TI").
		SetDynamicColors(true)

	// 3. Tabla de datos
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	headers := []string{"Código", "Tipo", "Marca/Modelo", "Serie", "Estatus", "Ubicación", "Acciones"}
	for col, h := range headers {
		cell := tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetExpansion(1)
		table.SetCell(0, col, cell)
	}

	// Datos de ejemplo
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

		// Estatus con color
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
	tableBox.SetBorder(true).SetTitle(" Inventario ")

	// 4. Barra de estado
	statusBar := tview.NewTextView().
		SetText(" [yellow]↑↓[white] Navegar  [yellow]Enter[white] Ver detalles  [yellow]n[white] Nuevo  [yellow]f[white] Filtros  [yellow]d[white] Eliminar  [yellow]?[white] Ayuda").
		SetDynamicColors(true)

	// Layout principal
	p.view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topBar, 1, 0, false).
		AddItem(nil, 1, 0, false). // Espaciador
		AddItem(header, 2, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tableBox, 0, 1, true).
		AddItem(statusBar, 1, 0, false)
}
