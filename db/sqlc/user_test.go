package db

import (
	"context"
	"testing"
	"time"

	"github.com/escalopa/go-bank/util"
	"github.com/stretchr/testify/require"
)

func validateUserBasic(t *testing.T, user User) {
	require.NotEmpty(t, user.Username)
	require.NotEmpty(t, user.HashedPassword)
	require.NotEmpty(t, user.FullName)
	require.NotEmpty(t, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}

func createRandomUser(t *testing.T) User {
	hashPassword, err := util.GenerateHashPassword(util.RandomString(10))
	require.NoError(t, err)

	user1 := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user2, err := testQueries.CreateUser(context.Background(), user1)
	require.NoError(t, err)

	// Check account values
	validateUserBasic(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	return user2
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)

	validateUserBasic(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
