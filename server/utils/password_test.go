package utils

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password := RandomString(12)

	hashedPasswordOne, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPasswordOne)
	require.NotEqual(t, password, hashedPasswordOne)

	err = CheckPassword(password, hashedPasswordOne)
	require.NoError(t, err)

	wrongPassword := RandomString(12)
	err = CheckPassword(wrongPassword, hashedPasswordOne)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPasswordTwo, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPasswordTwo)
	require.NotEqual(t, password, hashedPasswordOne, hashedPasswordTwo)
}
