package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// ItemHandler holds the dependencies for Item APIs
type ItemHandler struct {
	itemRepo *repository.ItemRepository
}

// NewItemHandler creates a new handler instance
func NewItemHandler() *ItemHandler {
	return &ItemHandler{itemRepo: repository.NewItemRepository()}
}

// CreateItem handles POST /v1/items (Tambah/Hapus Daftar Belanja)
// Note: UserID saat ini di-hardcode. Perlu middleware Autentikasi di masa depan.
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req model.Item
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid, pastikan semua field terisi."})
		return
	}

	// 1. Logika Bisnis: Perhitungan Total Harga (sesuai TK1)
	req.TotalCost = float64(req.Quantity) * req.UnitPrice
	req.PurchasedDate = time.Now()

	// 2. Set User ID (Placeholder sementara, harus dari JWT token)
	req.UserID = 1 // ASUMSI: User ID 1 untuk testing

	// 3. Simpan ke Repository
	if err := h.itemRepo.CreateItem(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan item belanja."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item belanja berhasil ditambahkan.", "total_cost": req.TotalCost})
}
