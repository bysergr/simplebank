package db

import (
	"context"
	"database/sql"
	"github.com/bysergr/simple-bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, accountID int64) Entry {
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	createRandomEntry(t, randomAccount.ID)
}

func TestGetEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	randomEntry := createRandomEntry(t, randomAccount.ID)

	entry, err := testQueries.GetEntry(context.Background(), randomEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randomEntry.ID, entry.ID)
	require.Equal(t, randomAccount.ID, entry.AccountID)
	require.Equal(t, randomEntry.Amount, entry.Amount)
	require.WithinDuration(t, randomEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	randomEntry := createRandomEntry(t, randomAccount.ID)

	updateEntryParams := UpdateEntryParams{
		ID:     randomEntry.ID,
		Amount: utils.RandomMoney(),
	}

	entry, err := testQueries.UpdateEntry(context.Background(), updateEntryParams)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randomEntry.ID, entry.ID)
	require.Equal(t, randomEntry.AccountID, entry.AccountID)
	require.Equal(t, updateEntryParams.Amount, entry.Amount)
	require.WithinDuration(t, randomEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	randomEntry := createRandomEntry(t, randomAccount.ID)

	err := testQueries.DeleteEntry(context.Background(), randomEntry.ID)
	require.NoError(t, err)

	databaseEntry, err := testQueries.GetEntry(context.Background(), randomEntry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, databaseEntry)
}

func TestListEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, randomAccount.ID)
	}

	listEntriesParams := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), listEntriesParams)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
