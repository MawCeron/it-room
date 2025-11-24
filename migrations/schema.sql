-- ============================================
-- GearRoom: IT Asset Management System
-- Database: SQLite
-- ============================================

-- ======================================================
-- 1. Lookup Tables / Catalogs
-- ======================================================

CREATE TABLE IF NOT EXISTS asset_categories (
    category_id INTEGER PRIMARY KEY AUTOINCREMENT,
    code_prefix TEXT NOT NULL UNIQUE,           -- E.g.: EQ, PR, NT, AV
    description TEXT NOT NULL                   -- E.g.: Computing, Printing, Network, Audio/Video
);

CREATE TABLE IF NOT EXISTS asset_types (
    type_id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    type_name TEXT NOT NULL UNIQUE,             -- E.g.: Laptop, Desktop, Printer, Router
    FOREIGN KEY (category_id) REFERENCES asset_categories(category_id)
);

CREATE TABLE IF NOT EXISTS asset_statuses (
    status_id INTEGER PRIMARY KEY AUTOINCREMENT,
    status_name TEXT NOT NULL UNIQUE            -- E.g.: Assigned, Available, Retired, In Repair
);

CREATE TABLE IF NOT EXISTS maintenance_types (
    maintenance_type_id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_name TEXT NOT NULL UNIQUE              -- E.g.: Preventive, Corrective, Upgrade
);

-- ======================================================
-- 2. Main Entities
-- ======================================================

CREATE TABLE IF NOT EXISTS locations (
    location_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL                          -- E.g.: Office, Warehouse, Remote
);

CREATE TABLE IF NOT EXISTS employees (
    employee_id INTEGER PRIMARY KEY AUTOINCREMENT,
    full_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS assets (
    asset_id TEXT PRIMARY KEY,                  -- UUID as TEXT
    asset_code TEXT NOT NULL UNIQUE,
    type_id INTEGER NOT NULL,
    status_id INTEGER NOT NULL,
    serial_number TEXT NOT NULL UNIQUE,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    processor TEXT,
    ram_gb INTEGER,
    storage_tb REAL,
    storage_type TEXT,                          -- HDD, SSD, Hybrid
    operating_system TEXT,
    purchase_date TEXT NOT NULL,                -- ISO 8601: YYYY-MM-DD
    warranty_end_date TEXT,
    location_id INTEGER NOT NULL,
    notes TEXT,
    
    FOREIGN KEY (type_id) REFERENCES asset_types(type_id),
    FOREIGN KEY (status_id) REFERENCES asset_statuses(status_id),
    FOREIGN KEY (location_id) REFERENCES locations(location_id)
);

-- ======================================================
-- 3. Asset Assignments (Movement Log)
-- ======================================================

CREATE TABLE IF NOT EXISTS asset_assignments (
    assignment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    asset_id TEXT NOT NULL,
    employee_id INTEGER NOT NULL,
    assignment_date TEXT NOT NULL DEFAULT (datetime('now')),  -- ISO 8601
    return_date TEXT,                           -- Filled when asset is returned
    notes TEXT,
    
    FOREIGN KEY (asset_id) REFERENCES assets(asset_id),
    FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
);

-- ======================================================
-- 4. Maintenance
-- ======================================================

CREATE TABLE IF NOT EXISTS maintenance_logs (
    log_id INTEGER PRIMARY KEY AUTOINCREMENT,
    asset_id TEXT NOT NULL,
    maintenance_type_id INTEGER NOT NULL,
    maintenance_date TEXT NOT NULL DEFAULT (datetime('now')),
    cost REAL,
    description TEXT NOT NULL,
    performed_by TEXT,
    
    FOREIGN KEY (asset_id) REFERENCES assets(asset_id),
    FOREIGN KEY (maintenance_type_id) REFERENCES maintenance_types(maintenance_type_id)
);

-- ======================================================
-- 5. Software Licenses
-- ======================================================

CREATE TABLE IF NOT EXISTS software_licenses (
    license_id INTEGER PRIMARY KEY AUTOINCREMENT,
    software_name TEXT NOT NULL,
    license_key TEXT NOT NULL UNIQUE,
    license_type TEXT NOT NULL,                 -- Retail, Volume, Subscription
    seats_purchased INTEGER NOT NULL,
    purchase_date TEXT NOT NULL,
    expiration_date TEXT,
    notes TEXT
);

CREATE TABLE IF NOT EXISTS license_assignments (
    assignment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    license_id INTEGER NOT NULL,
    asset_id TEXT NOT NULL,
    assignment_date TEXT NOT NULL DEFAULT (datetime('now')),
    removal_date TEXT,                          -- NULL while license is active
    notes TEXT,
    
    FOREIGN KEY (license_id) REFERENCES software_licenses(license_id),
    FOREIGN KEY (asset_id) REFERENCES assets(asset_id)
);

-- Prevents multiple active assignments of the same License-Asset pair
CREATE UNIQUE INDEX IF NOT EXISTS idx_license_assignments_active_pair
ON license_assignments (license_id, asset_id)
WHERE removal_date IS NULL;

-- ======================================================
-- 6. Consumables
-- ======================================================

CREATE TABLE IF NOT EXISTS consumable_types (
    consumable_type_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    part_number TEXT,
    manufacturer TEXT,
    last_purchase_date TEXT
);

CREATE TABLE IF NOT EXISTS consumable_usage (
    usage_id INTEGER PRIMARY KEY AUTOINCREMENT,
    consumable_type_id INTEGER NOT NULL,
    asset_id TEXT NOT NULL,
    installation_date TEXT NOT NULL DEFAULT (datetime('now')),
    notes TEXT,
    
    FOREIGN KEY (consumable_type_id) REFERENCES consumable_types(consumable_type_id),
    FOREIGN KEY (asset_id) REFERENCES assets(asset_id)
);

-- ======================================================
-- 7. Additional Indexes for Performance
-- ======================================================

CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type_id);
CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status_id);
CREATE INDEX IF NOT EXISTS idx_assets_location ON assets(location_id);
CREATE INDEX IF NOT EXISTS idx_asset_assignments_asset ON asset_assignments(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_assignments_employee ON asset_assignments(employee_id);
CREATE INDEX IF NOT EXISTS idx_maintenance_logs_asset ON maintenance_logs(asset_id);
CREATE INDEX IF NOT EXISTS idx_license_assignments_license ON license_assignments(license_id);
CREATE INDEX IF NOT EXISTS idx_license_assignments_asset ON license_assignments(asset_id);

-- ======================================================

INSERT INTO asset_categories
(code_prefix, description)
VALUES
	('EQ', 'Computer Equipment'),
	('PR', 'Printers and Multifuction Devices'),
	('SC', 'Scanners'),
	('NT', 'Newtork Devices'),
	('MD', 'Mobile Devices'),
	('MS', 'Monitors and Screens'),
	('SR', 'Servers'),
	('AC', 'Accessories'),
	('UP', 'UPS'),
	('AV', 'Audio and Video Equipment'),
	('SD', 'Storage Devices');

INSERT INTO asset_statuses
(status_name)
VALUES('Assigned'),
('Available'),
('Under Maintenance'),
('Retired');