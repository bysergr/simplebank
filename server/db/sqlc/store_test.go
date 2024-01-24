package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStoreSQl(testDB)

	firstAccount := createRandomAccount(t)
	secondAccount := createRandomAccount(t)
	fmt.Println(">> before:", firstAccount.Balance, secondAccount.Balance)

	// run n concurrent transfer transactions
	n := 15
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: firstAccount.ID,
				ToAccountID:   secondAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, firstAccount.ID, transfer.FromAccountID)
		require.Equal(t, secondAccount.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.Amount, -amount)
		require.Equal(t, fromEntry.AccountID, firstAccount.ID)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.Amount, amount)
		require.Equal(t, toEntry.AccountID, secondAccount.ID)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check updated account's balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, firstAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, secondAccount.ID, toAccount.ID)

		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := firstAccount.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - secondAccount.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	// check the final updated balance
	updatedFirstAccount, err := store.GetAccount(context.Background(), firstAccount.ID)
	require.NoError(t, err)

	updatedSecondAccount, err := store.GetAccount(context.Background(), secondAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedFirstAccount.Balance, updatedSecondAccount.Balance)
	require.Equal(t, firstAccount.Balance-int64(n)*amount, updatedFirstAccount.Balance)
	require.Equal(t, secondAccount.Balance+int64(n)*amount, updatedSecondAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStoreSQl(testDB)

	firstAccount := createRandomAccount(t)
	secondAccount := createRandomAccount(t)
	fmt.Println(">> before:", firstAccount.Balance, secondAccount.Balance)

	// run n concurrent transfer transactions
	n := 20
	amount := int64(10)

	errs := make(chan error)
	for i := 0; i < n; i++ {
		fromAccountID := firstAccount.ID
		toAccountID := secondAccount.ID

		if i%2 == 0 {
			fromAccountID = secondAccount.ID
			toAccountID = firstAccount.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedFirstAccount, err := store.GetAccount(context.Background(), firstAccount.ID)
	require.NoError(t, err)

	updatedSecondAccount, err := store.GetAccount(context.Background(), secondAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedFirstAccount.Balance, updatedSecondAccount.Balance)
	require.Equal(t, firstAccount.Balance, updatedFirstAccount.Balance)
	require.Equal(t, secondAccount.Balance, updatedSecondAccount.Balance)
}
