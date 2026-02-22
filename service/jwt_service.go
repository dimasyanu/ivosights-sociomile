package service

import (
	"fmt"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`

	jwt.RegisteredClaims
}

type JwtService struct {
	config *config.JwtConfig
}

func NewJwtService(c *config.JwtConfig) *JwtService {
	return &JwtService{
		config: c,
	}
}

func (s *JwtService) GenerateJWT(email string, duration time.Duration) (string, error) {
	jwtKey := []byte(s.config.SecretKey)

	// Set the expiration time for the token
	expirationTime := time.Now().Add(duration)

	// Create the claims (payload)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the exp (expiration time) claim is a NumericDate.
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "localhost", // The issuer of the token
			Subject:   email,       // The subject of the token (e.g., user ID)
		},
	}

	// Declare the token with the signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key to get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *JwtService) ValidateJWT(tokenString string) (string, error) {
	jwtKey := []byte(s.config.SecretKey)

	claims := &Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and not something else
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate the token and claims
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Email, nil
}

func (s *JwtService) RefreshJWT(tokenString string) (string, error) {
	email, err := s.ValidateJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	return s.GenerateJWT(email, time.Hour*24)
}
