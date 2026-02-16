package handlers

import (
	"ahadpos-go/internal/services"
	"ahadpos-go/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BarangHandler struct {
	barangService *services.BarangService
}

func NewBarangHandler(barangService *services.BarangService) *BarangHandler {
	return &BarangHandler{barangService: barangService}
}

// CreateBarang creates a new product
func (h *BarangHandler) CreateBarang(c *gin.Context) {
	var req services.CreateBarangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.ErrValidationFailed))
		return
	}

	userID, _ := c.Get("user_id")
	barang, err := h.barangService.Create(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, utils.NewSuccessResponse(barang))
}

// GetBarang retrieves a product by ID
func (h *BarangHandler) GetBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.ErrValidationFailed))
		return
	}

	barang, err := h.barangService.GetByID(uint(id))
	if err != nil {
		c.JSON(utils.HTTPStatus(err), utils.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(barang))
}

// GetBarangList retrieves a paginated list of products
func (h *BarangHandler) GetBarangList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var kategoriID uint
	if kategoriIDStr := c.Query("kategori_id"); kategoriIDStr != "" {
		kategoriID, _ = strconv.ParseUint(kategoriIDStr, 10, 32)
	}

	var status int = -1
	if statusStr := c.Query("status"); statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	req := &services.GetBarangListRequest{
		Page:       page,
		PageSize:   pageSize,
		KategoriID: uint(kategoriID),
		Status:     status,
		Search:     c.Query("search"),
	}

	barang, total, err := h.barangService.GetList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	meta := &utils.Meta{
		Page:      page,
		PageSize:   pageSize,
		Total:     total,
		TotalPages: int(totalPages),
	}

	c.JSON(http.StatusOK, utils.NewPaginatedResponse(barang, meta))
}

// UpdateBarang updates an existing product
func (h *BarangHandler) UpdateBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.ErrValidationFailed))
		return
	}

	var req services.UpdateBarangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.ErrValidationFailed))
		return
	}

	userID, _ := c.Get("user_id")
	barang, err := h.barangService.Update(uint(id), &req, userID.(uint))
	if err != nil {
		c.JSON(utils.HTTPStatus(err), utils.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(barang))
}

// DeleteBarang deletes a product
func (h *BarangHandler) DeleteBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.ErrValidationFailed))
		return
	}

	if err := h.barangService.Delete(uint(id)); err != nil {
		c.JSON(utils.HTTPStatus(err), utils.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(gin.H{"message": "Product deleted successfully"}))
}
