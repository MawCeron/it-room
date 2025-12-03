package models

import (
	"time"
)

// Asset represents an IT asset in the inventory
type Asset struct {
	AssetID         string     `db:"asset_id"`
	AssetTag        string     `db:"asset_tag"`
	TypeID          int        `db:"type_id"`
	StatusID        int        `db:"status_id"`
	SerialNumber    string     `db:"serial_number"`
	Maker           string     `db:"make"`
	Model           string     `db:"model"`
	PurchaseDate    time.Time  `db:"purchase_date"`
	WarrantyEndDate *time.Time `db:"warranty_end_date"` // Nullable
	LocationID      int        `db:"location_id"`
	Notes           *string    `db:"notes"` // Nullable
}

type AssetCategory struct {
	CategoryId  int    `db:"category_id"`
	CodePrefix  string `db:"code_prefix"`
	Description string `db:"description"`
}

type AssetType struct {
	TypeID     int    `db:"type_id"`
	CategoryID int    `db:"category_id"`
	TypeName   string `db:"type_name"`
}

type Location struct {
	LocationID int    `db:"location_id"`
	Name       string `db:"name"`
	Type       string `db:"type"`
}
