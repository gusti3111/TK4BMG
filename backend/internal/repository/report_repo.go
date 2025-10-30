package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gusti3111/TKBMG/backend/internal/db"
	"github.com/gusti3111/TKBMG/backend/internal/model"
)

// ReportRepository menangani operasi database yang kompleks untuk laporan/agregasi
type ReportRepository struct {
	db *sql.DB
}

// NewReportRepository membuat instance repository baru
func NewReportRepository() *ReportRepository {
	return &ReportRepository{db: db.DB}
}

// GetSpendingByCategory menghitung total pengeluaran per kategori
// Ini dipanggil oleh GetDashboardCharts untuk Pie Chart
func (r *ReportRepository) GetSpendingByCategory(ctx context.Context, userID int, startDate time.Time, endDate time.Time) ([]model.SpendingByCategory, error) {
	// (Query ini mengasumsikan Anda memiliki tabel 'referensi_kategori' sesuai ERD TK2)
	query := `
		SELECT 
			COALESCE(rk.nama_kategori, 'Tanpa Kategori') as nama_kategori, 
			SUM(i.total_harga) as total
		FROM 
			items i
		LEFT JOIN 
			referensi_kategori rk ON i.id_kategori = rk.id_kategori
		WHERE 
			i.id_user = $1 AND i.purchased_date BETWEEN $2 AND $3
		GROUP BY 
			rk.nama_kategori
		ORDER BY 
			total DESC`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		log.Printf("Error querying spending by category: %v", err)
		return nil, fmt.Errorf("failed to get pie chart data: %w", err)
	}
	defer rows.Close()

	var results []model.SpendingByCategory
	for rows.Next() {
		var item model.SpendingByCategory
		if err := rows.Scan(&item.Kategori, &item.Total); err != nil {
			log.Printf("Error scanning category spending: %v", err)
			continue
		}
		results = append(results, item)
	}
	return results, rows.Err()
}

// GetSpendingByWeek menghitung total pengeluaran per minggu (4 minggu terakhir)
// Ini dipanggil oleh GetDashboardCharts untuk Bar Chart
func (r *ReportRepository) GetSpendingByWeek(ctx context.Context, userID int, numWeeks int) ([]model.SpendingByWeek, error) {
	// Query ini spesifik untuk PostgreSQL (menggunakan TO_CHAR).
	// Jika Anda menggunakan MySQL, ganti 'TO_CHAR' dengan 'WEEK' atau 'DATE_FORMAT'
	// Jika Anda menggunakan SQLite, ini akan lebih rumit (perlu 'strftime')
	// Asumsi PostgreSQL:
	query := `
		SELECT 
			TO_CHAR(purchased_date, 'YYYY-WW') as minggu, 
			SUM(total_harga) as total
		FROM 
			items
		WHERE 
			id_user = $1 AND purchased_date >= $2
		GROUP BY 
			minggu
		ORDER BY 
			minggu DESC
		LIMIT $3`

	// Tentukan tanggal mulai (misal: 4 minggu * 7 hari = 28 hari lalu)
	startDate := time.Now().AddDate(0, 0, -(numWeeks * 7))

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, numWeeks)
	if err != nil {
		log.Printf("Error querying spending by week: %v", err)
		return nil, fmt.Errorf("failed to get bar chart data: %w", err)
	}
	defer rows.Close()

	var results []model.SpendingByWeek
	for rows.Next() {
		var item model.SpendingByWeek
		if err := rows.Scan(&item.MingguKe, &item.Total); err != nil {
			log.Printf("Error scanning weekly spending: %v", err)
			continue
		}
		results = append(results, item)
	}
	return results, rows.Err()
}
