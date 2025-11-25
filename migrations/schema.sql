-- ============================================
-- GearRoom: IT Asset Management System
-- Database Schema (SQLite)
-- ============================================

-- ======================================================
-- 1. Lookup Tables / Catalogs
-- ======================================================

-- Asset categories (high-level groups of equipment)
CREATE TABLE IF NOT EXISTS asset_categories (
    category_id INTEGER PRIMARY KEY AUTOINCREMENT,
    code_prefix TEXT NOT NULL UNIQUE,           -- Short prefix for asset tags
    description TEXT NOT NULL                   -- Human-readable category name
);

-- Asset types (specific product classes grouped by category)
CREATE TABLE IF NOT EXISTS asset_types (
    type_id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    type_name TEXT NOT NULL UNIQUE,
    FOREIGN KEY (category_id) REFERENCES asset_categories(category_id)
);

-- Asset statuses (workflow lifecycle)
CREATE TABLE IF NOT EXISTS asset_statuses (
    status_id INTEGER PRIMARY KEY AUTOINCREMENT,
    status_name TEXT NOT NULL UNIQUE
);

-- Maintenance types
CREATE TABLE IF NOT EXISTS maintenance_types (
    maintenance_type_id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_name TEXT NOT NULL UNIQUE
);

-- ======================================================
-- 2. Main Entities
-- ======================================================

-- Physical/organizational locations
CREATE TABLE IF NOT EXISTS locations (
    location_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL                          -- Example: Main, Remote
);

-- People who can receive assets or licenses
CREATE TABLE IF NOT EXISTS employees (
    employee_id INTEGER PRIMARY KEY AUTOINCREMENT,
    full_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

-- Assets (core entity)
CREATE TABLE IF NOT EXISTS assets (
    asset_id TEXT PRIMARY KEY,                  -- UUID stored as TEXT
    asset_tag TEXT NOT NULL UNIQUE,             -- Human-friendly asset tag
    type_id INTEGER NOT NULL,
    status_id INTEGER NOT NULL,
    serial_number TEXT NOT NULL UNIQUE,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    purchase_date TEXT NOT NULL,                -- ISO 8601 (YYYY-MM-DD)
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
    assignment_date TEXT NOT NULL DEFAULT (datetime('now')),
    return_date TEXT,
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
    removal_date TEXT,
    notes TEXT,
    
    FOREIGN KEY (license_id) REFERENCES software_licenses(license_id),
    FOREIGN KEY (asset_id) REFERENCES assets(asset_id)
);

-- Prevent multiple active (non-removed) assignments for the same pair
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
-- 7. Indexes for Performance
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
-- 8. Inserts (Categories, Types, Statuses, Locations)
-- ======================================================

-- Categories
INSERT INTO asset_categories (code_prefix, description) VALUES
    ('EQ', 'Computer Equipment'),
    ('PR', 'Printers and Multifunction Devices'),
    ('SC', 'Scanners'),
    ('NT', 'Network Devices'),
    ('MD', 'Mobile Devices'),
    ('MS', 'Monitors and Screens'),
    ('SR', 'Servers'),
    ('SD', 'Storage Devices'),
    ('AC', 'Accessories'),
    ('AV', 'Audio and Video Equipment'),
    ('TI', 'Tools and Infrastructure'),
    ('PE', 'Power Equipment');

-- Types (ordered by category_id)
INSERT INTO asset_types (category_id, type_name) VALUES
-- EQ (1)
(1, 'Laptop'),
(1, 'Desktop'),
(1, 'All-in-One'),
(1, 'Mini PC'),
(1, 'Workstation'),
(1, 'Thin Client'),

-- PR (2)
(2, 'Laser Printer'),
(2, 'Inkjet Printer'),
(2, 'Multifunction Printer'),
(2, 'Label Printer'),
(2, 'Plotter'),

-- SC (3)
(3, 'Document Scanner'),
(3, 'Flatbed Scanner'),
(3, 'Barcode Scanner'),
(3, 'ID Scanner'),

-- NT (4)
(4, 'Router'),
(4, 'Switch'),
(4, 'Firewall'),
(4, 'Access Point'),
(4, 'Network Appliance'),
(4, 'Modem'),

-- MD (5)
(5, 'Smartphone'),
(5, 'Tablet'),
(5, 'Rugged Device'),
(5, 'Handheld Terminal'),

-- MS (6)
(6, 'Monitor'),
(6, 'TV Display'),
(6, 'Digital Signage Display'),

-- SR (7)
(7, 'Rack Server'),
(7, 'Tower Server'),
(7, 'Blade Server'),
(7, 'Microserver'),

-- SD (8)
(8, 'External HDD'),
(8, 'External SSD'),
(8, 'Flash Drive'),
(8, 'Memory Card'),
(8, 'NAS'),

-- AC (9)
(9, 'Keyboard'),
(9, 'Mouse'),
(9, 'Headset'),
(9, 'Webcam'),
(9, 'Docking Station'),
(9, 'Power Adapter'),

-- AV (10)
(10, 'Projector'),
(10, 'Conference Speaker'),
(10, 'Conference Camera'),
(10, 'Microphone'),
(10, 'Mixer'),
(10, 'Amplifier'),

-- TI (11)
(11, 'Network Tools'),
(11, 'Electrical Tools'),
(11, 'Server Rack'),
(11, 'Patch Panel'),
(11, 'Cable Tester'),
(11, 'Tool Kit'),
(11, 'Label Maker'),

-- PE (12)
(12, 'UPS'),
(12, 'NoBreak'),
(12, 'Voltage Regulator'),
(12, 'Surge Protector'),
(12, 'Power Strip'),
(12, 'PDU');

-- Statuses
INSERT INTO asset_statuses (status_name) VALUES
('Assigned'),
('Available'),
('Under Maintenance'),
('Retired');

-- Locations
INSERT INTO locations (name, type) VALUES
('Main', 'Local'),
('Remote', 'Remote');

-- Maintenance Types
INSERT INTO maintenance_types (type_name) VALUES
('Preventive'),
('Corrective'),
('Upgrade');
