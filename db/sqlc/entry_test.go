package db

import (
	"context"
	"testing"
	"time"

	"github.com/arturbaldoramos/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err, "failed to create entry")
	require.NotEmpty(t, entry, "entry should not be empty")

	require.Equal(t, arg.AccountID, entry.AccountID, "account_id should be the same")
	require.Equal(t, arg.Amount, entry.Amount, "amount should be the same")

	require.NotZero(t, entry.ID, "entry id should not be zero")
	require.NotZero(t, entry.CreatedAt, "created at should not be zero")

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err, "failed to get entry")
	require.NotEmpty(t, entry2, "entry should not be empty")

	require.Equal(t, entry1.ID, entry2.ID, "id should be the same")
	require.Equal(t, entry1.AccountID, entry2.AccountID, "account_id should be the same")
	require.Equal(t, entry1.Amount, entry2.Amount, "amount should be the same")

	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, time.Second, "created at should be the same")
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err, "failed to list entries")
	require.Len(t, entries, 5, "length of entries should be 5")

	for _, entry := range entries {
		require.NotEmpty(t, entry, "entry should not be empty")
	}
}
