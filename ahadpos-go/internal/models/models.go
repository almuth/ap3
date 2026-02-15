package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// User represents a user account
type User struct {
	BaseModel
	ID       uint   `json:"id" gorm:"primarykey"`
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
	Email    string `json:"email"`
	Role     string `json:"role" gorm:"not null;default:'user'"`
	Status   int    `json:"status" gorm:"not null;default:1"`
}

// Config represents application configuration
type Config struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Nama      string    `json:"nama" gorm:"uniqueIndex;not null"`
	Nilai     string    `json:"nilai"`
	Deskripsi string    `json:"deskripsi"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy uint      `json:"updated_by" gorm:"not null"`
}

// KategoriBarang represents product categories
type KategoriBarang struct {
	BaseModel
	ID        uint   `json:"id" gorm:"primarykey"`
	Nama      string `json:"nama" gorm:"uniqueIndex;not null;collate:NOCASE"`
	UpdatedBy uint   `json:"updated_by" gorm:"not null"`
}

// SatuanBarang represents product units
type SatuanBarang struct {
	BaseModel
	ID        uint   `json:"id" gorm:"primarykey"`
	Nama      string `json:"nama" gorm:"uniqueIndex;not null;collate:NOCASE"`
	UpdatedBy uint   `json:"updated_by" gorm:"not null"`
}

// RakBarang represents product racks
type RakBarang struct {
	BaseModel
	ID        uint   `json:"id" gorm:"primarykey"`
	Nama      string `json:"nama" gorm:"uniqueIndex;not null;collate:NOCASE"`
	UpdatedBy uint   `json:"updated_by" gorm:"not null"`
}

// StrukturBarang represents product structure/hierarchy
type StrukturBarang struct {
	BaseModel
	ID        uint   `json:"id" gorm:"primarykey"`
	ParentID  *uint  `json:"parent_id" gorm:"index"`
	Nama      string `json:"nama" gorm:"not null"`
	Level     int    `json:"level" gorm:"not null"`
	Lft       int    `json:"lft" gorm:"not null"`
	Rgt       int    `json:"rgt" gorm:"not null"`
	UpdatedBy uint   `json:"updated_by" gorm:"not null"`
	Parent    *StrukturBarang `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children  []StrukturBarang `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// Barang represents products
type Barang struct {
	BaseModel
	ID                uint             `json:"id" gorm:"primarykey"`
	Barcode           string           `json:"barcode" gorm:"uniqueIndex;not null;collate:NOCASE"`
	Nama              string           `json:"nama" gorm:"not null"`
	StrukturID        *uint            `json:"struktur_id" gorm:"index"`
	KategoriID        *uint            `json:"kategori_id" gorm:"index"`
	SatuanID          uint             `json:"satuan_id" gorm:"not null;index"`
	RakID             *uint            `json:"rak_id" gorm:"index"`
	RestockPoint      int              `json:"restock_point" gorm:"default:0"`
	RestockLevel      int              `json:"restock_level" gorm:"default:0"`
	RestockMin        int              `json:"restock_min" gorm:"default:0"`
	VariantCoefficient float64          `json:"variant_coefficient" gorm:"default:1.0"`
	Status            int              `json:"status" gorm:"not null;default:1;check:status IN (0, 1)"`
	UpdatedBy         uint             `json:"updated_by" gorm:"not null"`
	Kategori          *KategoriBarang  `json:"kategori,omitempty" gorm:"foreignKey:KategoriID"`
	Satuan            *SatuanBarang    `json:"satuan,omitempty" gorm:"foreignKey:SatuanID"`
	Rak               *RakBarang       `json:"rak,omitempty" gorm:"foreignKey:RakID"`
	Struktur          *StrukturBarang  `json:"struktur,omitempty" gorm:"foreignKey:StrukturID"`
	HargaJual         []BarangHargaJual `json:"harga_jual,omitempty" gorm:"foreignKey:BarangID"`
	InventoryBalance  *InventoryBalance `json:"inventory_balance,omitempty" gorm:"foreignKey:BarangID"`
}

// BarangHargaJual represents product selling prices
type BarangHargaJual struct {
	BaseModel
	ID        uint    `json:"id" gorm:"primarykey"`
	BarangID  uint    `json:"barang_id" gorm:"not null;index"`
	Harga     int     `json:"harga" gorm:"not null"`
	UpdatedBy uint    `json:"updated_by" gorm:"not null"`
	Barang    *Barang `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
}

