package models

import (
	"database/sql"
	"time"
)

type PasswordResetPin struct {
	Email     string    `json:"email"`
	Pin       string    `json:"pin"`
	ExpiredAt time.Time `json:"expired_at"`
}

type PasswordResetService struct {
	DB *sql.DB
}

func NewPasswordResetService(db *sql.DB) *PasswordResetService {
	return &PasswordResetService{DB: db}
}

func (s *PasswordResetService) Create(pin *PasswordResetPin) error {
	_, err := s.DB.Exec(`INSERT INTO password_reset_pins (email, pin, expired_at) VALUES (?, ?, ?)`, pin.Email, pin.Pin, pin.ExpiredAt)
	return err
}

func (s *PasswordResetService) GetByEmail(email string) (*PasswordResetPin, error) {
	var pin PasswordResetPin
	err := s.DB.QueryRow(`SELECT email, pin, expired_at FROM password_reset_pins WHERE email = ?`, email).Scan(&pin.Email, &pin.Pin, &pin.ExpiredAt)
	if err != nil {
		return nil, err
	}
	return &pin, nil
}

func (s *PasswordResetService) DeleteByEmail(email string) error {
	_, err := s.DB.Exec(`DELETE FROM password_reset_pins WHERE email = ?`, email)
	return err
}
