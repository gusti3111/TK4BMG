package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Pustaka JWT
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Kunci rahasia untuk JWT.
// PENTING: Kunci ini HARUS SAMA PERSIS dengan kunci yang Anda gunakan di `auth_middleware.go`!
// (Sebaiknya simpan ini di file .env dan impor dari config)
var jwtSecretKey = []byte("kunci_rahasia_anda_yang_sangat_aman")

// AuthHandler menangani logika autentikasi.
type AuthHandler struct {
	userRepo *repository.UserRepository
}

// NewAuthHandler membuat instance AuthHandler baru.
func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// Register
// Endpoint: POST /api/v1/register
func (h *AuthHandler) Register(c *gin.Context) {
	// 1. Bind JSON langsung ke Model Entitas
	var req model.User // Menggunakan model.User (sesuai ERD TK2)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid", "details": err.Error()})
		return
	}

	// 2. Validasi input dasar (Menggunakan 'Nama' sesuai ERD TK2)
	// (Pastikan struct model.User Anda menggunakan json:"nama", json:"email", dll)
	if req.Username == "" || req.Password == "" || req.Name == "" || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field (username, password, nama, email) wajib diisi"})
		return
	}

	// 3. Set Role Default (Logika Bisnis di Handler)
	role := "user"

	// 4. Panggil Repository (PERBAIKAN ARSITEKTUR)
	// Buat request sesuai kebutuhan repository
	registerReq := &model.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
		Email:    req.Email,
	}

	err := h.userRepo.CreateUser(c.Request.Context(), registerReq, role)
	if err != nil {
		log.Printf("Error memanggil CreateUser: %v", err)
		// TODO: Cek error spesifik (misal: "duplicate key" atau "unique constraint")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user", "details": err.Error()})
		return
	}

	// 5. Kirim Respons Sukses (JSON Valid untuk frontend)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil! Silakan login.",
		"user_id": req.ID, // ID ini didapat dari 'RETURNING id_user' di repo CreateUser
	})
}

// Login
// Endpoint: POST /api/v1/login
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
	user, err := h.userRepo.GetUserByUsername(c.Request.Context(), req.Username)
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
