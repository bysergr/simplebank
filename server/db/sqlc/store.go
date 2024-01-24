package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type StoreSQL struct {
	*Queries
	db *sql.DB
}

func NewStoreSQl(db *sql.DB) Store {
	return &StoreSQL{
		Queries: New(db),
		db:      db,
	}
}

func (store *StoreSQL) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *StoreSQL) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = transferMoney(
				ctx,
				arg.FromAccountID,
				-arg.Amount,
				arg.ToAccountID,
				arg.Amount,
				queries)
			if err != nil {
				return nil
			}

			return nil
		}

		result.ToAccount, result.FromAccount, err = transferMoney(
			ctx,
			arg.ToAccountID,
			arg.Amount,
			arg.FromAccountID,
			-arg.Amount,
			queries)
		if err != nil {
			return nil
		}

		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	firstAccountID,
	firstAmount,
	secondAccountID,
	secondAmount int64,
	queries *Queries,
) (firstAccount Account, secondAccount Account, err error) {
	// Update the balance to the first account
	firstAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     firstAccountID,
		Amount: firstAmount,
	})
	if err != nil {
		return
	}

	// Update the balance to the second account
	secondAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     secondAccountID,
		Amount: secondAmount,
	})
	if err != nil {
		return
	}

	return
}
