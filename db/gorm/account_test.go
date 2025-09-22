package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/Viczdera/bank-gorm/util"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createTestAccount(t *testing.T) Account {
	user := createTestUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	acct1 := createTestAccount(t)
	acctFind, err := testStore.GetAccount(context.Background(), acct1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acctFind)

	require.Equal(t, acct1.ID, acctFind.ID)
	require.Equal(t, acct1.Balance, acctFind.Balance)
	require.Equal(t, acct1.Currency, acctFind.Currency)
	require.Equal(t, acct1.Owner, acctFind.Owner)
	require.WithinDuration(t, acct1.CreatedAt, acctFind.CreatedAt, time.Second)
}

func TestGetAccountNotFound(t *testing.T) {
	_, err := testStore.GetAccount(context.Background(), 999999)
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUpdateAccount(t *testing.T) {
	acct1 := createTestAccount(t)
	args := UpdateAccountParams{
		ID:      acct1.ID,
		Balance: util.RandomBalance(),
	}

	acctUpdate, err := testStore.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, acctUpdate)

	require.Equal(t, args.Balance, acctUpdate.Balance)
	require.Equal(t, acct1.ID, acctUpdate.ID)
	require.Equal(t, acct1.Currency, acctUpdate.Currency)
	require.Equal(t, acct1.Owner, acctUpdate.Owner)
	require.WithinDuration(t, acct1.CreatedAt, acctUpdate.CreatedAt, time.Second)
}

func TestAddAccountBalance(t *testing.T) {
	acct1 := createTestAccount(t)
	originalBalance := acct1.Balance
	amount := util.RandomBalance()

	args := AddAccountBalanceParams{
		ID:     acct1.ID,
		Amount: amount,
	}

	acctUpdate, err := testStore.AddAccountBalance(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, acctUpdate)

	require.Equal(t, originalBalance+amount, acctUpdate.Balance)
	require.Equal(t, acct1.ID, acctUpdate.ID)
	require.Equal(t, acct1.Currency, acctUpdate.Currency)
	require.Equal(t, acct1.Owner, acctUpdate.Owner)
}

func TestDeleteAccount(t *testing.T) {
	acct1 := createTestAccount(t)
	err := testStore.DeleteAccount(context.Background(), acct1.ID)

	require.NoError(t, err)

	acctFind, err := testStore.GetAccount(context.Background(), acct1.ID)

	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
	require.Empty(t, acctFind)
}

func TestListAccounts(t *testing.T) {
	// Create a user first
	user := createTestUser(t)

	// Create multiple accounts for the same user
	var accounts []Account
	for i := 0; i < 10; i++ {
		args := CreateAccountParams{
			Owner:    user.Username,
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		account, err := testStore.CreateAccount(context.Background(), args)
		require.NoError(t, err)
		accounts = append(accounts, account)
	}

	args := ListAccountsParams{
		Owner:  user.Username,
		Limit:  5,
		Offset: 0,
	}

	foundAccounts, err := testStore.ListAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, foundAccounts, 5)

	for _, account := range foundAccounts {
		require.NotEmpty(t, account)
		require.Equal(t, user.Username, account.Owner)
	}
}

func TestListAccountsWithPagination(t *testing.T) {
	// Create a user first
	user := createTestUser(t)

	// Create multiple accounts for the same user
	for i := 0; i < 10; i++ {
		args := CreateAccountParams{
			Owner:    user.Username,
			Balance:  util.RandomBalance(),
			Currency: util.RandomCurrency(),
		}
		_, err := testStore.CreateAccount(context.Background(), args)
		require.NoError(t, err)
	}

	// Test first page
	args1 := ListAccountsParams{
		Owner:  user.Username,
		Limit:  5,
		Offset: 0,
	}

	accounts1, err := testStore.ListAccounts(context.Background(), args1)
	require.NoError(t, err)
	require.Len(t, accounts1, 5)

	// Test second page
	args2 := ListAccountsParams{
		Owner:  user.Username,
		Limit:  5,
		Offset: 5,
	}

	accounts2, err := testStore.ListAccounts(context.Background(), args2)
	require.NoError(t, err)
	require.Len(t, accounts2, 5)

	// Ensure different accounts on different pages
	for _, acc1 := range accounts1 {
		for _, acc2 := range accounts2 {
			require.NotEqual(t, acc1.ID, acc2.ID)
		}
	}
}
