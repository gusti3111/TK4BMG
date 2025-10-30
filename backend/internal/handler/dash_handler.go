package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	// Pastikan path impor ini sesuai dengan struktur proyek Anda
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

// DashboardHandler menangani logika untuk endpoint dasbor
type DashboardHandler struct {
	itemRepo   *repository.ItemRepository
	budgetRepo *repository.BudgetRepository
	reportRepo *repository.ReportRepository
	// Anda bisa tambahkan repo lain di sini jika diperlukan
}

// NewDashboardHandler membuat instance DashboardHandler baru
// Anda akan memanggil ini di main.go saat menginisialisasi handler
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
	// 1. Dapatkan ID pengguna dari middleware JWT (AuthMiddleware)
	// Kita asumsikan middleware Anda menyimpan 'userID' di context
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Println("Error: User ID not found in token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	// Lakukan konversi tipe data dengan aman
	userID, ok := userIDValue.(int) // Sesuaikan tipe data ini (int, int64, string) dengan apa yang Anda simpan di token
	if !ok {
		log.Println("Error: Invalid User ID format in token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format"})
		return
	}

	ctx := c.Request.Context()

	// 2. Dapatkan Budget Mingguan Terakhir
	// (Memanggil fungsi dari budget_repository.go)
	// Kita tentukan rentang minggu ini
	now := time.Now()
	budget, err := h.budgetRepo.GetBudgetByDate(ctx, userID, now)
	var budgetAmount float64
	if err != nil {
		// Ini bukan error fatal, mungkin user belum set budget
		log.Printf("Info: No budget found for user %d for this week: %v", userID, err)
		budgetAmount = 0.0
	} else if budget != nil {
		budgetAmount = budget.Amount
	}

	// 3. Dapatkan Total Belanja Mingguan
	// (Memanggil fungsi dari item_repository.go)
	// Tentukan rentang tanggal (misalnya, 7 hari terakhir)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7) // 7 hari ke belakang

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

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// GetDashboardCharts
// Ini adalah handler untuk endpoint: GET /api/v1/dashboard/charts
func (h *DashboardHandler) GetDashboardCharts(c *gin.Context) {
	// 1. Dapatkan ID pengguna dari middleware JWT
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Println("Error: User ID not found in token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		log.Println("Error: Invalid User ID format in token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User ID format"})
		return
	}

	ctx := c.Request.Context()

	// 2. Tentukan rentang tanggal (misalnya, 30 hari terakhir untuk charts)
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0) // 1 bulan ke belakang

	// 3. Dapatkan data Pie Chart (Pengeluaran per Kategori)
	// (Memanggil fungsi dari report_repository.go)
	pieDataRepo, err := h.reportRepo.GetSpendingByCategory(ctx, userID, startDate, endDate)
	if err != nil {
		log.Printf("Error getting pie chart data for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pie chart data"})
		return
	}
	// Konversi tipe data dari Repo ke Model (DTO)
	pieData := make([]model.PieChartItem, len(pieDataRepo))
	for i, item := range pieDataRepo {
		pieData[i] = model.PieChartItem{Name: item.Kategori, Value: item.Total}
	}

	// 4. Dapatkan data Bar Chart (Pengeluaran per Minggu - misal 4 minggu terakhir)
	// (Memanggil fungsi dari report_repository.go)
	barDataRepo, err := h.reportRepo.GetSpendingByWeek(ctx, userID, 4) // Ambil 4 minggu terakhir
	if err != nil {
		log.Printf("Error getting bar chart data for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bar chart data"})
		return
	}
	// Konversi tipe data dari Repo ke Model (DTO)
	barData := make([]model.BarChartItem, len(barDataRepo))
	for i, item := range barDataRepo {
		barData[i] = model.BarChartItem{Name: item.MingguKe, Pengeluaran: item.Total}
	}

	// 5. Siapkan Respons (Menggunakan model dari dashboard_models.go)
	charts := model.ChartResponse{
		PieChart: pieData,
		BarChart: barData,
	}

	c.JSON(http.StatusOK, gin.H{"data": charts})
}
