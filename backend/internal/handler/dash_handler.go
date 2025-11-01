package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	// Pastikan path impor ini sesuai dengan struktur proyek Anda
	"github.com/gusti3111/TKBMG/backend/internal/helper" // <-- 1. IMPORT HELPER
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// DashboardHandler menangani logika untuk endpoint dasbor
type DashboardHandler struct {
	itemRepo   *repository.ItemRepository
	budgetRepo *repository.BudgetRepository
	reportRepo *repository.ReportRepository
}

// NewDashboardHandler membuat instance DashboardHandler baru
func NewDashboardHandler(
	itemRepo *repository.ItemRepository,
	budgetRepo *repository.BudgetRepository,
	reportRepo *repository.ReportRepository,
) *DashboardHandler {
	return &DashboardHandler{
		itemRepo:   itemRepo,
		budgetRepo: budgetRepo,
		reportRepo: reportRepo,
	}
}

// GetDashboardSummary
// Ini adalah handler untuk endpoint: GET /api/v1/dashboard/summary
func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	// ======================================================
	// =============== PERBAIKAN 1: GUNAKAN HELPER ==========
	// ======================================================
	userID, ok := helper.GetUserID(c)
	if !ok {
		// helper.GetUserID(c) sudah mengirim respons error 401
		log.Println("[DashboardHandler] Gagal mengambil UserID dari helper")
		return
	}
	// ======================================================

	ctx := c.Request.Context()

	// 2. Dapatkan Budget Mingguan Terakhir
	now := time.Now()
	budget, err := h.budgetRepo.GetBudgetByDate(ctx, userID, now)

	var budgetAmount float64
	var startDate, endDate time.Time

	// ======================================================
	// =============== PERBAIKAN 2: SINKRONKAN TANGGAL ======
	// ======================================================
	if err != nil {
		// Ini bukan error fatal, mungkin user belum set budget
		log.Printf("Info: No budget found for user %d for this week: %v", userID, err)
		budgetAmount = 0.0
		// Jika tidak ada budget, hitung belanja 7 hari terakhir sebagai default
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -7) // 7 hari ke belakang
	} else if budget != nil {
		budgetAmount = budget.Amount
		// Gunakan rentang tanggal DARI BUDGET untuk menghitung belanja
		startDate = budget.StartDate
		endDate = budget.EndDate
	}
	// ======================================================

	// 3. Dapatkan Total Belanja Mingguan
	// (Memanggil fungsi dari item_repository.go)
	totalBelanja, err := h.itemRepo.GetTotalSpendingByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		log.Printf("Error getting total spending for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate spending"})
		return
	}

	// 4. Hitung Sisa Budget
	sisaBudget := budgetAmount - totalBelanja

	// 5. Siapkan Respons (Menggunakan model dari dashboard_models.go)
	summary := model.SummaryResponse{
		TotalBelanja: totalBelanja,
		Budget:       budgetAmount,
		SisaBudget:   sisaBudget,
	}

	// ======================================================
	// =============== PERBAIKAN 3: (KONFIRMASI) ============
	// ======================================================
	// Respons "data" ini SUDAH BENAR. Jangan diubah.
	// Frontend (SetBudget.jsx) harus membaca dari data.data.
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// GetDashboardCharts
// Ini adalah handler untuk endpoint: GET /api/v1/dashboard/charts
func (h *DashboardHandler) GetDashboardCharts(c *gin.Context) {
	// ======================================================
	// =============== PERBAIKAN 1: GUNAKAN HELPER ==========
	// ======================================================
	userID, ok := helper.GetUserID(c)
	if !ok {
		log.Println("[DashboardHandler] Gagal mengambil UserID dari helper")
		return
	}
	// ======================================================

	ctx := c.Request.Context()

	// 2. Tentukan rentang tanggal (misalnya, 30 hari terakhir untuk charts)
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0) // 1 bulan ke belakang

	// 3. Dapatkan data Pie Chart (Pengeluaran per Kategori)
	pieDataRepo, err := h.reportRepo.GetSpendingByCategory(ctx, userID, startDate, endDate)
	if err != nil {
		log.Printf("Error getting pie chart data for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pie chart data"})
		return
	}
	pieData := make([]model.PieChartItem, len(pieDataRepo))
	for i, item := range pieDataRepo {
		pieData[i] = model.PieChartItem{Name: item.Kategori, Value: item.Total}
	}

	// 4. Dapatkan data Bar Chart (Pengeluaran per Minggu - misal 4 minggu terakhir)
	barDataRepo, err := h.reportRepo.GetSpendingByWeek(ctx, userID, 4) // Ambil 4 minggu terakhir
	if err != nil {
		log.Printf("Error getting bar chart data for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bar chart data"})
		return
	}
	barData := make([]model.BarChartItem, len(barDataRepo))
	for i, item := range barDataRepo {
		barData[i] = model.BarChartItem{Name: item.MingguKe, Pengeluaran: item.Total}
	}

	// 5. Siapkan Respons
	charts := model.ChartResponse{
		PieChart: pieData,
		BarChart: barData,
	}

	// Respons "data" ini juga sudah benar
	c.JSON(http.StatusOK, gin.H{"data": charts})
}
