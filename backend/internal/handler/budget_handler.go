package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gusti3111/TKBMG/backend/internal/helper"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// BudgetHandler menangani logika HTTP untuk Anggaran.
type BudgetHandler struct {
	repo *repository.BudgetRepository
}

// NewBudgetHandler membuat instance BudgetHandler baru.
func NewBudgetHandler(r *repository.BudgetRepository) *BudgetHandler {
	return &BudgetHandler{repo: r}
}

// SetBudget menangani POST /api/v1/budgets
// Ini akan memanggil 'UpsertBudgetForCurrentWeek' dari repository.
func (h *BudgetHandler) SetBudget(c *gin.Context) {
	// 1. Ambil UserID dari token
	userID, ok := helper.GetUserID(c)
	if !ok {
		return // Helper sudah mengirim respons 401
	}

	// 2. Bind JSON body
	// Kita hanya perlu 'jumlah_anggaran' dari user.
	// Repo Anda (UpsertBudgetForCurrentWeek) sudah pintar
	// menghitung start_date dan end_date (minggu ini) secara otomatis.
	var req struct {
		Amount float64 `json:"jumlah_anggaran" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah anggaran harus lebih besar dari 0"})
		return
	}

	// 3. Panggil Repository (logika inti)
	err := h.repo.UpsertBudgetForCurrentWeek(c.Request.Context(), userID, req.Amount)
	if err != nil {
		log.Printf("[BudgetHandler] Gagal upsert budget: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan anggaran", "details": err.Error()})
		return
	}

	// 4. Kirim respons sukses
	c.JSON(http.StatusOK, gin.H{
		"message":         "Anggaran untuk minggu ini berhasil disimpan/diperbarui",
		"jumlah_anggaran": req.Amount,
	})
}
