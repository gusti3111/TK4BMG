package service

import (
	"context"
	"fmt"
	"log"

	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication and authorization business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new service instance
func NewAuthService() *AuthService {
	return &AuthService{userRepo: repository.NewUserRepository()}
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a raw password with its hashed version
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RegisterUser handles the registration business logic
func (s *AuthService) RegisterUser(ctx context.Context, req *model.RegisterRequest) error {
	// 1. Cek apakah user sudah ada
	if existingUser, _ := s.userRepo.GetUserByUsername(ctx, req.Username); existingUser != nil {
		return fmt.Errorf("username sudah terdaftar")
	}

	// 2. Hash password sebelum disimpan
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fmt.Errorf("gagal memproses password")
	}

	// 3. Simpan user ke database
	return s.userRepo.CreateUser(ctx, req, hashedPassword)
}

// Login handles the login business logic
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 1. Ambil user dari database
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("terjadi kesalahan sistem")
	}
	if user == nil {
		return nil, fmt.Errorf("username atau password salah") // General error for security
	}

	// 2. Bandingkan password
	if !CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("username atau password salah")
	}

	// 3. Generate Token (Placeholder JWT)
	// Implementasi JWT dihilangkan untuk kesederhanaan, diganti dengan placeholder
	token := fmt.Sprintf("jwt_token_for_user_%d", user.ID)

	return &model.LoginResponse{Token: token, Role: user.Role}, nil
}
