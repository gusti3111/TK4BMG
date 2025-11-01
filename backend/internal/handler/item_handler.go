package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gusti3111/TKBMG/backend/internal/helper"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// ItemHandler menangani operasi HTTP untuk tabel items
type ItemHandler struct {
	repo         *repository.ItemRepository
	categoryRepo *repository.CategoryRepository
}

// NewItemHandler membuat handler baru
func NewItemHandler(repo *repository.ItemRepository, categoryRepo *repository.CategoryRepository) *ItemHandler {
	return &ItemHandler{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// ======================================================================
// CREATE ITEM (POST /api/v1/items)
// ======================================================================
func (h *ItemHandler) CreateItem(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	var req model.Item
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ItemHandler] Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	if req.ItemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama item tidak boleh kosong"})
		return
	}

	// Pastikan kategori valid jika ada
	if req.CategoryID != 0 {
		exists, err := h.repo.CategoryExists(c.Request.Context(), req.CategoryID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa kategori"})
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}
	}

	req.UserID = userID
	req.PurchasedDate = time.Now()
	req.TotalCost = float64(req.Quantity) * req.UnitPrice

	// NOTE: Repository hanya menyimpan nama_item, jumlah_item, harga_satuan
	// Kolom lain (userID, kategori, total_harga, purchased_date) belum disimpan di DB
	if err := h.repo.CreateItem(c.Request.Context(), &req); err != nil {
		log.Printf("[ItemHandler] Error saving item: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Item berhasil ditambahkan",
		"data":    req,
	})
}

// ======================================================================
// GET ITEMS (GET /api/v1/items)
// ======================================================================
func (h *ItemHandler) GetItems(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	items, err := h.repo.GetItemsByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[ItemHandler] Error fetching items: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar item"})
		return
	}

	if items == nil {
		items = []model.Item{}
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ======================================================================
// UPDATE ITEM (PUT /api/v1/items/:id)
// ======================================================================
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
		return
	}

	var req model.Item
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	req.ID = itemID
	req.UserID = userID
	req.TotalCost = float64(req.Quantity) * req.UnitPrice
	req.PurchasedDate = time.Now()

	if err := h.repo.UpdateItem(c.Request.Context(), &req); err != nil {
		log.Printf("[ItemHandler] Error updating item: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item berhasil diperbarui", "data": req})
}

// ======================================================================
// DELETE ITEM (DELETE /api/v1/items/:id)
// ======================================================================
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
		return
	}

	if err := h.repo.DeleteItem(c.Request.Context(), itemID, userID); err != nil {
		log.Printf("[ItemHandler] Error deleting item: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item berhasil dihapus"})
}
