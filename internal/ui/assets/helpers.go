package assets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// statusCell returns a colored table cell based on the asset status ID
// Status IDs: 1=Assigned, 2=Available, 3=Under Maintenance, 4=Retired
func (p *AssetsPage) statusCell(statusID int) *tview.TableCell {
	switch statusID {
	case 2:
		return tview.NewTableCell("Available").SetTextColor(tcell.ColorGreen)
	case 1:
		return tview.NewTableCell("Assigned").SetTextColor(tcell.ColorDodgerBlue)
	case 3:
		return tview.NewTableCell("Under Maintenance").SetTextColor(tcell.ColorOrange)
	case 4:
		return tview.NewTableCell("Retired").SetTextColor(tcell.ColorRed)
	}
	return tview.NewTableCell("Unknown").SetTextColor(tcell.ColorGray)
}
