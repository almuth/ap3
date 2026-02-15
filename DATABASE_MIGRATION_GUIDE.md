# Database Migration Guide: MySQL to SQLite

## Overview

This guide provides detailed instructions for migrating the AhadPOS 3 database from MySQL to SQLite, including schema conversion, data transformation, and validation steps.

## Prerequisites

### Tools Required
```bash
# MySQL tools
sudo apt-get install mysql-client mysqldump

# SQLite tools
sudo apt-get install sqlite3

# Data transformation tools (optional)
pip install pandas mysql-connector-python
```

## Phase 1: Analysis & Planning

### Step 1: Export Current MySQL Schema
```bash
# Export structure only
mysqldump -h localhost -u root -p --no-data ahadpos3 > mysql_schema.sql

# Export data only
mysqldump -h localhost -u root -p --no-create-info ahadpos3 > mysql_data.sql

# Export complete with procedures
mysqldump -h localhost -u root -p \
    --routines --triggers --events \
    ahadpos3 > mysql_complete.sql
```

### Step 2: Document Database Schema

Create a schema inventory:
```sql
-- List all tables
SELECT
    TABLE_NAME,
    TABLE_ROWS,
    ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024, 2) AS Size_MB
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'ahadpos3'
ORDER BY TABLE_NAME;
```

```sql
-- List all columns with types
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_KEY,
    COLUMN_DEFAULT,
    EXTRA
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = 'ahadpos3'
ORDER BY TABLE_NAME, ORDINAL_POSITION;
```

```sql
-- List all foreign keys
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = 'ahadpos3'
    AND REFERENCED_TABLE_NAME IS NOT NULL
ORDER BY TABLE_NAME;
```

### Step 3: Identify MySQL-Specific Features

Check for:
- ENUM types
- AUTO_INCREMENT on non-primary key columns
- TIMESTAMP defaults (ON UPDATE CURRENT_TIMESTAMP)
- ENGINE specifications (InnoDB)
- CHARACTER SET and COLLATE settings
- FULLTEXT indexes
- Stored procedures/functions
- Triggers
- Views

## Phase 2: SQLite Schema Conversion

### Type Mapping Reference

| MySQL Type | SQLite Type | Notes |
|------------|-------------|-------|
| TINYINT(1) | INTEGER | For boolean values |
| TINYINT | INTEGER | 0-255 |
| SMALLINT | INTEGER | -32768 to 32767 |
| INT, INTEGER | INTEGER | Primary key |
| BIGINT | INTEGER | Large integers |
| DECIMAL(M,D) | REAL | For currency, consider INTEGER in cents |
| FLOAT | REAL | Floating point |
| DOUBLE | REAL | Double precision |
| CHAR(N) | TEXT | Fixed-length strings |
| VARCHAR(N) | TEXT | Variable-length strings |
| TEXT | TEXT | Long text |
| DATE | TEXT | Format: 'YYYY-MM-DD' |
| DATETIME | TEXT | Format: 'YYYY-MM-DD HH:MM:SS' |
| TIMESTAMP | TEXT | Same as DATETIME |
| BLOB | BLOB | Binary data |
| ENUM | TEXT + CHECK | Use TEXT with CHECK constraint |

### Schema Conversion Examples

#### Example 1: Basic Table Conversion

