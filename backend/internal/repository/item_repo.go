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

// ItemRepository handles database operations related to Item and Budget
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new repository instance
func NewItemRepository() *ItemRepository {
	return &ItemRepository{db: db.DB}
}

// CreateItem saves a new item into the Items table
func (r *ItemRepository) CreateItem(ctx context.Context, item *model.Item) error {
	categoryID := item.CategoryID
	if categoryID == 0 {
		categoryID = 1 // default kategori
	}

	query := `
		INSERT INTO items (id_user, id_kategori, nama_item, jumlah_item, harga_satuan)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.UserID,
		item.CategoryID,
		item.ItemName,
		item.Quantity,
		item.UnitPrice,
		item.TotalCost,
		item.PurchasedDate,
	)

	if err != nil {
		log.Printf("Error inserting item: %v", err)
		return fmt.Errorf("failed to save shopping item: %w", err)
	}

	return nil
}

// GetItemsByUserID fetches all shopping items for a specific user within a timeframe (simple version)
func (r *ItemRepository) GetItemsByUserID(ctx context.Context, userID int) ([]model.Item, error) {
	// Query ini bisa dioptimalkan dengan filter tanggal di masa depan (TK4 Rework)
	query := `SELECT id_item, id_kategori, nama_item, jumlah_item, harga_satuan FROM items WHERE id_user = $1 ORDER BY purchased_date DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying items: %v", err)
		return nil, fmt.Errorf("failed to fetch shopping items")
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.CategoryID,

			&item.ItemName,
			&item.Quantity,
			&item.UnitPrice,
		)
		if err != nil {
			log.Printf("Error scanning item row: %v", err)
			continue
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error during row iteration: %w", rows.Err())
	}

	return items, nil
}

// --- FUNGSI BARU YANG DIMINTA ---

// GetTotalSpendingByDateRange menghitung total pengeluaran user dalam rentang waktu
// Fungsi ini dipanggil oleh DashboardHandler

func (r *ItemRepository) GetTotalSpendingByDateRange(ctx context.Context, userID int, startDate time.Time, endDate time.Time) (float64, error) {
	// COALESCE digunakan untuk memastikan 0 dikembalikan jika tidak ada data (SUM = NULL)
	query := `SELECT COALESCE(SUM(total_harga), 0) 
	          FROM items 
	          WHERE id_user = $1 AND purchased_date BETWEEN $2 AND $3`

	var totalSpending float64

	// Format tanggal ke string YYYY-MM-DD untuk query SQL
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	err := r.db.QueryRowContext(ctx, query, userID, startDateStr, endDateStr).Scan(&totalSpending)
	if err != nil {
		// ErrNoRows tidak akan terjadi karena COALESCE, tapi kita tangani error lain
		log.Printf("Error calculating total spending for user %d: %v", userID, err)
		return 0, fmt.Errorf("failed to calculate total spending: %w", err)
	}

	return totalSpending, nil
}
func (r *ItemRepository) UpdateItem(ctx context.Context, item *model.Item) error {
	query := `
		UPDATE items
		SET id_kategori = $1, nama_item = $2, jumlah_item = $3, harga_satuan = $4, total_harga = $5
		WHERE id_item = $6 AND id_user = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		item.CategoryID,
		item.ItemName,
		item.Quantity,
		item.UnitPrice,
		item.TotalCost,
		item.ID,
		item.UserID,
	)
	if err != nil {
		log.Printf("Error updating item: %v", err)
		return fmt.Errorf("failed to update shopping item: %w", err)
	}
	return nil

}

func (r *ItemRepository) DeleteItem(ctx context.Context, itemID int, userID int) error {
	query := `
		DELETE FROM items
		WHERE id_item = $1 AND id_user = $2
	`
	_, err := r.db.ExecContext(ctx, query, itemID, userID)
	if err != nil {
		log.Printf("Error deleting item: %v", err)
		return fmt.Errorf("failed to delete shopping item: %w", err)
	}
	return nil

}

// CategoryExists memeriksa apakah kategori dengan ID tertentu ada milik user tertentu
func (r *ItemRepository) CategoryExists(ctx context.Context, categoryID int, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM referensi_kategori WHERE id_kategori = $1 AND id_user = $2`
	var count int
	err := r.db.QueryRowContext(ctx, query, categoryID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Note: Repository untuk Budget, Category, dan Report akan dibuat di tahap selanjutnya
// karena fokus awal adalah pada fitur dasar (Login/Register/Tambah Item) dan Rework Laporan.
