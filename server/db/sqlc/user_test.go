package db

import (
	"context"
	"testing"
	"time"

	"github.com/bysergr/simple-bank/utils"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(12))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	randomUser := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), randomUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomUser.Username, user.Username)
	require.Equal(t, randomUser.Email, user.Email)
	require.Equal(t, randomUser.HashedPassword, user.HashedPassword)
	require.Equal(t, randomUser.FullName, user.FullName)

	require.WithinDuration(t, randomUser.CreatedAt, user.CreatedAt, time.Second)
}
