package model

import (
	"database/sql"
	"time"
)

// Item represents the data structure for the "Items" entity
type Item struct {
	ID            int            `json:"id_item"`
	UserID        int            `json:"id_user"`
	CategoryID    int            `json:"id_kategori"`
	ItemName      string         `json:"nama_item" binding:"required"`
	Quantity      int            `json:"jumlah_item" binding:"required"`
	UnitPrice     float64        `json:"harga_satuan" binding:"required"`
	TotalCost     float64        `json:"total_harga"` // Dihitung di backend (TK1 requirement)
	PurchasedDate time.Time      `json:"purchased_date"`
	CategoryName  sql.NullString `json:"nama_kategori,omitempty"` // For reporting joins
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
