package ui

import (
	"fmt"
	"time"

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

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()

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

	tableBox := tview.NewFlex().AddItem(table, 0, 1, true)
	tableBox.SetBorder(true).SetTitle(" [::b]Assets[::-] - IT equipment inventory management ")

	statusBar := tview.NewTextView().
		SetText(" [yellow]↑↓[white] Navigate  [yellow]Enter[white] View details [yellow]f[white] Filters  [yellow]n[white] New Asset  [yellow]a[white] Change Assignation  [red]r[white] Retire Asset  [yellow]?[white] Help").
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

func (p *AssetsPage) showNewAssetForm() {
	p.showAssetForm(nil)
}

func (p *AssetsPage) showEditAssetForm(asset *models.Asset) {
	p.showAssetForm(asset)
}

func (p *AssetsPage) showAssetForm(asset *models.Asset) {
	isEdit := asset != nil
	title := "New Asset"
	if isEdit {
		title = "Edit Asset"
	}

	form := tview.NewForm()
	assetsRepo := repo.NewAssetRepo(p.db.Conn)
	categories, _ := assetsRepo.GetAssetCategories()

	categoryOptions := make([]string, len(categories))
	categoryIDs := make([]int, len(categories))
	categoryPrefixes := make([]string, len(categories))
	for i, c := range categories {
		categoryOptions[i] = c.Description
		categoryIDs[i] = c.CategoryId
		categoryPrefixes[i] = c.CodePrefix
	}

	var assetCode, serialNumber, maker, model string
	purchaseDate := time.Now().Format("2006-01-02")
	var selectedOption int

	types, _ := assetsRepo.GetAssetTypes(categoryIDs[selectedOption])
	typeOptions := make([]string, len(types))
	typeIDs := make([]int, len(types))
	for i, t := range types {
		typeOptions[i] = t.TypeName
		typeIDs[i] = t.TypeID
	}

	typeDropDown := tview.NewDropDown().
		SetLabel("Type").
		SetOptions(typeOptions, nil).
		SetFieldWidth(30)

	if assetCode == "" {
		assetCode = categoryPrefixes[0] + "-"
	}

	assetTagInput := tview.NewInputField().
		SetLabel("Asset Tag").
		SetText(assetCode).
		SetFieldWidth(30)

	form.AddDropDown("Category", categoryOptions, selectedOption, func(option string, optionIndex int) {
		types, _ = assetsRepo.GetAssetTypes(categoryIDs[optionIndex])
		typeOptions := make([]string, len(types))
		typeIDs := make([]int, len(types))
		for i, t := range types {
			typeOptions[i] = t.TypeName
			typeIDs[i] = t.TypeID
		}

		typeDropDown.SetOptions(typeOptions, nil)
		typeDropDown.SetCurrentOption(0)
		assetTagInput.SetText(categoryPrefixes[optionIndex] + "-")
	})
	form.AddFormItem(typeDropDown)

	// 3. Agregar resto de campos
	form.AddFormItem(assetTagInput)
	form.AddInputField("Make", maker, 30, nil, nil)
	form.AddInputField("Model", model, 30, nil, nil)
	form.AddInputField("Serial Number", serialNumber, 30, nil, nil)
	form.AddInputField("Purchase Date (YYYY-MM-DD)", purchaseDate, 15, nil, nil)

	form.AddButton("Save", nil)
	form.AddButton("Cancel", func() {
		p.pages.RemovePage("assetForm")
	})

	form.SetBorder(true).SetTitle(" " + title + " ")
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 20, 1, true).
			AddItem(nil, 0, 1, false), 80, 1, true).
		AddItem(nil, 0, 1, false)

	p.pages.AddPage("assetForm", flex, true, true)

}
