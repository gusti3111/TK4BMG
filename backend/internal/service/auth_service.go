package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gusti3111/TKBMG/backend/internal/config" // Import config terpusat
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService menangani logika bisnis terkait otentikasi
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService adalah constructor untuk AuthService.
// INI FUNGSI YANG HILANG.
func NewAuthService() *AuthService {
	// Kita asumsikan NewUserRepository() ada di paket repository Anda
	return &AuthService{userRepo: repository.NewUserRepository()}
}

// === FUNGSI HELPER PASSWORD (HILANG) ===

// HashPassword membuat hash bcrypt dari password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash membandingkan plain password dengan hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// === FUNGSI UTAMA SERVICE ===

// RegisterUser memvalidasi dan membuat user baru
func (s *AuthService) RegisterUser(ctx context.Context, req *model.RegisterRequest) error {
	if existingUser, _ := s.userRepo.GetUserByUsername(ctx, req.Username); existingUser != nil {
		return fmt.Errorf("username sudah terdaftar")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fmt.Errorf("gagal memproses password")
	}

	// Kita asumsikan CreateUser ada di repo Anda
	return s.userRepo.CreateUser(ctx, req, hashedPassword)
}

// Login memvalidasi kredensial dan mengembalikan token JWT
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// Asumsi GetUserByUsername mengembalikan error jika user tidak ada
		log.Printf("Error getting user: %v", err)
		return nil, fmt.Errorf("username atau password salah")
	}
	if user == nil {
		// Fallback jika repo mengembalikan nil, nil
		return nil, fmt.Errorf("username atau password salah")
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("username atau password salah")
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	// Gunakan secret key dari config
	tokenString, err := token.SignedString(config.JWTSecretKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return nil, fmt.Errorf("gagal membuat token")
	}

	return &model.LoginResponse{Token: tokenString, Role: user.Role}, nil
}

// GetUserByUsername mengambil data user (tanpa password)
func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	// Pastikan fungsi repo Anda tidak mengembalikan password
	return s.userRepo.GetUserByUsername(ctx, username)
}
