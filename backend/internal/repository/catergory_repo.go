package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gusti3111/TKBMG/backend/internal/db"
	"github.com/gusti3111/TKBMG/backend/internal/model"
)

// CategoryRepository menangani operasi database untuk 'referensi_kategori'.
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository membuat instance CategoryRepository baru.
func NewCategoryRepository() *CategoryRepository {
	// db.DB adalah koneksi global yang diinisialisasi di paket db
	return &CategoryRepository{db: db.DB}
}

// CreateKategori menambahkan kategori baru ke database untuk user tertentu.
// Ini dipanggil oleh halaman 'Referensi Belanja'.
func (r *CategoryRepository) CreateKategori(ctx context.Context, kategori *model.Category) error {
	query := `INSERT INTO referensi_kategori ( nama_kategori)
	          VALUES ($1) RETURNING id_kategori`

	// Asumsi Anda telah menambahkan kolom created_at di DB:
	// query := `INSERT INTO referensi_kategori (id_user, nama_kategori, created_at)
	// 	          VALUES ($1, $2, NOW()) RETURNING id_kategori`

	// Pastikan untuk mengambil UserID dari token di Handler dan mengisinya ke struct 'kategori'
	err := r.db.QueryRowContext(ctx, query, kategori.UserID, kategori.CategoryName).Scan(&kategori.ID)
	if err != nil {
		log.Printf("Error creating kategori: %v", err)
		return fmt.Errorf("failed to save kategori: %w", err)
	}
	return nil
}

// GetKategoriByUserID mengambil semua kategori milik user tertentu.
// Ini dipanggil oleh halaman 'DaftarBelanja' (untuk dropdown) dan 'ReferensiBelanja'.
func (r *CategoryRepository) GetKategoriByUserID(ctx context.Context, userID int) ([]model.Category, error) {
	query := `SELECT id_kategori, id_user, nama_kategori FROM referensi_kategori WHERE id_user = $1 ORDER BY nama_kategori ASC`

	// Sesuaikan query jika Anda menambahkan created_at
	// query := `SELECT id_kategori, id_user, nama_kategori, created_at FROM referensi_kategori WHERE id_user = $1 ORDER BY nama_kategori ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying kategori for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to fetch kategori: %w", err)
	}
	defer rows.Close()

	var kategoriList []model.Category
	for rows.Next() {
		var k model.Category
		// Sesuaikan Scan jika Anda menambahkan created_at
		// if err := rows.Scan(&k.ID, &k.UserID, &k.NamaKategori, &k.CreatedAt); err != nil { ... }
		if err := rows.Scan(&k.ID, &k.UserID, &k.CategoryName); err != nil {
			log.Printf("Error scanning kategori row: %v", err)
			continue
		}
		kategoriList = append(kategoriList, k)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error during row iteration: %w", rows.Err())
	}

	return kategoriList, nil
}

// UpdateKategori memperbarui nama kategori milik user tertentu.
func (r *CategoryRepository) UpdateKategori(ctx context.Context, kategori *model.Category) error {
	query := `UPDATE referensi_kategori SET nama_kategori = $1 WHERE id_kategori = $2 AND id_user = $3`

	result, err := r.db.ExecContext(ctx, query, kategori.CategoryName, kategori.ID, kategori.UserID)
	if err != nil {
		log.Printf("Error updating kategori: %v", err)
		return fmt.Errorf("failed to update kategori: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("kategori not found or user not authorized")
	}

	return nil
}

// DeleteKategori menghapus kategori milik user tertentu.
func (r *CategoryRepository) DeleteKategori(ctx context.Context, kategoriID int, userID int) error {
	// PERHATIAN: Query ini akan gagal jika 'id_kategori' masih digunakan di tabel 'items' (Constraint Error).
	// Di 'Handler', Anda harus memeriksa dulu apakah kategori ini digunakan,
	// atau set 'ON DELETE SET NULL' pada Foreign Key 'items.id_kategori'.
	query := `DELETE FROM referensi_kategori WHERE id_kategori = $1 AND id_user = $2`

	result, err := r.db.ExecContext(ctx, query, kategoriID, userID)
	if err != nil {
		log.Printf("Error deleting kategori: %v", err)
		return fmt.Errorf("failed to delete kategori: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("kategori not found or user not authorized")
	}

	return nil
}
