package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	StatusActive UserStatus = "active"
	StatusLocked UserStatus = "locked"
)

// User represents a user in the library system
type User struct {
	ID           int64      `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`                  // Never expose in JSON
	Password     string     `json:"password,omitempty"` // Only used for input
	FullName     string     `json:"full_name"`
	Role         UserRole   `json:"role"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// UserService handles user-related database operations
type UserService struct {
	DB *sql.DB
}

// NewUserService creates a new UserService
func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// HashPassword generates a bcrypt hash from a password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Create adds a new user to the database
func (s *UserService) Create(user *User) error {
	// Hash the password
	var err error
	user.PasswordHash, err = HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Clear the plain text password
	user.Password = ""

	// Set default status to active
	if user.Status == "" {
		user.Status = StatusActive
	}

	query := `
		INSERT INTO users 
		(email, password_hash, full_name, role, status) 
		VALUES (?, ?, ?, ?, ?)
	`

	// Use parameterized query to prevent SQL injection
	result, err := s.DB.Exec(
		query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.Status,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	user.ID = id
	return nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int64) (*User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, status, created_at, updated_at 
		FROM users WHERE id = ?
	`

	var user User
	err := s.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with ID %d", id)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, status, created_at, updated_at 
		FROM users WHERE email = ?
	`

	var user User
	err := s.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with email %s", email)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

// Update updates an existing user
func (s *UserService) Update(user *User) error {
	// Check if password needs to be updated
	if user.Password != "" {
		var err error
		user.PasswordHash, err = HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = "" // Clear plaintext password
	}

	query := `
		UPDATE users 
		SET full_name = ?, role = ?, status = ?
	`
	params := []interface{}{
		user.FullName,
		user.Role,
		user.Status,
	}

	// Add password hash to update if it was changed
	if user.PasswordHash != "" {
		query += ", password_hash = ?"
		params = append(params, user.PasswordHash)
	}

	// Add email to update if it was changed
	if user.Email != "" {
		query += ", email = ?"
		params = append(params, user.Email)
	}

	query += " WHERE id = ?"
	params = append(params, user.ID)

	result, err := s.DB.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID %d", user.ID)
	}

	return nil
}

// Delete removes a user by ID
func (s *UserService) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found with ID %d", id)
	}

	return nil
}

// List returns a list of users with pagination and filters
func (s *UserService) List(limit, offset int, search string, role UserRole, status UserStatus) ([]User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, email, password_hash, full_name, role, status, created_at, updated_at 
		FROM users 
		WHERE 1=1
	`
	params := []interface{}{}

	// Add search filter
	if search != "" {
		query += " AND (full_name LIKE ? OR email LIKE ?)"
		searchParam := "%" + search + "%"
		params = append(params, searchParam, searchParam)
	}

	// Add role filter
	if role != "" {
		query += " AND role = ?"
		params = append(params, role)
	}

	// Add status filter
	if status != "" {
		query += " AND status = ?"
		params = append(params, status)
	}

	// Add order and pagination
	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	params = append(params, limit, offset)

	rows, err := s.DB.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.Role,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// Count returns the total number of users with filters
func (s *UserService) Count(search string, role UserRole, status UserStatus) (int, error) {
	query := "SELECT COUNT(*) FROM users WHERE 1=1"
	params := []interface{}{}

	// Add search filter
	if search != "" {
		query += " AND (full_name LIKE ? OR email LIKE ?)"
		searchParam := "%" + search + "%"
		params = append(params, searchParam, searchParam)
	}

	// Add role filter
	if role != "" {
		query += " AND role = ?"
		params = append(params, role)
	}

	// Add status filter
	if status != "" {
		query += " AND status = ?"
		params = append(params, status)
	}

	var count int
	err := s.DB.QueryRow(query, params...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

// Authenticate verifies a user's credentials and returns the user if valid
func (s *UserService) Authenticate(email, password string) (*User, error) {
	user, err := s.GetByEmail(email)
	if err != nil {
		// Don't reveal that the email doesn't exist
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is locked
	if user.Status == StatusLocked {
		return nil, fmt.Errorf("account is locked")
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
