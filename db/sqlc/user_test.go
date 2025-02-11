package db

import (
	"context"
	"testing"
	"time"

	"github.com/arturbaldoramos/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err, "failed to create account")
	require.NotEmpty(t, user, "account should not be empty")

	require.Equal(t, arg.Username, user.Username, "username should be the same")
	require.Equal(t, arg.HashedPassword, user.HashedPassword, "hashed password should be the same")
	require.Equal(t, arg.FullName, user.FullName, "full name should be the same")
	require.Equal(t, arg.Email, user.Email, "email should be the same")

	require.True(t, user.PasswordChangedAt.IsZero(), "id should not be zero")
	require.NotZero(t, user.CreatedAt, "created at should not be zero")

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err, "failed to get account")
	require.NotEmpty(t, user2, "account should not be empty")

	require.Equal(t, user1.Username, user2.Username, "username should be the same")
	require.Equal(t, user1.HashedPassword, user2.HashedPassword, "hashed password should be the same")
	require.Equal(t, user1.FullName, user2.FullName, "full name should be the same")
	require.Equal(t, user1.Email, user2.Email, "email should be the same")

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second, "password changed at should be the same")
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second, "created at should be the same")
}
