package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gusti3111/TKBMG/backend/internal/model"
	"github.com/gusti3111/TKBMG/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte("your-very-secret-key")

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repository.NewUserRepository()}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) RegisterUser(ctx context.Context, req *model.RegisterRequest) error {
	if existingUser, _ := s.userRepo.GetUserByUsername(ctx, req.Username); existingUser != nil {
		return fmt.Errorf("username sudah terdaftar")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fmt.Errorf("gagal memproses password")
	}

	return s.userRepo.CreateUser(ctx, req, hashedPassword)
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("terjadi kesalahan sistem")
	}
	if user == nil {
		return nil, fmt.Errorf("username atau password salah")
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("username atau password salah")
	}

	// === Ganti dengan JWT sungguhan ===
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token")
	}

	return &model.LoginResponse{Token: tokenString, Role: user.Role}, nil
}

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetUserByUsername(ctx, username)
}
