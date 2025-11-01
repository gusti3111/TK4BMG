package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	// Import internal packages
	"github.com/gusti3111/TKBMG/backend/internal/helper"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// CategoryHandler menangani logika HTTP untuk referensi_kategori.
type CategoryHandler struct {
	repo *repository.CategoryRepository
}

// NewCategoryHandler membuat instance CategoryHandler baru.
func NewCategoryHandler(r *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: r}
}

// ======================================================================
// CREATE CATEGORY (POST /api/v1/kategori)
// ======================================================================
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// Ambil user ID dari token JWT (via middleware)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return // Sudah mengembalikan response 401 dari helper
	}

	var req model.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CategoryHandler] Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid", "details": err.Error()})
		return
	}

	// Validasi input
	if req.CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori tidak boleh kosong"})
		return
	}

	req.UserID = userID

	// Simpan ke database via repository
	if err := h.repo.CreateKategori(c.Request.Context(), &req); err != nil {
		log.Printf("[CategoryHandler] Gagal membuat kategori: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil ditambahkan",
		"data":    req,
	})
}

// ======================================================================
// GET ALL CATEGORY (GET /api/v1/kategori)
// ======================================================================
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	kategoriList, err := h.repo.GetKategoriByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[CategoryHandler] Gagal mengambil daftar kategori: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar kategori", "details": err.Error()})
		return
	}

	if kategoriList == nil {
		kategoriList = []model.Category{}
	}

	c.JSON(http.StatusOK, gin.H{"data": kategoriList})
}

// ======================================================================
// UPDATE CATEGORY (PUT /api/v1/kategori/:id)
// ======================================================================
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	kategoriID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid"})
		return
	}

	var req model.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid", "details": err.Error()})
		return
	}

	if req.CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori tidak boleh kosong"})
		return
	}

	req.ID = kategoriID
	req.UserID = userID

	if err := h.repo.UpdateKategori(c.Request.Context(), &req); err != nil {
		log.Printf("[CategoryHandler] Error update kategori: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil diperbarui",
		"data":    req,
	})
}

// ======================================================================
// DELETE CATEGORY (DELETE /api/v1/kategori/:id)
// ======================================================================
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}

	kategoriID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid"})
		return
	}

	// TODO (opsional): Cek apakah kategori digunakan oleh tabel Items sebelum dihapus

	if err := h.repo.DeleteKategori(c.Request.Context(), kategoriID, userID); err != nil {
		log.Printf("[CategoryHandler] Error delete kategori: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}
