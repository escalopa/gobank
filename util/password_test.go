package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestGeneratePassword(t *testing.T) {
	password := RandomString(10)

	hashedPassword, err := GenerateHashPassword(password)
	require.NoError(t, err)

	require.NotEqual(t, password, hashedPassword)
	require.NotEmpty(t, hashedPassword)

	err = CheckHashedPassword(hashedPassword, password)
	require.NoError(t, err)

	password = RandomString(32)
	err = CheckHashedPassword(hashedPassword, password)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