**MySQL:**
```sql
CREATE TABLE `barang` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `barcode` varchar(30) CHARACTER SET utf8 NOT NULL,
    `nama` varchar(45) CHARACTER SET utf8 NOT NULL,
    `kategori_id` int(10) unsigned NOT NULL,
    `satuan_id` int(10) unsigned NOT NULL,
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0=tidak aktif; 1=aktif',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_by` int(10) unsigned NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT '2000-01-01 00:00:00',
    PRIMARY KEY (`id`),
    UNIQUE KEY `barcode_UNIQUE` (`barcode`),
    KEY `fk_barang_kategori_idx` (`kategori_id`),
    KEY `fk_barang_satuan_idx` (`satuan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
```

**SQLite:**
```sql
CREATE TABLE barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    barcode TEXT NOT NULL COLLATE NOCASE,
    nama TEXT NOT NULL,
    kategori_id INTEGER NOT NULL,
    satuan_id INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 1 CHECK(status IN (0, 1)),
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    UNIQUE(barang),
    FOREIGN KEY (kategori_id) REFERENCES kategori_barang(id),
    FOREIGN KEY (satuan_id) REFERENCES satuan_barang(id)
);

CREATE INDEX idx_barang_kategori ON barang(kategori_id);
CREATE INDEX idx_barang_satuan ON barang(satuan_id);

-- Trigger for auto-updating updated_at
CREATE TRIGGER update_barang_updated_at
AFTER UPDATE ON barang
BEGIN
    UPDATE barang SET updated_at = datetime('now') WHERE id = NEW.id;
END;
```

#### Example 2: Table with Decimal Fields

**MySQL:**
```sql
CREATE TABLE `penjualan_detail` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `penjualan_id` int(10) unsigned NOT NULL,
    `barang_id` int(10) unsigned NOT NULL,
    `qty` int(10) unsigned NOT NULL,
    `harga_satuan` decimal(18,2) NOT NULL,
    `diskon` decimal(18,2) DEFAULT '0.00',
    `total` decimal(18,2) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

**SQLite (Option 1: Use REAL - simpler but potential precision issues):**
```sql
CREATE TABLE penjualan_detail (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    penjualan_id INTEGER NOT NULL,
    barang_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    harga_satuan REAL NOT NULL,
    diskon REAL DEFAULT 0.0,
    total REAL NOT NULL,
    FOREIGN KEY (penjualan_id) REFERENCES penjualan(id),
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);
```

**SQLite (Option 2: Use INTEGER in cents - recommended for financial accuracy):**
```sql
CREATE TABLE penjualan_detail (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    penjualan_id INTEGER NOT NULL,
    barang_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    harga_satuan INTEGER NOT NULL,  -- Store in cents (e.g., 15000 for Rp 150.00)
    diskon INTEGER DEFAULT 0,        -- Store in cents
    total INTEGER NOT NULL,          -- Store in cents
    FOREIGN KEY (penjualan_id) REFERENCES penjualan(id),
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);
```

#### Example 3: ENUM Conversion

**MySQL:**
```sql
CREATE TABLE `user` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `role` enum('admin','kasir','supervisor') NOT NULL DEFAULT 'user',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;
```

**SQLite:**
```sql
CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('admin', 'kasir', 'supervisor'))
);
```

#### Example 4: Text Indexes

**MySQL:**
```sql
CREATE TABLE `kategori_barang` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `nama` varchar(45) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `nama` (`nama`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

**SQLite (Case-insensitive unique):**
```sql
CREATE TABLE kategori_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE,
    UNIQUE(nama)
);
```

## Phase 3: Complete SQLite Schema

### Core Tables Schema

```sql
-- ============================================
-- BEGIN TRANSACTION
-- ============================================
PRAGMA foreign_keys = OFF;

-- Users table
CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT,
    role TEXT NOT NULL DEFAULT 'user',
    status INTEGER NOT NULL DEFAULT 1,
    menu_id INTEGER
);

-- Config table
CREATE TABLE config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL,
    nilai TEXT,
    deskripsi TEXT,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL
);

-- Kategori Barang (Product Categories)
CREATE TABLE kategori_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);

CREATE TRIGGER update_kategori_barang_updated_at
AFTER UPDATE ON kategori_barang
BEGIN
    UPDATE kategori_barang SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Satuan Barang (Product Units)
CREATE TABLE satuan_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);

CREATE TRIGGER update_satuan_barang_updated_at
AFTER UPDATE ON satuan_barang
BEGIN
    UPDATE satuan_barang SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Rak Barang (Product Racks)
CREATE TABLE rak_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);

CREATE TRIGGER update_rak_barang_updated_at
AFTER UPDATE ON rak_barang
BEGIN
    UPDATE rak_barang SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Struktur Barang (Product Structure/Classification)
CREATE TABLE struktur_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    nama TEXT NOT NULL,
    level INTEGER NOT NULL,
    lft INTEGER NOT NULL,
    rgt INTEGER NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES struktur_barang(id)
);

CREATE TRIGGER update_struktur_barang_updated_at
AFTER UPDATE ON struktur_barang
BEGIN
    UPDATE struktur_barang SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Barang (Products)
CREATE TABLE barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    barcode TEXT NOT NULL COLLATE NOCASE UNIQUE,
    nama TEXT NOT NULL,
    struktur_id INTEGER,
    kategori_id INTEGER,
    satuan_id INTEGER NOT NULL,
    rak_id INTEGER,
    restock_point INTEGER NOT NULL DEFAULT 0,
    restock_level INTEGER NOT NULL DEFAULT 0,
    restock_min INTEGER NOT NULL DEFAULT 0,
    variant_coefficient REAL DEFAULT 1.0,
    status INTEGER NOT NULL DEFAULT 1 CHECK(status IN (0, 1)),
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (struktur_id) REFERENCES struktur_barang(id),
    FOREIGN KEY (kategori_id) REFERENCES kategori_barang(id),
    FOREIGN KEY (satuan_id) REFERENCES satuan_barang(id),
    FOREIGN KEY (rak_id) REFERENCES rak_barang(id)
);

CREATE INDEX idx_barang_kategori ON barang(kategori_id);
CREATE INDEX idx_barang_satuan ON barang(satuan_id);
CREATE INDEX idx_barang_rak ON barang(rak_id);
CREATE INDEX idx_barang_status ON barang(status);

CREATE TRIGGER update_barang_updated_at
AFTER UPDATE ON barang
BEGIN
    UPDATE barang SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Harga Jual (Selling Prices)
CREATE TABLE barang_harga_jual (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    barang_id INTEGER NOT NULL,
    harga INTEGER NOT NULL,  -- Stored in cents
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);

CREATE INDEX idx_harga_jual_barang ON barang_harga_jual(barang_id);

CREATE TRIGGER update_barang_harga_jual_updated_at
AFTER UPDATE ON barang_harga_jual
BEGIN
    UPDATE barang_harga_jual SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Profil (Customers and Suppliers)
CREATE TABLE profil (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tipe_id INTEGER NOT NULL,
    nama TEXT NOT NULL,
    alamat TEXT,
    kota TEXT,
    telepon TEXT,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);

CREATE INDEX idx_profil_tipe ON profil(tipe_id);

CREATE TRIGGER update_profil_updated_at
AFTER UPDATE ON profil
BEGIN
    UPDATE profil SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Penjualan (Sales)
CREATE TABLE penjualan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor TEXT NOT NULL UNIQUE,
    tanggal TEXT NOT NULL,
    profil_id INTEGER NOT NULL,
    hutang_piutang_id INTEGER,
    transfer_mode INTEGER NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 0 CHECK(status IN (0, 1, 2)),
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (profil_id) REFERENCES profil(id),
    FOREIGN KEY (hutang_piutang_id) REFERENCES hutang_piutang(id)
);

CREATE INDEX idx_penjualan_tanggal ON penjualan(tanggal);
CREATE INDEX idx_penjualan_profil ON penjualan(profil_id);
CREATE INDEX idx_penjualan_status ON penjualan(status);

CREATE TRIGGER update_penjualan_updated_at
AFTER UPDATE ON penjualan
BEGIN
    UPDATE penjualan SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Penjualan Detail (Sale Items)
CREATE TABLE penjualan_detail (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    penjualan_id INTEGER NOT NULL,
    barang_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    harga_satuan INTEGER NOT NULL,
    diskon INTEGER DEFAULT 0,
    diskon_persen REAL DEFAULT 0,
    total INTEGER NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (penjualan_id) REFERENCES penjualan(id) ON DELETE CASCADE,
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);

CREATE INDEX idx_penjualan_detail_penjualan ON penjualan_detail(penjualan_id);
CREATE INDEX idx_penjualan_detail_barang ON penjualan_detail(barang_id);

-- Pembelian (Purchases)
CREATE TABLE pembelian (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor TEXT NOT NULL UNIQUE,
    tanggal TEXT NOT NULL,
    profil_id INTEGER NOT NULL,
    referensi TEXT,
    tanggal_referensi TEXT,
    hutang_piutang_id INTEGER,
    status INTEGER NOT NULL DEFAULT 0 CHECK(status IN (0, 1, 2)),
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (profil_id) REFERENCES profil(id),
    FOREIGN KEY (hutang_piutang_id) REFERENCES hutang_piutang(id)
);

CREATE INDEX idx_pembelian_tanggal ON pembelian(tanggal);
CREATE INDEX idx_pembelian_profil ON pembelian(profil_id);

CREATE TRIGGER update_pembelian_updated_at
AFTER UPDATE ON pembelian
BEGIN
    UPDATE pembelian SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Pembelian Detail (Purchase Items)
CREATE TABLE pembelian_detail (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pembelian_id INTEGER NOT NULL,
    barang_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    harga_satuan INTEGER NOT NULL,
    diskon INTEGER DEFAULT 0,
    total INTEGER NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (pembelian_id) REFERENCES pembelian(id) ON DELETE CASCADE,
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);

CREATE INDEX idx_pembelian_detail_pembelian ON pembelian_detail(pembelian_id);
CREATE INDEX idx_pembelian_detail_barang ON pembelian_detail(barang_id);

-- Hutang Piutang (Debts/Credits)
CREATE TABLE hutang_piutang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor TEXT NOT NULL UNIQUE,
    tanggal TEXT NOT NULL,
    profil_id INTEGER NOT NULL,
    tipe INTEGER NOT NULL CHECK(tipe IN (0, 1)), -- 0=hutang, 1=piutang
    jumlah INTEGER NOT NULL,
    sisa INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 0 CHECK(status IN (0, 1, 2)), -- 0=aktif, 1=sebagian, 2=lunas
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (profil_id) REFERENCES profil(id)
);

CREATE INDEX idx_hutang_piutang_profil ON hutang_piutang(profil_id);
CREATE INDEX idx_hutang_piutang_status ON hutang_piutang(status);

CREATE TRIGGER update_hutang_piutang_updated_at
AFTER UPDATE ON hutang_piutang
BEGIN
    UPDATE hutang_piutang SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Inventory Balance
CREATE TABLE inventory_balance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    barang_id INTEGER UNIQUE NOT NULL,
    qty INTEGER NOT NULL DEFAULT 0,
    qty_po INTEGER NOT NULL DEFAULT 0,
    qty_so INTEGER NOT NULL DEFAULT 0,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);

CREATE TRIGGER update_inventory_balance_updated_at
AFTER UPDATE ON inventory_balance
BEGIN
    UPDATE inventory_balance SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Kas Bank (Cash/Bank Accounts)
CREATE TABLE kas_bank (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL,
    nomor_rekening TEXT,
    tipe INTEGER NOT NULL,
    saldo INTEGER NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);

-- Penerimaan (Income/Receipts)
CREATE TABLE penerimaan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor TEXT NOT NULL UNIQUE,
    tanggal TEXT NOT NULL,
    kas_bank_id INTEGER,
    kategori_id INTEGER,
    profil_id INTEGER,
    jumlah INTEGER NOT NULL,
    keterangan TEXT,
    status INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (kas_bank_id) REFERENCES kas_bank(id),
    FOREIGN KEY (kategori_id) REFERENCES kategori_penerimaan(id),
    FOREIGN KEY (profil_id) REFERENCES profil(id)
);

CREATE INDEX idx_penerimaan_tanggal ON penerimaan(tanggal);

CREATE TRIGGER update_penerimaan_updated_at
AFTER UPDATE ON penerimaan
BEGIN
    UPDATE penerimaan SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Pengeluaran (Expenses)
CREATE TABLE pengeluaran (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor TEXT NOT NULL UNIQUE,
    tanggal TEXT NOT NULL,
    kas_bank_id INTEGER,
    kategori_id INTEGER,
    profil_id INTEGER,
    jumlah INTEGER NOT NULL,
    keterangan TEXT,
    status INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (kas_bank_id) REFERENCES kas_bank(id),
    FOREIGN KEY (kategori_id) REFERENCES kategori_pengeluaran(id),
    FOREIGN KEY (profil_id) REFERENCES profil(id)
);

CREATE INDEX idx_pengeluaran_tanggal ON pengeluaran(tanggal);

CREATE TRIGGER update_pengeluaran_updated_at
AFTER UPDATE ON pengeluaran
BEGIN
    UPDATE pengeluaran SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Additional supporting tables would go here...
-- (kategori_penerimaan, kategori_pengeluaran, etc.)

PRAGMA foreign_keys = ON;
-- ============================================
-- END TRANSACTION
-- ============================================
```

## Phase 4: Data Migration

### Method 1: CSV Export/Import (Recommended)

#### Step 1: Export MySQL Data to CSV

Create a script to export all tables:
```bash
#!/bin/bash
DB_HOST="localhost"
DB_USER="root"
DB_PASS=""
DB_NAME="ahadpos3"
OUTPUT_DIR="mysql_exports"

mkdir -p $OUTPUT_DIR

# Get list of tables
TABLES=$(mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "SHOW TABLES FROM $DB_NAME" | awk '{print $1}' | grep -v "Tables_in")

# Export each table to CSV
for TABLE in $TABLES; do
    echo "Exporting $TABLE..."
    mysql -h $DB_HOST -u $DB_USER -p$DB_PASS $DB_NAME -e "SELECT * FROM $TABLE" \
        | sed 's/\t/","/g;s/^/"/;s/$/"/;s/\n//g' \
        > $OUTPUT_DIR/${TABLE}.csv
done

echo "Export complete!"
```

#### Step 2: Transform Data Types

Create a Python script for transformation:
```python
#!/usr/bin/env python3
import pandas as pd
import sqlite3
import os

# Configuration
mysql_dir = "mysql_exports"
sqlite_db = "ahadpos3.db"

# Type conversions needed
TYPE_CONVERSIONS = {
    'barang': {
        'restock_point': 'int',
        'restock_level': 'int',
        'restock_min': 'int',
        'status': 'int',
    },
    'penjualan_detail': {
        'harga_satuan': 'int',  # Convert to cents
        'diskon': 'int',        # Convert to cents
        'total': 'int',         # Convert to cents
    },
    # Add more conversions as needed
}

# Decimal to cents conversion for financial fields
def convert_decimal_to_cents(value):
    """Convert decimal string to integer cents"""
    if pd.isna(value) or value == '':
        return 0
    return int(float(value) * 100)

# Import data into SQLite
conn = sqlite3.connect(sqlite_db)
cursor = conn.cursor()

# Process each CSV file
for csv_file in os.listdir(mysql_dir):
    if not csv_file.endswith('.csv'):
        continue

    table_name = csv_file.replace('.csv', '')
    csv_path = os.path.join(mysql_dir, csv_file)

    print(f"Processing {table_name}...")

    # Read CSV
    df = pd.read_csv(csv_path)

    # Apply type conversions
    if table_name in TYPE_CONVERSIONS:
        conversions = TYPE_CONVERSIONS[table_name]
        for col, dtype in conversions.items():
            if col in df.columns:
                if dtype == 'int':
                    df[col] = pd.to_numeric(df[col], errors='coerce').fillna(0).astype('int64')

    # Handle financial fields (decimal to cents)
    financial_tables = ['penjualan_detail', 'pembelian_detail', 'barang_harga_jual']
    if table_name in financial_tables:
        for col in ['harga_satuan', 'diskon', 'total', 'harga']:
            if col in df.columns:
                df[col] = df[col].apply(convert_decimal_to_cents)

    # Insert into SQLite
    df.to_sql(table_name, conn, if_exists='append', index=False)

    print(f"  Imported {len(df)} rows into {table_name}")

conn.commit()
conn.close()
print("Data migration complete!")
```

### Method 2: Direct SQL Transformation

Create a transformation script:
```bash
#!/bin/bash
MYSQL_HOST="localhost"
MYSQL_USER="root"
MYSQL_PASS=""
MYSQL_DB="ahadpos3"
SQLITE_DB="ahadpos3.db"

# Create SQLite database
sqlite3 $SQLITE_DB "VACUUM;"

# Process each table
mysql -h $MYSQL_HOST -u $MYSQL_USER -p$MYSQL_PASS $MYSQL_DB -N -e "SHOW TABLES" | while read table; do
    echo "Processing table: $table"

    # Get CREATE TABLE statement and transform
    mysql -h $MYSQL_HOST -u $MYSQL_USER -p$MYSQL_PASS $MYSQL_DB -e "SHOW CREATE TABLE $table\G" | \
        grep "Create Table" | \
        sed 's/.*CREATE TABLE //' | \
        sed 's/ ENGINE=.*$//' | \
        sed 's/ AUTO_INCREMENT//' | \
        sed 's/ unsigned//g' | \
        sed 's/ int(/ INTEGER(/g' | \
        sed 's/ tinyint(1)/ INTEGER/g' | \
        sed 's/ tinyint(/ INTEGER(/g' | \
        sed 's/ varchar(/ TEXT/g' | \
        sed 's/ text/ TEXT/g' | \
        sed 's/ date/ TEXT/g' | \
        sed 's/ datetime/ TEXT/g' | \
        sed 's/ timestamp/ TEXT/g' | \
        sed 's/ decimal(/ REAL(/g' | \
        sed 's/ CHARACTER SET.*COLLATE.*//' | \
        sed 's/,/\n,/g' > /tmp/${table}.sql

    # Export data with transformed dates
    mysql -h $MYSQL_HOST -u $MYSQL_USER -p$MYSQL_PASS $MYSQL_DB -e "SELECT * FROM $table" | \
        sed 's/\t/|/g' | \
        while IFS='|' read -a row; do
            # Transform each row as needed
            echo "${row[@]}"
        done > /tmp/${table}_data.txt

    # Import into SQLite
    sqlite3 $SQLITE_DB <<EOF
$(cat /tmp/${table}.sql)
.import --csv /tmp/${table}_data.txt $table
EOF

    echo "Completed table: $table"
done
```

## Phase 5: Validation

### Record Count Validation

```sql
-- MySQL: Count records in each table
SELECT
    'barang' as table_name, COUNT(*) as count FROM barang
UNION ALL
SELECT 'penjualan', COUNT(*) FROM penjualan
UNION ALL
SELECT 'penjualan_detail', COUNT(*) FROM penjualan_detail
UNION ALL
SELECT 'pembelian', COUNT(*) FROM pembelian
UNION ALL
SELECT 'user', COUNT(*) FROM user
UNION ALL
SELECT 'config', COUNT(*) FROM config;
```

```sql
-- SQLite: Count records in each table
SELECT
    'barang' as table_name, COUNT(*) as count FROM barang
UNION ALL
SELECT 'penjualan', COUNT(*) FROM penjualan
UNION ALL
SELECT 'penjualan_detail', COUNT(*) FROM penjualan_detail
UNION ALL
SELECT 'pembelian', COUNT(*) FROM pembelian
UNION ALL
SELECT 'user', COUNT(*) FROM user
UNION ALL
SELECT 'config', COUNT(*) FROM config;
```

### Data Integrity Checks

```sql
-- Check for NULL in required fields
SELECT 'barang' as table_name, COUNT(*) as null_barcode
FROM barang WHERE barcode IS NULL OR barcode = ''
UNION ALL
SELECT 'barang', COUNT(*) FROM barang WHERE nama IS NULL OR nama = ''
UNION ALL
SELECT 'penjualan', COUNT(*) FROM penjualan WHERE nomor IS NULL OR nomor = '';

-- Check foreign key integrity
SELECT COUNT(*) as orphaned_penjualan_detail
FROM penjualan_detail pd
LEFT JOIN penjualan p ON pd.penjualan_id = p.id
WHERE p.id IS NULL;

SELECT COUNT(*) as orphaned_penjualan_detail_barang
FROM penjualan_detail pd
LEFT JOIN barang b ON pd.barang_id = b.id
WHERE b.id IS NULL;

-- Check data consistency
-- Total of penjualan_detail should match penjualan totals
SELECT p.id, p.nomor,
    (SELECT COALESCE(SUM(total), 0) FROM penjualan_detail WHERE penjualan_id = p.id) as detail_total
FROM penjualan p;
```

### Performance Checks

```sql
-- Analyze SQLite database
PRAGMA table_info(barang);
PRAGMA foreign_key_list(barang);
PRAGMA index_list(barang);

-- Check database size
SELECT page_count * page_size as bytes FROM pragma_page_count(), pragma_page_size();

-- Check query performance
EXPLAIN QUERY PLAN SELECT * FROM barang WHERE barcode = '123456789';
EXPLAIN QUERY PLAN
SELECT * FROM penjualan_detail
JOIN penjualan ON penjualan_detail.penjualan_id = penjualan.id
WHERE penjualan.tanggal >= '2024-01-01';
```

## Phase 6: Post-Migration Tasks

### 1. Create Initial Indexes

```sql
-- Performance indexes
CREATE INDEX idx_barang_nama ON barang(nama);
CREATE INDEX idx_penjualan_tanggal_profil ON penjualan(tanggal, profil_id);
CREATE INDEX idx_penjualan_detail_barang_penjualan ON penjualan_detail(barang_id, penjualan_id);
```

### 2. Recalculate Inventory Balance

```sql
-- Recalculate inventory balance from transactions
INSERT OR REPLACE INTO inventory_balance (barang_id, qty, updated_at)
SELECT
    barang_id,
    COALESCE(SUM(qty), 0) as qty,
    datetime('now') as updated_at
FROM pembelian_detail
GROUP BY barang_id;
```

### 3. Update Statistics

```sql
ANALYZE;
VACUUM;
```

### 4. Verify Business Logic

```sql
-- Check that all sales have details
SELECT p.id, p.nomor
FROM penjualan p
LEFT JOIN penjualan_detail pd ON p.id = pd.penjualan_id
WHERE pd.id IS NULL;

-- Check that no detail has NULL references
SELECT COUNT(*) FROM penjualan_detail WHERE penjualan_id IS NULL;
SELECT COUNT(*) FROM penjualan_detail WHERE barang_id IS NULL;
```

## Rollback Plan

If issues are detected:

```bash
# 1. Stop the application
systemctl stop ahadpos-api

# 2. Backup the SQLite database (for analysis)
cp ahadpos3.db ahadpos3.db.failed

# 3. Restore MySQL database
mysql -u root -p ahadpos3 < mysql_backup.sql

# 4. Update application config to point to MySQL
# Edit config/database.php

# 5. Restart PHP application
systemctl restart apache2

# 6. Analyze issues
# Review logs, data mismatches, etc.
```

## Migration Checklist

### Pre-Migration
- [ ] Backup MySQL database
- [ ] Document all custom queries
- [ ] Identify data inconsistencies
- [ ] Plan for data transformation
- [ ] Test migration on copy of database
- [ ] Schedule downtime window

### During Migration
- [ ] Export MySQL schema
- [ ] Convert schema to SQLite
- [ ] Create SQLite database
- [ ] Export MySQL data
- [ ] Transform data types
- [ ] Import to SQLite
- [ ] Apply indexes and triggers

### Post-Migration
- [ ] Validate record counts
- [ ] Check data integrity
- [ ] Verify foreign keys
- [ ] Test critical queries
- [ ] Performance test
- [ ] Recalculate derived data
- [ ] Run application tests
- [ ] Update application configuration
- [ ] Monitor for issues

### Cleanup
- [ ] Archive migration scripts
- [ ] Document any manual fixes
- [ ] Update documentation
- [ ] Train team on SQLite specifics

## Tips and Best Practices

1. **Test First**: Always test on a copy of production data
2. **Incremental Migration**: Migrate in phases if possible
3. **Validate Early**: Check data after each major step
4. **Keep Backups**: Maintain MySQL backup until confident
5. **Monitor Performance**: Watch for slow queries in SQLite
6. **Use Transactions**: Wrap bulk operations in transactions
7. **Index Wisely**: SQLite handles indexes differently than MySQL
8. **Date Format**: Consistently use ISO 8601 format (YYYY-MM-DD HH:MM:SS)
9. **Text Case**: Consider using COLLATE NOCASE for case-insensitive comparisons
10. **Decimal Precision**: Store financial values as integers in cents

## Common Issues and Solutions

### Issue 1: ENUM Conversion
**Problem**: SQLite doesn't support ENUM types
**Solution**: Use TEXT with CHECK constraint

### Issue 2: Auto-increment on Non-Primary Key
**Problem**: SQLite AUTOINCREMENT only works on PRIMARY KEY
**Solution**: Use separate sequence table or application logic

### Issue 3: ON UPDATE CURRENT_TIMESTAMP
**Problem**: SQLite doesn't support this
**Solution**: Use triggers (as shown in examples)

### Issue 4: Case-Insensitive Unique
**Problem**: UNIQUE is case-sensitive by default
**Solution**: Use COLLATE NOCASE on the column

### Issue 5: Decimal Precision
**Problem**: Floating point arithmetic can have precision issues
**Solution**: Store financial amounts as integers in cents

## Conclusion

This migration guide provides a comprehensive approach to converting the AhadPOS 3 database from MySQL to SQLite. Follow each phase carefully, validate at each step, and maintain backups throughout the process. The key to success is thorough testing and validation before deploying to production.
