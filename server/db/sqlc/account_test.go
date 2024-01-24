package db

import (
	"context"
	"database/sql"
	"github.com/bysergr/simple-bank/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	randomUser := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    randomUser.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), randomAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randomAccount.ID, account.ID)
	require.Equal(t, randomAccount.Owner, account.Owner)
	require.Equal(t, randomAccount.Balance, account.Balance)
	require.Equal(t, randomAccount.Currency, account.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)

	updateAccountParams := UpdateAccountParams{
		ID:      randomAccount.ID,
		Balance: utils.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), updateAccountParams)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randomAccount.ID, account.ID)
	require.Equal(t, randomAccount.Owner, account.Owner)
	require.Equal(t, updateAccountParams.Balance, account.Balance)
	require.Equal(t, randomAccount.Currency, account.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), randomAccount.ID)
	require.NoError(t, err)

	databaseAccount, err := testQueries.GetAccount(context.Background(), randomAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, databaseAccount)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	listAccountsParams := ListAccountsParams{
		Limit:  5,
		Offset: 0,
		Owner:  lastAccount.Owner,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), listAccountsParams)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
