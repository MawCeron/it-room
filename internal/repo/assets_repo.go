package repo

import (
	"database/sql"
	"time"

	"github.com/MawCeron/it-room/internal/models"
)

type AssetRepo struct{ db *sql.DB }

func NewAssetRepo(db *sql.DB) *AssetRepo {
	return &AssetRepo{db: db}
}

func (r *AssetRepo) List() ([]*models.Asset, error) {
	rows, err := r.db.Query(`SELECT asset_id, asset_code, type_id, status_id, serial_number, make, model, processor, ram_gb, storage_tb, storage_type, operating_system, purchase_date, warranty_end_date, location_id, notes
FROM assets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Asset
	for rows.Next() {
		var a models.Asset
		var purchaseDate, warrantyEndDate sql.NullString

		if err := rows.Scan(&a.AssetID, &a.AssetCode, &a.TypeID, &a.StatusID,
			&a.SerialNumber, &a.Make, &a.Model, &a.Processor, &a.RamGB, &a.StorageTB,
			&a.StorageType, &a.OperatingSystem, &purchaseDate, &warrantyEndDate,
			&a.LocationID, &a.Notes); err != nil {
			return nil, err
		}

		// Parsear las fechas
		if purchaseDate.Valid {
			if t, err := time.Parse("2006-01-02", purchaseDate.String); err == nil {
				a.PurchaseDate = t
			}
		}

		if warrantyEndDate.Valid {
			if t, err := time.Parse("2006-01-02", warrantyEndDate.String); err == nil {
				a.WarrantyEndDate = &t
			}
		}

		out = append(out, &a)
	}

	return out, nil
}
