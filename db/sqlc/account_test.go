package db

import (
	"context"
	"testing"
	"time"

	"github.com/arturbaldoramos/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err, "failed to create account")
	require.NotEmpty(t, account, "account should not be empty")
	require.Equal(t, arg.Owner, account.Owner, "owner should be the same")
	require.Equal(t, arg.Balance, account.Balance, "balance should be the same")

	require.NotZero(t, account.ID, "account id should not be zero")
	require.NotZero(t, account.CreatedAt, "created at should not be zero")

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err, "failed to get account")
	require.NotEmpty(t, account2, "account should not be empty")

	require.Equal(t, account1.ID, account2.ID, "id should be the same")
	require.Equal(t, account1.Owner, account2.Owner, "owner should be the same")
	require.Equal(t, account1.Balance, account2.Balance, "balance should be the same")
	require.Equal(t, account1.Currency, account2.Currency, "currency should be the same")

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second, "created at should be the same")
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err, "failed to update account")
	require.NotEmpty(t, account2, "account should not be empty")

	require.Equal(t, account1.ID, account2.ID, "id should be the same")
	require.Equal(t, account1.Owner, account2.Owner, "owner should be the same")
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency, "currency should be the same")

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second, "created at should be the same")
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err, "failed to delete account")

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err, "account should be deleted")
	require.Empty(t, account2, "account should be empty")
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err, "failed to list accounts")
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account, "account should not be empty")
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
