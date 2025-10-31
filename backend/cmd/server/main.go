package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// Import package internal
	"github.com/gusti3111/TKBMG/backend/internal/db"
	"github.com/gusti3111/TKBMG/backend/internal/handler"
	"github.com/gusti3111/TKBMG/backend/internal/middleware"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
)

func main() {
	// 1. Setup Koneksi Database
	// Ini akan menggunakan kredensial dari .env dan DB_HOST: db (nama service Docker)
	err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Kesalahan Fatal saat koneksi DB: %v", err)
	}
	defer db.CloseDB() // Pastikan koneksi ditutup saat aplikasi berhenti

	// 2. Setup Gin Router
	r := gin.Default()

	// 3. Definisikan API Routes
	setupRoutes(r)

	// 4. Konfigurasi dan Jalankan Server
	server := &http.Server{
		Addr:         ":8080", // Backend berjalan di port 8080
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Backend BMG berjalan di http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

// setupRoutes mendefinisikan semua API endpoint dan menerapkan middleware
func setupRoutes(r *gin.Engine) {
	// Inisiasi Handlers (Menggunakan Canvas yang sudah Anda buat)
	authHandler := handler.NewAuthHandler(
		repository.NewUserRepository(),
	)
	itemHandler := handler.NewItemHandler()

	// Pastikan NewItemRepository diambil dari package yang benar (misal: repository)
	// Jika sudah ada package repository, gunakan seperti berikut:
	// import "github.com/gusti3111/TKBMG/backend/internal/repository"
	// dashHandler := handler.NewDashboardHandler(
	// 	repository.NewItemRepository(),
	// )

	// Jika belum ada, buat fungsi NewItemRepository di handler package atau ganti dengan implementasi yang sesuai.
	// Untuk sementara, gunakan nil jika tidak diperlukan:
	dashHandler := handler.NewDashboardHandler(
		repository.NewItemRepository(),
		repository.NewBudgetRepository(),
		repository.NewReportRepository(),
	)

	// 1. Terapkan CORS ke Seluruh Router (untuk komunikasi Frontend)
	r.Use(middleware.CORSMiddleware())

	// Endpoint Dasar (Health Check)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "BMG Backend API"})
	})

	// --- V1 API Group: RUTE PUBLIK (Tidak Terlindungi) ---
	publicV1 := r.Group("/api/v1")
	{
		// Login dan Register tidak memerlukan token
		publicV1.POST("/register", authHandler.Register)
		publicV1.POST("/login", authHandler.Login)
	}

	// --- V1 API Group: RUTE TERLINDUNGI ---
	secureV1 := r.Group("/api/v1")
	// Terapkan AuthMiddleware() untuk rute yang memerlukan login (sesuai TK4)
	secureV1.Use(middleware.AuthMiddleware())
	{
		// 1. Item Belanja (Membutuhkan User ID dari Token)
		secureV1.POST("/items", itemHandler.CreateItem)

		// === TAMBAHKAN DUA RUTE INI ===
		secureV1.GET("/dashboard/summary", dashHandler.GetDashboardSummary)
		secureV1.GET("/dashboard/charts", dashHandler.GetDashboardCharts)

		// 2. Laporan (Area Prioritas Rework dari TK4 - membutuhkan optimasi)
		// secureV1.GET("/reports/weekly", h) // Handler placeholder

		// 3. Budget (Membutuhkan User ID dari Token)
		secureV1.POST("/budgets", handleSetBudget)
	}
}

// --- Handler Placeholder Tambahan ---

func handleSetBudget(c *gin.Context) {
	// Implementasi: Set Budget Mingguan (Membutuhkan User ID dari Token)
	userID := c.GetInt("user_id") // Mengambil User ID dari AuthMiddleware
	c.JSON(http.StatusOK, gin.H{"message": "Set Budget logic placeholder. User ID:", "user_id": userID})
}
