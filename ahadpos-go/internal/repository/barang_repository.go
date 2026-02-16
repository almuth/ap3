package repository

import (
	"ahadpos-go/internal/models"
	"ahadpos-go/pkg/database"
)

type BarangRepository struct {
	*BaseRepository
}

func NewBarangRepository() *BarangRepository {
	return &BarangRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *BarangRepository) Create(barang *models.Barang) error {
	return r.DB.Create(barang).Error
}

func (r *BarangRepository) GetByID(id uint) (*models.Barang, error) {
	var barang models.Barang
	err := r.DB.Preload("Kategori").Preload("Satuan").Preload("Rak").Preload("Struktur").
		First(&barang, id).Error
	if err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *BarangRepository) GetAll(page, pageSize int, filters map[string]interface{}) ([]models.Barang, int64, error) {
	var barang []models.Barang
	var total int64

	query := r.DB.Model(&models.Barang{})

	// Apply filters
	if kategoriID, ok := filters["kategori_id"]; ok {
		query = query.Where("kategori_id = ?", kategoriID)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if search, ok := filters["search"]; ok {
		query = query.Where("nama LIKE ? OR barcode LIKE ?", "%"+search.(string)+"%", "%"+search.(string)+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch data
	offset := (page - 1) * pageSize
	err := query.Preload("Kategori").Preload("Satuan").Preload("Rak").
		Offset(offset).Limit(pageSize).
		Find(&barang).Error

	if err != nil {
		return nil, 0, err
	}

	return barang, total, nil
}

func (r *BarangRepository) Update(barang *models.Barang) error {
	return r.DB.Save(barang).Error
}

func (r *BarangRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Barang{}, id).Error
}

func (r *BarangRepository) GetByBarcode(barcode string) (*models.Barang, error) {
	var barang models.Barang
	err := r.DB.Preload("Kategori").Preload("Satuan").Preload("Rak").
		Where("barcode = ?", barcode).
		First(&barang).Error
	if err != nil {
		return nil, err
	}
	return &barang, nil
}
