package handler

import (
	"log"
	"net/http"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// ItemHandler holds the dependencies for Item APIs
type ItemHandler struct {
	itemRepo *repository.ItemRepository
	catRepo  *repository.CategoryRepository // Dibutuhkan untuk validasi
}

// NewItemHandler (BUG #2 FIXED) - Sekarang menerima repo
// Ini akan dipanggil di main.go
func NewItemHandler(ir *repository.ItemRepository, cr *repository.CategoryRepository) *ItemHandler {
	return &ItemHandler{
		itemRepo: ir,
		catRepo:  cr,
	}
}

// CreateItem handles POST /v1/items (Tambah/Hapus Daftar Belanja)
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req model.Item

	// (PERINGATAN: Ini akan gagal jika Bug #1 (Tag JSON) belum diperbaiki di model/item.go)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	// 1. (BUG #4 FIXED) Ambil UserID dari middleware, BUKAN hardcode
	userID, ok := getUserID(c) // Menggunakan helper
	if !ok {
		return // Error sudah dikirim oleh getUserID
	}

	// 2. Logika Bisnis: Perhitungan Total Harga (sesuai TK1)
	req.TotalCost = float64(req.Quantity) * req.UnitPrice
	req.PurchasedDate = time.Now()
	req.UserID = userID // (PERBAIKAN) Set User ID dari token
	// req.CategoryID didapat dari binding JSON

	// 3. (Opsional tapi direkomendasikan) Validasi apakah id_kategori milik user
	// (Implementasi ini bisa ditambahkan di catRepo)
	if req.CategoryID != 0 {
		exists, err := h.itemRepo.CategoryExists(c.Request.Context(), req.CategoryID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa kategori"})
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}
	}

	// 4. Simpan ke Repository (Memperbaiki logika redundan)
	if err := h.itemRepo.CreateItem(c.Request.Context(), &req); err != nil {
		// Ini adalah tempat error "Gagal menyimpan item belanja" (Foreign Key)
		log.Printf("Error memanggil CreateItem repo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan item belanja.", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item belanja berhasil ditambahkan.", "data": req})
}

// GetItemsByUserID handles GET /v1/items
// (Menggantikan fungsi GetItems Anda yang lama agar sesuai dengan frontend)
func (h *ItemHandler) GetItemsByUserID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	items, err := h.itemRepo.GetItemsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil item belanja", "details": err.Error()})
		return
	}
	if items == nil {
		items = []model.Item{} // Kembalikan array kosong, bukan nil
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// UpdateItem handles PUT /v1/items/:id (Fungsi yang Hilang)
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	userID, ok := getUserID(c)
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

	// Set data dari token/URL
	req.ID = itemID
	req.UserID = userID
	// Hitung ulang total
	req.TotalCost = float64(req.Quantity) * req.UnitPrice

	if err := h.itemRepo.UpdateItem(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui item", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item berhasil diperbarui", "data": req})
}

// // DeleteItem handles DELETE /v1/items/:id (Fungsi yang Hilang)
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
		return
	}

	if err := h.itemRepo.DeleteItem(c.Request.Context(), itemID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus item", "details": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// getUserID adalah helper internal (disalin dari category_handler.go)
// (Sebaiknya dipindahkan ke paket 'helper' atau 'utils' umum)
func getUserID(c *gin.Context) (int, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Println("Handler Error: User ID not found in token context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terotentikasi"})
		return 0, false
	}
	userID, ok := userIDValue.(int)
	if !ok {
		log.Println("Handler Error: Invalid User ID format in token context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format User ID tidak valid"})
		return 0, false
	}
	return userID, true
}