// Profil represents customers and suppliers
type Profil struct {
	BaseModel
	ID      uint   `json:"id" gorm:"primarykey"`
	TipeID  uint   `json:"tipe_id" gorm:"index;not null"`
	Nama    string `json:"nama" gorm:"not null"`
	Alamat  string `json:"alamat"`
	Kota    string `json:"kota"`
	Telepon string `json:"telepon"`
	UpdatedBy uint `json:"updated_by" gorm:"not null"`
}

// Penjualan represents sales (header)
type Penjualan struct {
	BaseModel
	ID               uint     `json:"id" gorm:"primarykey"`
	Nomor            string   `json:"nomor" gorm:"uniqueIndex;not null"`
	Tanggal          string   `json:"tanggal" gorm:"not null;index"`
	ProfilID         uint     `json:"profil_id" gorm:"index;not null"`
	HutangPiutangID  *uint    `json:"hutang_piutang_id" gorm:"index"`
	TransferMode     int      `json:"transfer_mode" gorm:"default:0"`
	Status           int      `json:"status" gorm:"not null;default:0;check:status IN (0, 1, 2);index"`
	UpdatedBy        uint     `json:"updated_by" gorm:"not null"`
	Profil           *Profil  `json:"profil,omitempty" gorm:"foreignKey:ProfilID"`
	HutangPiutang    *HutangPiutang `json:"hutang_piutang,omitempty" gorm:"foreignKey:HutangPiutangID"`
	Details          []PenjualanDetail `json:"details,omitempty" gorm:"foreignKey:PenjualanID"`
}

// PenjualanDetail represents sales line items
type PenjualanDetail struct {
	BaseModel
	ID          uint   `json:"id" gorm:"primarykey"`
	PenjualanID uint   `json:"penjualan_id" gorm:"not null;index"`
	BarangID    uint   `json:"barang_id" gorm:"not null;index"`
	Qty         int    `json:"qty" gorm:"not null"`
	HargaSatuan int    `json:"harga_satuan" gorm:"not null"`
	Diskon      int    `json:"diskon" gorm:"default:0"`
	DiskonPersen float64 `json:"diskon_persen" gorm:"default:0"`
	Total       int    `json:"total" gorm:"not null"`
	UpdatedBy   uint   `json:"updated_by" gorm:"not null"`
	Penjualan   *Penjualan `json:"penjualan,omitempty" gorm:"foreignKey:PenjualanID"`
	Barang      *Barang    `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
}

// Pembelian represents purchases (header)
type Pembelian struct {
	BaseModel
	ID              uint             `json:"id" gorm:"primarykey"`
	Nomor           string           `json:"nomor" gorm:"uniqueIndex;not null"`
	Tanggal         string           `json:"tanggal" gorm:"not null;index"`
	ProfilID        uint             `json:"profil_id" gorm:"index;not null"`
	Referensi       string           `json:"referensi"`
	TanggalReferensi string          `json:"tanggal_referensi"`
	HutangPiutangID *uint            `json:"hutang_piutang_id" gorm:"index"`
	Status          int              `json:"status" gorm:"not null;default:0;check:status IN (0, 1, 2)"`
	UpdatedBy       uint             `json:"updated_by" gorm:"not null"`
	Profil          *Profil          `json:"profil,omitempty" gorm:"foreignKey:ProfilID"`
	HutangPiutang   *HutangPiutang  `json:"hutang_piutang,omitempty" gorm:"foreignKey:HutangPiutangID"`
	Details         []PembelianDetail `json:"details,omitempty" gorm:"foreignKey:PembelianID"`
}

// PembelianDetail represents purchase line items
type PembelianDetail struct {
	BaseModel
	ID         uint      `json:"id" gorm:"primarykey"`
	PembelianID uint      `json:"pembelian_id" gorm:"not null;index"`
	BarangID   uint      `json:"barang_id" gorm:"not null;index"`
	Qty        int       `json:"qty" gorm:"not null"`
	HargaSatuan int       `json:"harga_satuan" gorm:"not null"`
	Diskon     int       `json:"diskon" gorm:"default:0"`
	Total      int       `json:"total" gorm:"not null"`
	UpdatedBy  uint      `json:"updated_by" gorm:"not null"`
	Pembelian  *Pembelian `json:"pembelian,omitempty" gorm:"foreignKey:PembelianID"`
	Barang     *Barang    `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
}

