package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Cell6969/go_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	ctx := context.Background()
	arg := CreateAccountParams{
		Owner:    util.GenerateRandomName(),
		Balance:  util.GenerateRandomMoney(),
		Currency: util.GenerateRandomCurrency(),
	}

	account, err := testQueries.CreateAccount(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)

	return account
}

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()
	testQueries.ResetAccountTable(ctx)
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	testQueries.ResetAccountTable(ctx)

	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(ctx, account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	ctx := context.Background()
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.GenerateRandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(ctx, arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, account1.ID, updatedAccount.ID)
	require.Equal(t, account1.Owner, updatedAccount.Owner)
	require.NotEqual(t, account1.Balance, updatedAccount.Balance)
	require.Equal(t, account1.Currency, updatedAccount.Currency)
	require.WithinDuration(t, account1.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()
	testQueries.ResetAccountTable(ctx)

	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(ctx, account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(ctx, account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(ctx, arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
