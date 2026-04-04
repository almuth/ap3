CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT,
    role TEXT NOT NULL DEFAULT 'user',
    status INTEGER NOT NULL DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL,
    nilai TEXT,
    deskripsi TEXT,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL
);
CREATE TABLE kategori_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);
CREATE TABLE satuan_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);
CREATE TABLE rak_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);
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
CREATE TABLE barang_harga_jual (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    barang_id INTEGER NOT NULL,
    harga INTEGER NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);
CREATE INDEX idx_harga_jual_barang ON barang_harga_jual(barang_id);
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
CREATE TABLE hutang_piutang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor TEXT NOT NULL UNIQUE,
    tanggal TEXT NOT NULL,
    profil_id INTEGER NOT NULL,
    tipe INTEGER NOT NULL CHECK(tipe IN (0, 1)),
    jumlah INTEGER NOT NULL,
    sisa INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 0 CHECK(status IN (0, 1, 2)),
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00',
    FOREIGN KEY (profil_id) REFERENCES profil(id)
);
CREATE INDEX idx_hutang_piutang_profil ON hutang_piutang(profil_id);
CREATE INDEX idx_hutang_piutang_status ON hutang_piutang(status);
CREATE TABLE inventory_balance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    barang_id INTEGER UNIQUE NOT NULL,
    qty INTEGER NOT NULL DEFAULT 0,
    qty_po INTEGER NOT NULL DEFAULT 0,
    qty_so INTEGER NOT NULL DEFAULT 0,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);
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
CREATE TABLE kategori_penerimaan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);
CREATE TABLE kategori_pengeluaran (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL COLLATE NOCASE UNIQUE,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT '2000-01-01 00:00:00'
);