// HutangPiutang represents debts and credits
type HutangPiutang struct {
	BaseModel
	ID       uint   `json:"id" gorm:"primarykey"`
	Nomor    string `json:"nomor" gorm:"uniqueIndex;not null"`
	Tanggal  string `json:"tanggal" gorm:"not null"`
	ProfilID uint   `json:"profil_id" gorm:"index;not null"`
	Tipe     int    `json:"tipe" gorm:"not null;check:tipe IN (0, 1)"`
	Jumlah   int    `json:"jumlah" gorm:"not null"`
	Sisa     int    `json:"sisa" gorm:"not null"`
	Status   int    `json:"status" gorm:"not null;default:0;check:status IN (0, 1, 2);index"`
	UpdatedBy uint  `json:"updated_by" gorm:"not null"`
	Profil   *Profil `json:"profil,omitempty" gorm:"foreignKey:ProfilID"`
}

// InventoryBalance represents inventory stock balance
type InventoryBalance struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	BarangID uint   `json:"barang_id" gorm:"uniqueIndex;not null"`
	Qty      int    `json:"qty" gorm:"default:0"`
	QtyPO    int    `json:"qty_po" gorm:"default:0"`
	QtySO    int    `json:"qty_so" gorm:"default:0"`
	UpdatedAt time.Time `json:"updated_at"`
	Barang   *Barang `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
}

// KasBank represents cash and bank accounts
type KasBank struct {
	BaseModel
	ID            uint   `json:"id" gorm:"primarykey"`
	Nama          string `json:"nama" gorm:"not null"`
	NomorRekening string `json:"nomor_rekening"`
	Tipe          int    `json:"tipe" gorm:"not null"`
	Saldo         int    `json:"saldo" gorm:"default:0"`
	Status        int    `json:"status" gorm:"not null;default:1"`
	UpdatedBy     uint   `json:"updated_by" gorm:"not null"`
}

// Penerimaan represents income/receipts
type Penerimaan struct {
	BaseModel
	ID          uint   `json:"id" gorm:"primarykey"`
	Nomor       string `json:"nomor" gorm:"uniqueIndex;not null"`
	Tanggal     string `json:"tanggal" gorm:"not null;index"`
	KasBankID   *uint  `json:"kas_bank_id" gorm:"index"`
	KategoriID  *uint  `json:"kategori_id" gorm:"index"`
	ProfilID    *uint  `json:"profil_id" gorm:"index"`
	Jumlah      int    `json:"jumlah" gorm:"not null"`
	Keterangan  string `json:"keterangan"`
	Status      int    `json:"status" gorm:"not null;default:1"`
	UpdatedBy   uint   `json:"updated_by" gorm:"not null"`
	KasBank     *KasBank `json:"kas_bank,omitempty" gorm:"foreignKey:KasBankID"`
	Kategori    *KategoriPenerimaan `json:"kategori,omitempty" gorm:"foreignKey:KategoriID"`
	Profil      *Profil `json:"profil,omitempty" gorm:"foreignKey:ProfilID"`
}

// Pengeluaran represents expenses
type Pengeluaran struct {
	BaseModel
	ID          uint   `json:"id" gorm:"primarykey"`
	Nomor       string `json:"nomor" gorm:"uniqueIndex;not null"`
	Tanggal     string `json:"tanggal" gorm:"not null;index"`
	KasBankID   *uint  `json:"kas_bank_id" gorm:"index"`
	KategoriID  *uint  `json:"kategori_id" gorm:"index"`
	ProfilID    *uint  `json:"profil_id" gorm:"index"`
	Jumlah      int    `json:"jumlah" gorm:"not null"`
	Keterangan  string `json:"keterangan"`
	Status      int    `json:"status" gorm:"not null;default:1"`
	UpdatedBy   uint   `json:"updated_by" gorm:"not null"`
	KasBank     *KasBank `json:"kas_bank,omitempty" gorm:"foreignKey:KasBankID"`
	Kategori    *KategoriPengeluaran `json:"kategori,omitempty" gorm:"foreignKey:KategoriID"`
	Profil      *Profil `json:"profil,omitempty" gorm:"foreignKey:ProfilID"`
}

// KategoriPenerimaan represents income categories
type KategoriPenerimaan struct {
	BaseModel
	ID        uint   `json:"id" gorm:"primarykey"`
	Nama      string `json:"nama" gorm:"uniqueIndex;not null;collate:NOCASE"`
	UpdatedBy uint   `json:"updated_by" gorm:"not null"`
}

// KategoriPengeluaran represents expense categories
type KategoriPengeluaran struct {
	BaseModel
	ID        uint   `json:"id" gorm:"primarykey"`
	Nama      string `json:"nama" gorm:"uniqueIndex;not null;collate:NOCASE"`
	UpdatedBy uint   `json:"updated_by" gorm:"not null"`
}
