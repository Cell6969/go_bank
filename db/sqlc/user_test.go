package db

import (
	"context"
	"testing"
	"time"

	"github.com/Cell6969/go_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	ctx := context.Background()
	arg := CreateUserParams{
		Username: util.GenerateRandomName(),
		Password: "secret",
		FullName: util.GenerateRandomName(),
		Email:    util.GenerateRandomEmail(),
	}

	user, err := testQueries.CreateUser(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	testQueries.ResetUserTable(ctx)
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	testQueries.ResetUserTable(ctx)

	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(ctx, user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
