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
	// 1. Koneksi Database
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Kesalahan Fatal saat koneksi DB: %v", err)
	}
	defer db.CloseDB()

	// 2. Setup Router Gin
	r := gin.Default()

	// 3. Setup Routes
	setupRoutes(r)

	// 4. Jalankan Server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("✅ Backend BMG berjalan di http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// --- Inisialisasi Repository ---
	itemRepo := repository.NewItemRepository()
	categoryRepo := repository.NewCategoryRepository()
	budgetRepo := repository.NewBudgetRepository()
	reportRepo := repository.NewReportRepository()

	// --- Inisialisasi Handler ---
	authHandler := handler.NewAuthHandler()
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	itemHandler := handler.NewItemHandler(itemRepo, categoryRepo)
	dashHandler := handler.NewDashboardHandler(itemRepo, budgetRepo, reportRepo)

	// Terapkan CORS untuk semua endpoint
	r.Use(middleware.CORSMiddleware())

	// --- HEALTH CHECK ---
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "BMG Backend API"})
	})

	// --- RUTE PUBLIK ---
	publicV1 := r.Group("/api/v1")
	{
		publicV1.POST("/register", authHandler.Register)
		publicV1.POST("/login", authHandler.Login)
	}

	// --- RUTE TERLINDUNGI (PERLU TOKEN) ---
	secureV1 := r.Group("/api/v1")
	secureV1.Use(middleware.AuthMiddleware())
	{
		// Dashboard
		secureV1.GET("/dashboard/summary", dashHandler.GetDashboardSummary)
		secureV1.GET("/dashboard/charts", dashHandler.GetDashboardCharts)

		// Items
		secureV1.POST("/items", itemHandler.CreateItem)
		secureV1.GET("/items/:id", itemHandler.GetItemsByUserID)
		// secureV1.PUT("/items/:id", itemHandler.UpdateItem)
		// secureV1.DELETE("/items/:id", itemHandler.DeleteItem)

		// Kategori
		secureV1.POST("/kategori", categoryHandler.CreateCategory)
		secureV1.GET("/kategori", categoryHandler.GetCategories)
		secureV1.PUT("/kategori/:id", categoryHandler.UpdateCategory)
		secureV1.DELETE("/kategori/:id", categoryHandler.DeleteCategory)

		// Budget
		secureV1.POST("/budgets", handleSetBudget)
	}
}

// Placeholder Handler Budget
func handleSetBudget(c *gin.Context) {
	userID := c.GetInt("user_id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Set Budget logic placeholder",
		"user_id": userID,
	})
}
