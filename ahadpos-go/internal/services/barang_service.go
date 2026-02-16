package services

import (
	"ahadpos-go/internal/models"
	"ahadpos-go/internal/repository"
	"ahadpos-go/pkg/utils"
	"errors"
)

type BarangService struct {
	barangRepository *repository.BarangRepository
}

func NewBarangService(barangRepo *repository.BarangRepository) *BarangService {
	return &BarangService{
		barangRepository: barangRepo,
	}
}

type CreateBarangRequest struct {
	Barcode           string  `json:"barcode" validate:"required,max=30"`
	Nama              string  `json:"nama" validate:"required,max=100"`
	StrukturID        *uint   `json:"struktur_id,omitempty"`
	KategoriID        *uint   `json:"kategori_id,omitempty"`
	SatuanID          uint    `json:"satuan_id" validate:"required"`
	RakID             *uint   `json:"rak_id,omitempty"`
	RestockPoint      int     `json:"restock_point"`
	RestockLevel      int     `json:"restock_level"`
	RestockMin        int     `json:"restock_min"`
	VariantCoefficient float64 `json:"variant_coefficient"`
	Status            int     `json:"status"`
}

type UpdateBarangRequest struct {
	Nama              *string  `json:"nama,omitempty"`
	StrukturID        *uint    `json:"struktur_id,omitempty"`
	KategoriID        *uint    `json:"kategori_id,omitempty"`
	SatuanID          *uint    `json:"satuan_id,omitempty"`
	RakID             *uint    `json:"rak_id,omitempty"`
	RestockPoint      *int     `json:"restock_point,omitempty"`
	RestockLevel      *int     `json:"restock_level,omitempty"`
	RestockMin        *int     `json:"restock_min,omitempty"`
	VariantCoefficient *float64 `json:"variant_coefficient,omitempty"`
	Status            *int     `json:"status,omitempty"`
}

type GetBarangListRequest struct {
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	KategoriID uint             `json:"kategori_id,omitempty"`
	Status     int              `json:"status,omitempty"`
	Search     string           `json:"search,omitempty"`
}

func (s *BarangService) Create(req *CreateBarangRequest, updatedBy uint) (*models.Barang, error) {
	// Create barang model
	barang := &models.Barang{
		Barcode:           req.Barcode,
		Nama:              req.Nama,
		StrukturID:        req.StrukturID,
		KategoriID:        req.KategoriID,
		SatuanID:          req.SatuanID,
		RakID:             req.RakID,
		RestockPoint:      req.RestockPoint,
		RestockLevel:      req.RestockLevel,
		RestockMin:        req.RestockMin,
		VariantCoefficient: req.VariantCoefficient,
		Status:            req.Status,
		UpdatedBy:         updatedBy,
	}

	if barang.Status == 0 {
		barang.Status = 1 // Default to active
	}

	if err := s.barangRepository.Create(barang); err != nil {
		return nil, err
	}

	// Load relations
	barang, err := s.barangRepository.GetByID(barang.ID)
	if err != nil {
		return nil, err
	}

	return barang, nil
}

func (s *BarangService) GetByID(id uint) (*models.Barang, error) {
	barang, err := s.barangRepository.GetByID(id)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return barang, nil
}

func (s *BarangService) GetList(req *GetBarangListRequest) ([]models.Barang, int64, error) {
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.KategoriID > 0 {
		filters["kategori_id"] = req.KategoriID
	}
	if req.Status >= 0 {
		filters["status"] = req.Status
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}

	return s.barangRepository.GetAll(req.Page, req.PageSize, filters)
}

func (s *BarangService) Update(id uint, req *UpdateBarangRequest, updatedBy uint) (*models.Barang, error) {
	// Get existing barang
	barang, err := s.barangRepository.GetByID(id)
	if err != nil {
		return nil, utils.ErrNotFound
	}

	// Update fields
	if req.Nama != nil {
		barang.Nama = *req.Nama
	}
	if req.StrukturID != nil {
		barang.StrukturID = req.StrukturID
	}
	if req.KategoriID != nil {
		barang.KategoriID = req.KategoriID
	}
	if req.SatuanID != nil {
		barang.SatuanID = *req.SatuanID
	}
	if req.RakID != nil {
		barang.RakID = req.RakID
	}
	if req.RestockPoint != nil {
		barang.RestockPoint = *req.RestockPoint
	}
	if req.RestockLevel != nil {
		barang.RestockLevel = *req.RestockLevel
	}
	if req.RestockMin != nil {
		barang.RestockMin = *req.RestockMin
	}
	if req.VariantCoefficient != nil {
		barang.VariantCoefficient = *req.VariantCoefficient
	}
	if req.Status != nil {
		barang.Status = *req.Status
	}

	barang.UpdatedBy = updatedBy

	if err := s.barangRepository.Update(barang); err != nil {
		return nil, err
	}

	// Reload with relations
	return s.barangRepository.GetByID(barang.ID)
}

func (s *BarangService) Delete(id uint) error {
	barang, err := s.barangRepository.GetByID(id)
	if err != nil {
		return utils.ErrNotFound
	}

	if err := s.barangRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *BarangService) GetByBarcode(barcode string) (*models.Barang, error) {
	barang, err := s.barangRepository.GetByBarcode(barcode)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return barang, nil
}
