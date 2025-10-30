package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gusti3111/TKBMG/backend/internal/db"
	"github.com/gusti3111/TKBMG/backend/internal/model" // Menggunakan 'model' bukan 'models'
)

// BudgetRepository menangani operasi database untuk 'anggaran'
type BudgetRepository struct {
	db *sql.DB
}

// NewBudgetRepository membuat instance repository baru
func NewBudgetRepository() *BudgetRepository {
	// Asumsi db.DB adalah koneksi db global Anda yang diinisialisasi di paket db
	return &BudgetRepository{db: db.DB}
}

// Helper function untuk mendapatkan rentang minggu (Asumsi Minggu s/d Sabtu)
// Anda bisa sesuaikan logika ini jika minggu Anda dimulai hari Senin
func getWeekRange(date time.Time) (time.Time, time.Time) {
	// Temukan Minggu (hari 0)
	startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
	// Temukan Sabtu (hari 6)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	// Format ke YYYY-MM-DD 00:00:00 dan 23:59:59 untuk perbandingan SQL
	start := time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, date.Location())
	end := time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, date.Location())

	return start, end
}

// GetBudgetByDate mengambil budget yang aktif untuk user pada tanggal tertentu
// Ini adalah fungsi yang akan dipanggil oleh GetDashboardSummary
func (r *BudgetRepository) GetBudgetByDate(ctx context.Context, userID int, date time.Time) (*model.Budget, error) {
	query := `SELECT id_anggaran, id_user, start_date, end_date, jumlah_anggaran
	          FROM anggaran 
	          WHERE id_user = $1 AND $2 BETWEEN start_date AND end_date
	          ORDER BY start_date DESC
	          LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, userID, date)
	var budget model.Budget

	err := row.Scan(
		&budget.ID,
		&budget.UserID,
		&budget.StartDate,
		&budget.EndDate,
		&budget.Amount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Ini bukan error fatal, hanya berarti tidak ada budget yang di-set
			log.Printf("Info: No budget found for user %d in this date range", userID)
			return nil, fmt.Errorf("no budget found for this date range")
		}
		log.Printf("Error scanning budget: %v", err)
		return nil, fmt.Errorf("failed to scan budget: %w", err)
	}

	return &budget, nil
}

// UpsertBudgetForCurrentWeek membuat atau memperbarui budget untuk minggu ini
// Ini akan dipanggil oleh handler Halaman "Set Budget" (POST /api/v1/budgets)
func (r *BudgetRepository) UpsertBudgetForCurrentWeek(ctx context.Context, userID int, amount float64) error {
	// Tentukan awal dan akhir minggu ini
	startOfWeek, endOfWeek := getWeekRange(time.Now())

	// 1. Cek apakah budget untuk minggu ini sudah ada
	var existingID int
	checkQuery := `SELECT id_anggaran FROM anggaran WHERE id_user = $1 AND start_date = $2`

	err := r.db.QueryRowContext(ctx, checkQuery, userID, startOfWeek).Scan(&existingID)

	// 2. Jika tidak ada (ErrNoRows), INSERT
	if err == sql.ErrNoRows {
		insertQuery := `INSERT INTO anggaran (id_user, start_date, end_date, jumlah_anggaran)
		                VALUES ($1, $2, $3, $4)`
		_, errInsert := r.db.ExecContext(ctx, insertQuery, userID, startOfWeek, endOfWeek, amount)
		if errInsert != nil {
			log.Printf("Error inserting new budget: %v", errInsert)
			return fmt.Errorf("failed to insert budget: %w", errInsert)
		}
		log.Printf("Successfully INSERTED budget for user %d", userID)
		return nil
	}

	// 3. Jika ada error lain saat mengecek
	if err != nil {
		log.Printf("Error checking existing budget: %v", err)
		return fmt.Errorf("failed to check budget: %w", err)
	}

	// 4. Jika ada (tidak error), UPDATE
	updateQuery := `UPDATE anggaran SET jumlah_anggaran = $1, end_date = $2
	                WHERE id_anggaran = $3 AND id_user = $4`
	_, errUpdate := r.db.ExecContext(ctx, updateQuery, amount, endOfWeek, existingID, userID)
	if errUpdate != nil {
		log.Printf("Error updating existing budget: %v", errUpdate)
		return fmt.Errorf("failed to update budget: %w", errUpdate)
	}

	log.Printf("Successfully UPDATED budget for user %d", userID)
	return nil
}
