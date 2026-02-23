package service

import (
	"testing"
	"time"

	"github.com/dimasyanu/ivosights-sociomile/config"
	"github.com/stretchr/testify/suite"
)

const testEmail = "test@example.com"

type JwtServiceTestSuite struct {
	svc *JwtService

	suite.Suite
}

func TestJwtServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceTestSuite))
}

func (suite *JwtServiceTestSuite) SetupSuite() {
	cfg := &config.JwtConfig{
		SecretKey: "test_secret_key",
	}
	suite.svc = NewJwtService(cfg)
}

func (s *JwtServiceTestSuite) TestGenerateToken() {
	token, err := s.svc.GenerateJWT(testEmail, time.Hour*24)
	s.NoError(err)
	s.NotEmpty(token)
}

func (s *JwtServiceTestSuite) TestValidateToken() {
	token, err := s.svc.GenerateJWT(testEmail, time.Hour*24)
	s.NoError(err)

	validatedEmail, err := s.svc.ValidateJWT(token)
	s.NoError(err)
	s.Equal(testEmail, validatedEmail)
}

func (s *JwtServiceTestSuite) TestValidateToken_InvalidToken() {
	invalidToken := "invalid.token.value"

	_, err := s.svc.ValidateJWT(invalidToken)
	s.Error(err)
}

func (s *JwtServiceTestSuite) TestValidateToken_ExpiredToken() {
	// Create a token with a short expiration time
	token, err := s.svc.GenerateJWT(testEmail, time.Millisecond*1)
	s.NoError(err)

	// Wait for the token to expire
	time.Sleep(time.Millisecond * 10)

	// Validate the expired token
	_, err = s.svc.ValidateJWT(token)
	s.Error(err)
}
