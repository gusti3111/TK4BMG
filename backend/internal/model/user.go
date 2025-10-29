package model

// User represents the data structure for the "User" entity in the database (TK2 ERD)
type User struct {
	ID       int    `json:"id_user"`
	Username string `json:"username"`
	Password string `json:"-"` // Omit password when returning JSON
	Name     string `json:"nama"`
	Email    string `json:"email"`
	Role     string `json:"role"` // "member" or "admin"
}

// LoginRequest defines the structure for incoming login data
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest defines the structure for incoming registration data
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// LoginResponse defines the data structure returned upon successful login
type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}
