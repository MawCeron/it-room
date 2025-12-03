package assets

import (
	"time"

	"github.com/MawCeron/it-room/internal/models"
	"github.com/MawCeron/it-room/internal/repo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (p *AssetsPage) showNewAssetForm() {
	p.showAssetForm(nil)
}

func (p *AssetsPage) showEditAssetForm(asset *models.Asset) {
	p.showAssetForm(asset)
}

const DateLayout = "2006-01-02"

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
	purchaseDate := time.Now().Format(DateLayout)
	warrantyEndDate := time.Now().AddDate(1, 0, 0).Format(DateLayout)
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
		SetFieldWidth(40)

	if assetCode == "" {
		assetCode = categoryPrefixes[0] + "-"
	}

	locationsRepo := repo.NewLocationRepo(p.db.Conn)
	locations, _ := locationsRepo.List()
	locationOptions := make([]string, len(locations))
	locationIDs := make([]int, len(locations))
	locationTypes := make([]string, len(locations))
	for i, l := range locations {
		locationOptions[i] = l.Name
		locationIDs[i] = l.LocationID
		locationTypes[i] = l.Type
	}

	locationDropDown := tview.NewDropDown().
		SetLabel("Location").
		SetOptions(locationOptions, nil).
		SetFieldWidth(40)

	assetTagInput := tview.NewInputField().
		SetLabel("Asset Tag").
		SetText(assetCode).
		SetFieldWidth(40)

	purchaseDateInput := tview.NewInputField().
		SetLabel("Purchase Date (YYYY-MM-DD)").
		SetText(purchaseDate).
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			if (lastChar >= '0' && lastChar <= '9') || lastChar == '-' || lastChar == 0 {
				if len(textToCheck) == 10 {
					_, err := time.Parse(DateLayout, textToCheck)
					return err == nil
				}
				return true
			}
			return false
		}).
		SetFieldWidth(40)

	warrantyEndInput := tview.NewInputField().
		SetLabel("Warranty End Date (YYYY-MM-DD)").
		SetText(warrantyEndDate).
		SetFieldWidth(40)

	purchaseDateInput.SetDoneFunc(func(key tcell.Key) {
		text := purchaseDateInput.GetText()

		// Parsear la fecha ingresada
		d, err := time.Parse("2006-01-02", text)
		if err != nil {
			return
		}

		// Crear la fecha + 1 año
		oneYearLater := d.AddDate(1, 0, 0)
		warrantyEndInput.SetText(oneYearLater.Format(DateLayout))
	})

	categoryDropDown := tview.NewDropDown().
		SetLabel("Category").
		SetOptions(categoryOptions, func(option string, optionIndex int) {
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
		}).
		SetCurrentOption(selectedOption).
		SetFieldWidth(40)

	form.AddFormItem(categoryDropDown)
	form.AddFormItem(typeDropDown)
	form.AddFormItem(assetTagInput)
	form.AddInputField("Make", maker, 40, nil, nil)
	form.AddInputField("Model", model, 40, nil, nil)
	form.AddInputField("Serial Number", serialNumber, 40, nil, nil)
	form.AddFormItem(purchaseDateInput)
	form.AddFormItem(warrantyEndInput)
	form.AddFormItem(locationDropDown)
	form.AddTextArea("Notes", "", 40, 0, 0, nil)

	form.AddButton("Save", nil)
	form.AddButton("Cancel", func() {
		p.pages.RemovePage("assetForm")
	})

	form.SetBorder(true).SetTitle(" " + title + " ")

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 30, 1, true).
			AddItem(nil, 0, 1, false), 80, 1, true).
		AddItem(nil, 0, 1, false)

	p.pages.AddPage("assetForm", flex, true, true)
}
