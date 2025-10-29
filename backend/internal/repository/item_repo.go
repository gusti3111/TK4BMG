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
	// TotalCost dihitung di service/handler sebelum dipanggil
	query := `INSERT INTO items (id_kategori, id_user, nama_item, jumlah_item, harga_satuan, total_harga, purchased_date)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		item.CategoryID,
		item.UserID,
		item.ItemName,
		item.Quantity,
		item.UnitPrice,
		item.TotalCost,
		time.Now().Format("2006-01-02"), // Format tanggal saat ini
	)

	if err != nil {
		log.Printf("Error inserting item: %v", err)
		return fmt.Errorf("failed to save shopping item")
	}
	return nil
}

// GetItemsByUserID fetches all shopping items for a specific user within a timeframe (simple version)
func (r *ItemRepository) GetItemsByUserID(ctx context.Context, userID int) ([]model.Item, error) {
	// Query ini bisa dioptimalkan dengan filter tanggal di masa depan (TK4 Rework)
	query := `SELECT id_item, id_kategori, id_user, nama_item, jumlah_item, harga_satuan, total_harga, purchased_date FROM items WHERE id_user = $1 ORDER BY purchased_date DESC`

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
			&item.UserID,
			&item.ItemName,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalCost,
			&item.PurchasedDate,
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

// Note: Repository untuk Budget, Category, dan Report akan dibuat di tahap selanjutnya
// karena fokus awal adalah pada fitur dasar (Login/Register/Tambah Item) dan Rework Laporan.
