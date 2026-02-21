package util

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	suite.Suite
}

func TestPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}

func (s *PasswordTestSuite) TestHashAndCheckPassword() {
	password := "my_secure_password"
	hashed, err := HashPassword(password)
	s.Require().NoError(err, "Hashing password should not return an error")
	s.Require().NotEmpty(hashed, "Hashed password should not be empty")

	err = CheckPasswordHash(password, hashed)
	s.Require().NoError(err, "Checking correct password should not return an error")

	err = CheckPasswordHash("wrong_password", hashed)
	s.Require().Error(err, "Checking incorrect password should return an error")
}
