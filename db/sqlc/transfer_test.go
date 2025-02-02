package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err, "failed to create transfer")
	require.NotEmpty(t, transfer, "transfer should not be empty")

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID, "from account id should be the same")
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID, "to account id should be the same")
	require.Equal(t, arg.Amount, transfer.Amount, "amount should be the same")

	require.NotZero(t, transfer.ID, "transfer id should not be zero")
	require.NotZero(t, transfer.CreatedAt, "created at should not be zero")

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err, "failed to get transfer")
	require.NotEmpty(t, transfer2, "transfer should not be empty")

	require.Equal(t, transfer1.ID, transfer2.ID, "id should be the same")
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID, "from account id should be the same")
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID, "to account id should be the same")
	require.Equal(t, transfer1.Amount, transfer2.Amount, "amount should be the same")

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second, "created at should be the same")
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
