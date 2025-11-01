package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gusti3111/TKBMG/backend/internal/helper"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
	"github.com/xuri/excelize/v2" // <-- 1. Import excelize
)

// ReportHandler menangani logika HTTP untuk Laporan.
type ReportHandler struct {
	reportRepo *repository.ReportRepository
}

// NewReportHandler membuat instance ReportHandler baru.
func NewReportHandler(r *repository.ReportRepository) *ReportHandler {
	return &ReportHandler{reportRepo: r}
}

// GenerateReport menangani GET /api/v1/reports/download
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	// 1. Ambil UserID dari token
	userID, ok := helper.GetUserID(c)
	if !ok {
		return // Helper sudah mengirim respons 401
	}

	// 2. (Opsional) Ambil query param
	// Contoh: /api/v1/reports/download?type=excel&weeks=4
	reportType := c.DefaultQuery("type", "excel") // Default ke excel
	numWeeksStr := c.DefaultQuery("weeks", "4")
	numWeeks, _ := strconv.Atoi(numWeeksStr)

	// 3. Ambil data dari Repository
	// Kita panggil fungsi yang sama dengan yang dipakai dashboard
	barData, err := h.reportRepo.GetSpendingByWeek(c.Request.Context(), userID, numWeeks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan"})
		return
	}

	// 4. Cek tipe laporan yang diminta
	if reportType == "excel" {
		// Panggil fungsi helper untuk membuat file Excel
		buffer, err := h.createExcelReport(barData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat file Excel"})
			return
		}

		// 5. Kirim file ke user
		fileName := fmt.Sprintf("Laporan_Mingguan_%d.xlsx", time.Now().Unix())

		// Set Header HTTP agar browser men-download file
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())

	} else if reportType == "pdf" {
		// Logika untuk PDF (menggunakan gofpdf) akan ada di sini
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Download PDF belum diimplementasikan"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe laporan tidak didukung"})
	}
}

// createExcelReport adalah helper untuk men-generate file Excel
func (h *ReportHandler) createExcelReport(data []model.SpendingByWeek) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Laporan Mingguan"
	index, _ := f.NewSheet(sheetName) // Buat sheet baru

	// Set Header Tabel
	f.SetCellValue(sheetName, "A1", "Minggu Ke")
	f.SetCellValue(sheetName, "B1", "Total Pengeluaran")

	// Set Style untuk Header
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", "B1", style)

	// Isi data laporan
	for i, item := range data {
		row := i + 2 // Mulai dari baris 2
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), item.MingguKe)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), item.Total)

		// Set format mata uang (Contoh: Rp 123.456)
		// Anda bisa buat ini lebih kompleks
		f.SetCellStyle(sheetName, "B"+strconv.Itoa(row), "B"+strconv.Itoa(row),
			style,
		)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // Hapus sheet default

	// Simpan ke buffer (memory)
	buffer, err := f.WriteToBuffer()
	if err != nil {
		log.Printf("Gagal menulis excel ke buffer: %v", err)
		return nil, err
	}

	return buffer, nil
}
