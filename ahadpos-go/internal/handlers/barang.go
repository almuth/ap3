package handlers

import (
	"net/http"
	"strconv"

	"github.com/ahadpos/go/internal/models"
	"github.com/gin-gonic/gin"
)

// BarangHandler handles product-related requests
type BarangHandler struct{}

func NewBarangHandler() *BarangHandler {
	return &BarangHandler{}
}

// CreateBarang creates a new product
func (h *BarangHandler) CreateBarang(c *gin.Context) {
	var barang models.Barang
	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	barang.UpdatedBy = userID.(uint)

	if err := db.Create(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, barang)
}

// GetBarang retrieves a product by ID
func (h *BarangHandler) GetBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var barang models.Barang
	if err := db.Preload("Kategori").Preload("Satuan").Preload("Rak").First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, barang)
}

// GetBarangList retrieves a paginated list of products
func (h *BarangHandler) GetBarangList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	var barang []models.Barang

	query := db.Model(&models.Barang{})
	
	// Filter by category
	if kategoriID := c.Query("kategori_id"); kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Search by name or barcode
	if search := c.Query("search"); search != "" {
		query = query.Where("nama LIKE ? OR barcode LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("Kategori").Preload("Satuan").Preload("Rak").
		Offset(offset).Limit(pageSize).Find(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     barang,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
		"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdateBarang updates an existing product
func (h *BarangHandler) UpdateBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var barang models.Barang
	if err := db.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var updateData models.Barang
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	updateData.UpdatedBy = userID.(uint)

	if err := db.Model(&barang).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, barang)
}

// DeleteBarang deletes a product
func (h *BarangHandler) DeleteBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := db.Delete(&models.Barang{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
