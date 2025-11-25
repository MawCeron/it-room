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

func (r *AssetRepo) GetAssetCategories() ([]*models.AssetCategory, error) {
	rows, err := r.db.Query(`SELECT category_id, code_prefix, description
FROM asset_categories;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.AssetCategory
	for rows.Next() {
		var c models.AssetCategory

		if err := rows.Scan(&c.CategoryId, &c.CodePrefix, &c.Description); err != nil {
			return nil, err
		}

		out = append(out, &c)
	}

	return out, nil
}

func (r *AssetRepo) GetAssetTypes(category int) ([]*models.AssetType, error) {
	rows, err := r.db.Query(`SELECT type_id, category_id, type_name
FROM asset_types WHERE category_id = ?`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.AssetType
	for rows.Next() {
		var t models.AssetType

		if err := rows.Scan(&t.TypeID, &t.CategoryID, &t.TypeName); err != nil {
			return nil, err
		}

		out = append(out, &t)
	}

	return out, nil
}
