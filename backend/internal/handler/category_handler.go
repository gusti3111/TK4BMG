package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// Pastikan path impor ini sesuai dengan struktur proyek Anda
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// CategoryHandler menangani logika HTTP untuk referensi_kategori.
type CategoryHandler struct {
	repo *repository.CategoryRepository
}

// NewCategoryHandler membuat instance CategoryHandler baru.
// Anda akan memanggil ini di main.go dan menyuntikkan (inject) repository.
func NewCategoryHandler(r *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: r}
}

// --- METODE-METODE HANDLER (CRUD) ---

// CreateCategory
// Endpoint: POST /api/v1/kategori
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// 1. Ambil UserID dari token (yang disimpan oleh middleware)
	userID, ok := getUserID(c)
	if !ok {
		return // Error response sudah dikirim oleh getUserID
	}

	// 2. Bind JSON request body ke struct
	var req model.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid", "details": err.Error()})
		return
	}

	// 3. Validasi dasar
	if req.CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori tidak boleh kosong"})
		return
	}

	// 4. Set UserID (dari token) ke struct
	req.UserID = userID

	// 5. Panggil Repository untuk menyimpan ke DB
	if err := h.repo.CreateKategori(c.Request.Context(), &req); err != nil {
		log.Printf("Error memanggil CreateKategori repo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori", "details": err.Error()})
		return
	}

	// 6. Kirim respons sukses
	c.JSON(http.StatusCreated, gin.H{"message": "Kategori berhasil ditambahkan", "data": req})
}

// GetCategories
// Endpoint: GET /api/v1/kategori
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	// 1. Ambil UserID dari token
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	// 2. Panggil Repository untuk mengambil data
	kategoriList, err := h.repo.GetKategoriByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error memanggil GetKategoriByUserID repo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar kategori", "details": err.Error()})
		return
	}

	// 3. Kirim respons sukses
	// Pastikan mengirim array kosong jika tidak ada data, bukan nil
	if kategoriList == nil {
		kategoriList = []model.Category{}
	}

	c.JSON(http.StatusOK, gin.H{"data": kategoriList})
}

// UpdateCategory
// Endpoint: PUT /api/v1/kategori/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	// 1. Ambil UserID dari token
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	// 2. Ambil ID Kategori dari parameter URL
	kategoriID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid"})
		return
	}

	// 3. Bind JSON request body
	var req model.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid", "details": err.Error()})
		return
	}

	// 4. Set ID dan UserID (dari token/URL) ke struct
	req.ID = kategoriID
	req.UserID = userID

	// 5. Panggil Repository untuk update
	if err := h.repo.UpdateKategori(c.Request.Context(), &req); err != nil {
		log.Printf("Error memanggil UpdateKategori repo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori", "details": err.Error()})
		return
	}

	// 6. Kirim respons sukses
	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil diperbarui", "data": req})
}

// DeleteCategory
// Endpoint: DELETE /api/v1/kategori/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	// 1. Ambil UserID dari token
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	// 2. Ambil ID Kategori dari parameter URL
	kategoriID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid"})
		return
	}

	// PERHATIAN (TK2): ERD Anda menunjukkan relasi antara Items dan Referensi_Kategori.
	// Jika Anda langsung menghapus kategori yang sedang dipakai, database akan error (Foreign Key Constraint).
	// TODO: Tambahkan logika di service/handler untuk:
	// 1. Memeriksa apakah kategori ini sedang dipakai oleh 'items'.
	// 2. Jika ya, GAGALKAN penghapusan ATAU set 'id_kategori' di 'items' menjadi NULL.

	// 3. Panggil Repository untuk menghapus
	if err := h.repo.DeleteKategori(c.Request.Context(), kategoriID, userID); err != nil {
		log.Printf("Error memanggil DeleteKategori repo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori", "details": err.Error()})
		return
	}

	// 4. Kirim respons sukses (No Content)
	c.JSON(http.StatusNoContent, nil)
}

// --- HELPER ---
