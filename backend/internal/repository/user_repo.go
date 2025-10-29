package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gusti3111/TKBMG/backend/internal/db"
	"github.com/gusti3111/TKBMG/backend/internal/model"
)

// UserRepository handles database operations related to User entity
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new repository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{db: db.DB}
}

// GetUserByUsername fetches a user by their username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id_user, username, password, nama, email, role FROM "User" WHERE username = $1`
	user := new(model.User)

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		log.Printf("Error querying user by username: %v", err)
		return nil, fmt.Errorf("database query error")
	}
	return user, nil
}

// CreateUser saves a new user to the database
func (r *UserRepository) CreateUser(ctx context.Context, req *model.RegisterRequest, hashedPassword string) error {
	query := `INSERT INTO "User" (username, password, nama, email, role) VALUES ($1, $2, $3, $4, 'member')`

	_, err := r.db.ExecContext(ctx, query, req.Username, hashedPassword, req.Name, req.Email)
	if err != nil {
		// Specific error handling for UNIQUE constraint violation (e.g., username/email already exists)
		// This requires more complex error checking depending on the DB driver, but for simplicity:
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("failed to create user")
	}
	return nil
}
