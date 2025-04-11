package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestValidPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = ValidatePassword(password, hashedPassword)
	require.NoError(t, err)
}

func TestNotValidPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	wrongPassword := RandomString(6)
	err = ValidatePassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
