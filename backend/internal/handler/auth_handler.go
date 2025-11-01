package handler

import (
	"log"
	"net/http"

	// "time" // Tidak perlu lagi
	// "github.com/golang-jwt/jwt/v5" // Tidak perlu lagi
	"github.com/gin-gonic/gin"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/service"
	// "golang.org/x/crypto/bcrypt" // Tidak perlu lagi
)

// var jwtSecretKey = []byte("your-very-secret-key") // HAPUS: Ini penyebab bug

// AuthHandler holds the dependencies for authentication APIs
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new handler instance
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: service.NewAuthService()}
}

// Register handles POST /v1/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "details": err.Error()})
		return
	}

	// Delegasikan semua logika registrasi ke service
	if err := h.authService.RegisterUser(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil"})
}

// Login handles POST /v1/login
// (INI FUNGSI YANG DIPERBAIKI TOTAL)
func (h *AuthHandler) Login(c *gin.Context) {

	// 1. Gunakan model.LoginRequest, bukan struct anonim
	var req model.LoginRequest

	// 2. Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	// 3. Panggil SATU fungsi service Login
	// Service akan menangani (get user, check pass, create token)
	loginResponse, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		// Service akan mengembalikan error "username atau password salah"
		log.Printf("Login gagal untuk user: %s, error: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 4. Kirim Token dari service sebagai Respons
	// loginResponse sudah berisi Token dan Role
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   loginResponse.Token,
		"role":    loginResponse.Role,
	})
}

// HAPUS FUNGSI createJWT()
// func createJWT(user *model.User) (string, error) { ... }
// Fungsi ini tidak diperlukan lagi di handler, karena service sudah menanganinya.
