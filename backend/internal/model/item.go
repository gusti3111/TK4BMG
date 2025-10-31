package model

import (
	"database/sql"
	"time"
)

// Item represents the data structure for the "Items" entity
// (BUG #1 FIXED: Mengganti json:\"nama_item\" menjadi json:"nama_item")
type Item struct {
	ID            int            `json:"id_item"`
	UserID        int            `json:"id_user"`
	CategoryID    int            `json:"id_kategori"`
	ItemName      string         `json:"nama_item" binding:"required"`
	Quantity      int            `json:"jumlah_item" binding:"required"`
	UnitPrice     float64        `json:"harga_satuan" binding:"required"`
	TotalCost     float64        `json:"total_harga"` // Dihitung di backend
	PurchasedDate time.Time      `json:"purchased_date"`
	CategoryName  sql.NullString `json:"nama_kategori,omitempty"` // Untuk join
}
type ItemRequest struct {
	CategoryID int     `json:"id_kategori" binding:"required"`
	ItemName   string  `json:"nama_item" binding:"required"`
	Quantity   int     `json:"jumlah_item" binding:"required"`
	UnitPrice  float64 `json:"harga_satuan" binding:"required"`
}

// Budget represents the data structure for the "Anggaran" entity
type Budget struct {
	ID        int       `json:"id_anggaran"`
	UserID    int       `json:"id_user"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Amount    float64   `json:"jumlah_anggaran" binding:"required"`
}

// Category represents the data structure for "Referensi_Kategori"
type Category struct {
	ID           int    `json:"id_kategori"`
	UserID       int    `json:"id_user"`
	CategoryName string `json:"nama_kategori" binding:"required"`
}

// === DTO (Data Transfer Objects) untuk Laporan/Dasbor ===

// SpendingByCategory adalah struct untuk data Pie Chart
type SpendingByCategory struct {
	Kategori string  `json:"kategori" db:"nama_kategori"`
	Total    float64 `json:"total" db:"total"`
}

// SpendingByWeek adalah struct untuk data Bar Chart
type SpendingByWeek struct {
	MingguKe string  `json:"minggu_ke" db:"minggu"` // Contoh: "W40"
	Total    float64 `json:"total" db:"total"`
}
type SummaryResponse struct {
	TotalBelanja float64 `json:"total_belanja"`
	Budget       float64 `json:"budget"`
	SisaBudget   float64 `json:"sisa_budget"`
}

// PieChartItem adalah DTO untuk satu potong data di Pie Chart.
type PieChartItem struct {
	Name  string  `json:"name"`  // Nama Kategori
	Value float64 `json:"value"` // Total pengeluaran
}

// BarChartItem adalah DTO untuk satu batang data di Bar Chart.
type BarChartItem struct {
	Name        string  `json:"name"`        // Nama minggu (misal: "W40")
	Pengeluaran float64 `json:"pengeluaran"` // Total pengeluaran
}

// ChartResponse adalah DTO pembungkus untuk kedua data chart.
type ChartResponse struct {
	PieChart []PieChartItem `json:"pie_chart"`
	BarChart []BarChartItem `json:"bar_chart"`
}
