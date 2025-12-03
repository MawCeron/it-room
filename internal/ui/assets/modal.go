package assets

import (
	"github.com/MawCeron/it-room/internal/models"
	"github.com/rivo/tview"
)

func (p *AssetsPage) showAssetModal(asset *models.Asset) {
	modal := tview.NewModal().
		SetText("Asset Details:\n" +
			"Code: " + asset.AssetTag + "\n" +
			"Model: " + asset.Model).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(idx int, label string) {
			p.pages.RemovePage("assetModal")
		})

	p.pages.AddPage("assetModal", modal, true, true)
}
