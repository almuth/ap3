-- Seed Data for AhadPOS 3 SQLite Database
-- Generated: 2024-02-15

-- Insert default categories (from PHP migration)
INSERT INTO kategori_barang (nama, updated_at, updated_by, created_at) VALUES
('umum', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('wafer', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('biskuit', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('sirup', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('mie', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('kopi', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('isotonik drink', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('makanan', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('gula', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('minuman', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('susu', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00');

-- Insert default units (satuan)
INSERT INTO satuan_barang (nama, updated_at, updated_by, created_at) VALUES
('pcs', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('kg', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('Ons', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('Kardus', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('pak', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('lusin', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('box', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('set', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('gr', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('ltr', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00'),
('ml', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00');

-- Insert default rack
INSERT INTO rak_barang (nama, updated_at, updated_by, created_at) VALUES
('Rak 1', '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00');

-- Insert admin user (password: admin123 - hashed with bcrypt)
INSERT INTO user (name, username, password, email, role, status, created_at, updated_at) VALUES
('Administrator', 'admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'admin@ahadpos.local', 'admin', 1, datetime('now'), datetime('now'));

-- Insert default customer (UMUM)
INSERT INTO profil (tipe_id, nama, alamat, kota, telepon, updated_at, updated_by, created_at) VALUES
(2, 'UMUM', NULL, NULL, NULL, '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00');

-- Insert default cash/bank account
INSERT INTO kas_bank (nama, nomor_rekening, tipe, saldo, status, updated_at, updated_by, created_at) VALUES
('Kas Utama', NULL, 1, 0, 1, '2000-01-01 00:00:00', 1, '2000-01-01 00:00:00');

-- Insert some config values
INSERT INTO config (nama, nilai, deskripsi, updated_at, updated_by) VALUES
('app.name', 'AhadPOS 3', 'Nama Aplikasi', datetime('now'), 1),
('app.version', '3.0.0', 'Versi Aplikasi', datetime('now'), 1),
('pos.cara_caribarang', 'barcode', 'Cara mencari barang: barcode atau nama', datetime('now'), 1);
