package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "my_secure_password"
	hashed, err := HashPassword(password)

	assert.NoError(t, err, "Hashing password should not return an error")
	assert.NotEmpty(t, hashed, "Hashed password should not be empty")

	err = CheckPasswordHash(password, hashed)
	assert.NoError(t, err, "Checking correct password should not return an error")

	err = CheckPasswordHash("wrong_password", hashed)
	assert.Error(t, err, "Checking incorrect password should return an error")
}
