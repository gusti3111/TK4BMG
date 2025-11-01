package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte("your-very-secret-key") // replace with a secure key or load from env

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if err := h.authService.RegisterUser(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil"})
}

// Login handles POST /v1/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 1. Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	// 2. Dapatkan User dari Repo (dari file user_repository.go Anda)
	user, err := h.authService.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		// User tidak ditemukan
		log.Printf("Login gagal (user not found): %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// 3. Bandingkan Password
	// (Diasumsikan repo GetUserByUsername mengembalikan HASH password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Password tidak cocok
		log.Printf("Login gagal (password mismatch) untuk user: %s", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// 4. Buat Token JWT
	token, err := createJWT(user)
	if err != nil {
		log.Printf("Error membuat token JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// 5. Kirim Token sebagai Respons
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   token,
		"user": gin.H{ // Kirim data user (tanpa password)
			"id":       user.ID,
			"username": user.Username,
			"nama":     user.Name, // (Menggunakan 'Nama' sesuai ERD TK2)
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// createJWT adalah fungsi helper untuk membuat token
func createJWT(user *model.User) (string, error) {
	// Buat claims (payload)
	claims := jwt.MapClaims{
		"sub":  user.ID, // Subject (User ID)
		"user": user.Username,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Expired dalam 24 jam
		"iat":  time.Now().Unix(),                     // Issued At
	}

	// Buat token dengan signing method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token dengan secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
