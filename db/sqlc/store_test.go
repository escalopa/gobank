package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1, account2 := createRandomAccount(t), createRandomAccount(t)
	amount := int64(10)

	results, errs := make(chan TransferTxResult), make(chan error)
	existed := make(map[int]bool)

	n := 5
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		validateTransferBasic(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// Check entry
		// From
		validateEntryBasic(t, result.FromEntry)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)

		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		// To
		validateEntryBasic(t, result.ToEntry)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// Check account
		// From
		validateAccountBasic(t, result.FromAccount)
		require.Equal(t, account1.ID, result.FromAccount.ID)

		// To
		validateAccountBasic(t, result.ToAccount)
		require.Equal(t, account2.ID, result.ToAccount.ID)

		// Check new balances
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2) // The diff in balance of the 2 accounts MUST be equal
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		// Check the no concurrent transactions used the same value of `K`
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check for updated balance

	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-(amount*int64(n)), updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+(amount*int64(n)), updatedToAccount.Balance)
}
