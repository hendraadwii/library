package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// Config holds the JWT configuration
type Config struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

// TokenPair contains an access token and a refresh token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // Expiry in seconds
}

// TokenManager handles JWT token operations
type TokenManager struct {
	Config Config
}

// NewTokenManager creates a new TokenManager
func NewTokenManager(config Config) *TokenManager {
	// Set default values if not provided
	if config.AccessTokenExpiry == 0 {
		config.AccessTokenExpiry = 15 * time.Minute
	}
	if config.RefreshTokenExpiry == 0 {
		config.RefreshTokenExpiry = 7 * 24 * time.Hour // 7 days
	}
	return &TokenManager{Config: config}
}

// GenerateTokenPair generates a new access and refresh token pair
func (tm *TokenManager) GenerateTokenPair(userID int64, email, role string) (*TokenPair, error) {
	// Create access token
	accessToken, err := tm.generateToken(userID, email, role, "access", tm.Config.AccessTokenSecret, tm.Config.AccessTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// Create refresh token
	refreshToken, err := tm.generateToken(userID, email, role, "refresh", tm.Config.RefreshTokenSecret, tm.Config.RefreshTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(tm.Config.AccessTokenExpiry.Seconds()),
	}, nil
}

// generateToken creates a new JWT token
func (tm *TokenManager) generateToken(userID int64, email, role, tokenType, secret string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    tm.Config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyAccessToken validates an access token and returns the claims
func (tm *TokenManager) VerifyAccessToken(tokenString string) (*Claims, error) {
	return tm.verifyToken(tokenString, tm.Config.AccessTokenSecret, "access")
}

// VerifyRefreshToken validates a refresh token and returns the claims
func (tm *TokenManager) VerifyRefreshToken(tokenString string) (*Claims, error) {
	return tm.verifyToken(tokenString, tm.Config.RefreshTokenSecret, "refresh")
}

// verifyToken validates a token and returns the claims
func (tm *TokenManager) verifyToken(tokenString, secret, expectedType string) (*Claims, error) {
	claims := &Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Verify token type
	if claims.TokenType != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.TokenType)
	}

	return claims, nil
}

// RefreshTokens generates new token pair from a valid refresh token
func (tm *TokenManager) RefreshTokens(refreshToken string) (*TokenPair, error) {
	// Verify the refresh token
	claims, err := tm.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return tm.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
} 