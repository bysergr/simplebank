package db

import (
	"context"
	"database/sql"
	"github.com/bysergr/simple-bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccountID, toAccountID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	firstRandomAccount := createRandomAccount(t)
	secondRandomAccount := createRandomAccount(t)

	createRandomTransfer(t, firstRandomAccount.ID, secondRandomAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	firstRandomAccount := createRandomAccount(t)
	secondRandomAccount := createRandomAccount(t)
	randomTransfer := createRandomTransfer(t, firstRandomAccount.ID, secondRandomAccount.ID)

	transfer, err := testQueries.GetTransfer(context.Background(), randomTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randomTransfer.ID, transfer.ID)
	require.Equal(t, firstRandomAccount.ID, transfer.FromAccountID)
	require.Equal(t, secondRandomAccount.ID, transfer.ToAccountID)
	require.Equal(t, randomTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, randomTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	firstRandomAccount := createRandomAccount(t)
	secondRandomAccount := createRandomAccount(t)
	randomTransfer := createRandomTransfer(t, firstRandomAccount.ID, secondRandomAccount.ID)

	updateTransferParams := UpdateTransferParams{
		ID:     randomTransfer.ID,
		Amount: utils.RandomMoney(),
	}

	transfer, err := testQueries.UpdateTransfer(context.Background(), updateTransferParams)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randomTransfer.ID, transfer.ID)
	require.Equal(t, randomTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, randomTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, updateTransferParams.Amount, transfer.Amount)
	require.WithinDuration(t, randomTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestDeleteTrans(t *testing.T) {
	firstRandomAccount := createRandomAccount(t)
	secondRandomAccount := createRandomAccount(t)
	randomTransfer := createRandomTransfer(t, firstRandomAccount.ID, secondRandomAccount.ID)

	err := testQueries.DeleteTransfer(context.Background(), randomTransfer.ID)
	require.NoError(t, err)

	databaseTransfer, err := testQueries.GetTransfer(context.Background(), randomTransfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, databaseTransfer)
}

func TestListTransfers(t *testing.T) {
	firstRandomAccount := createRandomAccount(t)
	secondRandomAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, firstRandomAccount.ID, secondRandomAccount.ID)
	}

	listTransfersParams := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), listTransfersParams)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, entry := range transfers {
		require.NotEmpty(t, entry)
	}
}
