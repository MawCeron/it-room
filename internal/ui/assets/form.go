package assets

import (
	"time"

	"github.com/MawCeron/it-room/internal/models"
	"github.com/MawCeron/it-room/internal/repo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const DateLayout = "2006-01-02"

// showNewAssetForm displays the form to create a new asset
func (p *AssetsPage) showNewAssetForm() {
	p.showAssetForm(nil)
}

// showEditAssetForm displays the form to edit an existing asset
func (p *AssetsPage) showEditAssetForm(asset *models.Asset) {
	p.showAssetForm(asset)
}

// showAssetForm displays the asset form (create or edit)
func (p *AssetsPage) showAssetForm(asset *models.Asset) {
	isEdit := asset != nil
	title := "New Asset"
	if isEdit {
		title = "Edit Asset"
	}

	// Create the form with all components
	form := p.buildAssetForm(asset)
	form.SetBorder(true).SetTitle(" " + title + " ")

	// Create centered layout
	flex := p.createCenteredLayout(form)

	p.pages.AddPage("assetForm", flex, true, true)
}

// buildAssetForm builds the complete form with all its fields
func (p *AssetsPage) buildAssetForm(asset *models.Asset) *tview.Form {
	form := tview.NewForm()
	assetsRepo := repo.NewAssetRepo(p.db.Conn)

	// Load categories
	categories, _ := assetsRepo.GetAssetCategories()
	categoryData := p.prepareCategoryData(categories)

	// Load initial types
	types, _ := assetsRepo.GetAssetTypes(categoryData.IDs[0])
	typeData := p.prepareTypeData(types)

	// Load locations
	locationsRepo := repo.NewLocationRepo(p.db.Conn)
	locations, _ := locationsRepo.List()
	locationData := p.prepareLocationData(locations)

	// Initialize default values
	defaultValues := p.getDefaultFormValues(asset, categoryData.Prefixes[0])

	// Create form fields
	assetTagInput := p.createAssetTagInput(defaultValues.AssetCode)
	purchaseDateInput := p.createPurchaseDateInput(defaultValues.PurchaseDate)
	warrantyEndInput := p.createWarrantyEndInput(defaultValues.WarrantyEndDate)
	typeDropDown := p.createTypeDropDown(typeData.Options)
	locationDropDown := p.createLocationDropDown(locationData.Options)

	// Link purchase date to warranty
	p.linkPurchaseDateToWarranty(purchaseDateInput, warrantyEndInput)

	// Create category dropdown with type update logic
	categoryDropDown := p.createCategoryDropDown(
		categoryData,
		assetsRepo,
		typeDropDown,
		&typeData,
		assetTagInput,
	)

	// Add all fields to the form
	p.addFormFields(form, formFields{
		categoryDropDown:  categoryDropDown,
		typeDropDown:      typeDropDown,
		assetTagInput:     assetTagInput,
		purchaseDateInput: purchaseDateInput,
		warrantyEndInput:  warrantyEndInput,
		locationDropDown:  locationDropDown,
		defaultValues:     defaultValues,
	})

	// Add buttons
	p.addFormButtons(form)

	return form
}

// formFields groups all form fields
type formFields struct {
	categoryDropDown  *tview.DropDown
	typeDropDown      *tview.DropDown
	assetTagInput     *tview.InputField
	purchaseDateInput *tview.InputField
	warrantyEndInput  *tview.InputField
	locationDropDown  *tview.DropDown
	defaultValues     formDefaultValues
}

// formDefaultValues contains the form default values
type formDefaultValues struct {
	AssetCode       string
	SerialNumber    string
	Maker           string
	Model           string
	PurchaseDate    string
	WarrantyEndDate string
}

// categoryData contains category information
type categoryData struct {
	Options  []string
	IDs      []int
	Prefixes []string
}

// typeData contains type information
type typeData struct {
	Options []string
	IDs     []int
}

// locationData contains location information
type locationData struct {
	Options []string
	IDs     []int
	Types   []string
}

// prepareCategoryData extracts and organizes category data
func (p *AssetsPage) prepareCategoryData(categories []*models.AssetCategory) categoryData {
	data := categoryData{
		Options:  make([]string, len(categories)),
		IDs:      make([]int, len(categories)),
		Prefixes: make([]string, len(categories)),
	}

	for i, c := range categories {
		data.Options[i] = c.Description
		data.IDs[i] = c.CategoryId
		data.Prefixes[i] = c.CodePrefix
	}

	return data
}

// prepareTypeData extracts and organizes type data
func (p *AssetsPage) prepareTypeData(types []*models.AssetType) typeData {
	data := typeData{
		Options: make([]string, len(types)),
		IDs:     make([]int, len(types)),
	}

	for i, t := range types {
		data.Options[i] = t.TypeName
		data.IDs[i] = t.TypeID
	}

	return data
}

// prepareLocationData extracts and organizes location data
func (p *AssetsPage) prepareLocationData(locations []*models.Location) locationData {
	data := locationData{
		Options: make([]string, len(locations)),
		IDs:     make([]int, len(locations)),
		Types:   make([]string, len(locations)),
	}

	for i, l := range locations {
		data.Options[i] = l.Name
		data.IDs[i] = l.LocationID
		data.Types[i] = l.Type
	}

	return data
}

// getDefaultFormValues gets the default values for the form
func (p *AssetsPage) getDefaultFormValues(asset *models.Asset, defaultPrefix string) formDefaultValues {
	defaults := formDefaultValues{
		AssetCode:       defaultPrefix + "-",
		PurchaseDate:    time.Now().Format(DateLayout),
		WarrantyEndDate: time.Now().AddDate(1, 0, 0).Format(DateLayout),
	}

	if asset != nil {
		defaults.AssetCode = asset.AssetTag
		defaults.SerialNumber = asset.SerialNumber
		defaults.Maker = asset.Maker
		defaults.Model = asset.Model
		// Add more fields as needed
	}

	return defaults
}

// createAssetTagInput creates the asset tag input field
func (p *AssetsPage) createAssetTagInput(defaultValue string) *tview.InputField {
	return tview.NewInputField().
		SetLabel("Asset Tag").
		SetText(defaultValue).
		SetFieldWidth(40)
}

// createPurchaseDateInput creates the purchase date input field
func (p *AssetsPage) createPurchaseDateInput(defaultValue string) *tview.InputField {
	return tview.NewInputField().
		SetLabel("Purchase Date (YYYY-MM-DD)").
		SetText(defaultValue).
		SetAcceptanceFunc(p.dateAcceptanceFunc).
		SetFieldWidth(40)
}

// createWarrantyEndInput creates the warranty end date input field
func (p *AssetsPage) createWarrantyEndInput(defaultValue string) *tview.InputField {
	return tview.NewInputField().
		SetLabel("Warranty End Date (YYYY-MM-DD)").
		SetText(defaultValue).
		SetFieldWidth(40)
}

// createTypeDropDown creates the type dropdown
func (p *AssetsPage) createTypeDropDown(options []string) *tview.DropDown {
	return tview.NewDropDown().
		SetLabel("Type").
		SetOptions(options, nil).
		SetFieldWidth(40)
}

// createLocationDropDown creates the location dropdown
func (p *AssetsPage) createLocationDropDown(options []string) *tview.DropDown {
	return tview.NewDropDown().
		SetLabel("Location").
		SetOptions(options, nil).
		SetFieldWidth(40)
}

// dateAcceptanceFunc validates the date format while typing
func (p *AssetsPage) dateAcceptanceFunc(textToCheck string, lastChar rune) bool {
	if (lastChar >= '0' && lastChar <= '9') || lastChar == '-' || lastChar == 0 {
		if len(textToCheck) == 10 {
			_, err := time.Parse(DateLayout, textToCheck)
			return err == nil
		}
		return true
	}
	return false
}

// linkPurchaseDateToWarranty links purchase date to warranty date
func (p *AssetsPage) linkPurchaseDateToWarranty(purchaseInput, warrantyInput *tview.InputField) {
	purchaseInput.SetDoneFunc(func(key tcell.Key) {
		text := purchaseInput.GetText()
		d, err := time.Parse(DateLayout, text)
		if err != nil {
			return
		}
		oneYearLater := d.AddDate(1, 0, 0)
		warrantyInput.SetText(oneYearLater.Format(DateLayout))
	})
}

// createCategoryDropDown creates the category dropdown with update logic
func (p *AssetsPage) createCategoryDropDown(
	catData categoryData,
	assetsRepo *repo.AssetRepo,
	typeDropDown *tview.DropDown,
	typeData *typeData,
	assetTagInput *tview.InputField,
) *tview.DropDown {
	return tview.NewDropDown().
		SetLabel("Category").
		SetOptions(catData.Options, func(option string, optionIndex int) {
			// Update types based on selected category
			types, _ := assetsRepo.GetAssetTypes(catData.IDs[optionIndex])
			*typeData = p.prepareTypeData(types)

			typeDropDown.SetOptions(typeData.Options, nil)
			typeDropDown.SetCurrentOption(0)

			// Update asset tag prefix
			assetTagInput.SetText(catData.Prefixes[optionIndex] + "-")
		}).
		SetCurrentOption(0).
		SetFieldWidth(40)
}

// addFormFields adds all fields to the form
func (p *AssetsPage) addFormFields(form *tview.Form, fields formFields) {
	form.AddFormItem(fields.categoryDropDown)
	form.AddFormItem(fields.typeDropDown)
	form.AddFormItem(fields.assetTagInput)
	form.AddInputField("Make", fields.defaultValues.Maker, 40, nil, nil)
	form.AddInputField("Model", fields.defaultValues.Model, 40, nil, nil)
	form.AddInputField("Serial Number", fields.defaultValues.SerialNumber, 40, nil, nil)
	form.AddFormItem(fields.purchaseDateInput)
	form.AddFormItem(fields.warrantyEndInput)
	form.AddFormItem(fields.locationDropDown)
	form.AddTextArea("Notes", "", 40, 0, 0, nil)
}

// addFormButtons adds buttons to the form
func (p *AssetsPage) addFormButtons(form *tview.Form) {
	form.AddButton("Save", nil) // TODO: Implement save logic
	form.AddButton("Cancel", func() {
		p.pages.RemovePage("assetForm")
	})
}

// createCenteredLayout creates a centered layout for the form
func (p *AssetsPage) createCenteredLayout(content tview.Primitive) *tview.Flex {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(content, 30, 1, true).
			AddItem(nil, 0, 1, false), 80, 1, true).
		AddItem(nil, 0, 1, false)
}
