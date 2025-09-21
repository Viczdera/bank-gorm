package gorm

import (
	"context"

	"gorm.io/gorm"
)

// Store provides all functions to execute database operations
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error)
}

// Querier interface defines all database query methods
type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccount(ctx context.Context, ID int64) error
	GetAccount(ctx context.Context, ID int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, ID int64) (Account, error)
	GetEntry(ctx context.Context, ID int64) (Entry, error)
	GetTransfer(ctx context.Context, ID int64) (Transfer, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

// GormStore implements Store interface using GORM
type GormStore struct {
	db *gorm.DB
}

// NewStore creates a new GormStore instance
func NewStore(db *gorm.DB) Store {
	return &GormStore{
		db: db,
	}
}

// User operations

func (store *GormStore) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	user := User{
		Username:       arg.Username,
		PasswordHashed: arg.PasswordHashed,
		FullName:       arg.FullName,
		Email:          arg.Email,
	}

	result := store.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

func (store *GormStore) GetUser(ctx context.Context, username string) (User, error) {
	var user User
	result := store.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

// Account operations

func (store *GormStore) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	account := Account{
		Owner:    arg.Owner,
		Balance:  arg.Balance,
		Currency: arg.Currency,
	}

	result := store.db.WithContext(ctx).Create(&account)
	if result.Error != nil {
		return Account{}, result.Error
	}

	return account, nil
}

func (store *GormStore) GetAccount(ctx context.Context, ID int64) (Account, error) {
	var account Account
	result := store.db.WithContext(ctx).Where("_id = ?", ID).First(&account)
	if result.Error != nil {
		return Account{}, result.Error
	}

	return account, nil
}

func (store *GormStore) GetAccountForUpdate(ctx context.Context, ID int64) (Account, error) {
	var account Account
	result := store.db.WithContext(ctx).Set("gorm:query_option", "FOR NO KEY UPDATE").Where("_id = ?", ID).First(&account)
	if result.Error != nil {
		return Account{}, result.Error
	}

	return account, nil
}

func (store *GormStore) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	var accounts []Account
	result := store.db.WithContext(ctx).Where("owner = ?", arg.Owner).
		Order("_id").
		Limit(int(arg.Limit)).
		Offset(int(arg.Offset)).
		Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (store *GormStore) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	var account Account
	result := store.db.WithContext(ctx).Model(&account).Where("_id = ?", arg.ID).
		Update("balance", arg.Balance).
		Where("_id = ?", arg.ID).First(&account)

	if result.Error != nil {
		return Account{}, result.Error
	}

	return account, nil
}

func (store *GormStore) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	var account Account

	// First get the current account
	if err := store.db.WithContext(ctx).Where("_id = ?", arg.ID).First(&account).Error; err != nil {
		return Account{}, err
	}

	// Update the balance
	result := store.db.WithContext(ctx).Model(&account).Where("_id = ?", arg.ID).
		Update("balance", gorm.Expr("balance + ?", arg.Amount))

	if result.Error != nil {
		return Account{}, result.Error
	}

	// Get the updated account
	if err := store.db.WithContext(ctx).Where("_id = ?", arg.ID).First(&account).Error; err != nil {
		return Account{}, err
	}

	return account, nil
}

func (store *GormStore) DeleteAccount(ctx context.Context, ID int64) error {
	result := store.db.WithContext(ctx).Where("_id = ?", ID).Delete(&Account{})
	return result.Error
}

// Entry operations

func (store *GormStore) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	entry := Entry{
		AccountID: arg.AccountID,
		Amount:    arg.Amount,
	}

	result := store.db.WithContext(ctx).Create(&entry)
	if result.Error != nil {
		return Entry{}, result.Error
	}

	return entry, nil
}

func (store *GormStore) GetEntry(ctx context.Context, ID int64) (Entry, error) {
	var entry Entry
	result := store.db.WithContext(ctx).Where("_id = ?", ID).First(&entry)
	if result.Error != nil {
		return Entry{}, result.Error
	}

	return entry, nil
}

func (store *GormStore) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	var entries []Entry
	result := store.db.WithContext(ctx).
		Order("account_id").
		Limit(int(arg.Limit)).
		Offset(int(arg.Offset)).
		Find(&entries)

	if result.Error != nil {
		return nil, result.Error
	}

	return entries, nil
}

// Transfer operations

func (store *GormStore) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	transfer := Transfer{
		FromAccount: arg.FromAccount,
		ToAccount:   arg.ToAccount,
		Amount:      arg.Amount,
	}

	result := store.db.WithContext(ctx).Create(&transfer)
	if result.Error != nil {
		return Transfer{}, result.Error
	}

	return transfer, nil
}

func (store *GormStore) GetTransfer(ctx context.Context, ID int64) (Transfer, error) {
	var transfer Transfer
	result := store.db.WithContext(ctx).Where("_id = ?", ID).First(&transfer)
	if result.Error != nil {
		return Transfer{}, result.Error
	}

	return transfer, nil
}

// execTx executes a function within a database transaction
func (store *GormStore) execTx(ctx context.Context, fn func(*gorm.DB) error) error {
	return store.db.WithContext(ctx).Transaction(fn)
}

// TransferTx performs a money transfer from one account to another
// It creates a transfer record, add account entries, and update accounts' balance within a transaction
func (store *GormStore) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(tx *gorm.DB) error {
		var err error

		// Create the transfer
		result.Transfer, err = store.createTransferTx(ctx, tx, CreateTransferParams{
			FromAccount: arg.FromAccount,
			ToAccount:   arg.ToAccount,
			Amount:      arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create entry for the sender (negative amount)
		result.FromEntry, err = store.createEntryTx(ctx, tx, CreateEntryParams{
			AccountID: arg.FromAccount,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create entry for the receiver (positive amount)
		result.ToEntry, err = store.createEntryTx(ctx, tx, CreateEntryParams{
			AccountID: arg.ToAccount,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update account balances
		// To avoid deadlock, always update accounts in consistent order (smaller ID first)
		if arg.FromAccount < arg.ToAccount {
			result.FromAccount, err = store.addAccountBalanceTx(ctx, tx, AddAccountBalanceParams{
				ID:     arg.FromAccount,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = store.addAccountBalanceTx(ctx, tx, AddAccountBalanceParams{
				ID:     arg.ToAccount,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = store.addAccountBalanceTx(ctx, tx, AddAccountBalanceParams{
				ID:     arg.ToAccount,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = store.addAccountBalanceTx(ctx, tx, AddAccountBalanceParams{
				ID:     arg.FromAccount,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

// Helper methods for transaction operations

func (store *GormStore) createTransferTx(ctx context.Context, tx *gorm.DB, arg CreateTransferParams) (Transfer, error) {
	transfer := Transfer{
		FromAccount: arg.FromAccount,
		ToAccount:   arg.ToAccount,
		Amount:      arg.Amount,
	}

	result := tx.WithContext(ctx).Create(&transfer)
	if result.Error != nil {
		return Transfer{}, result.Error
	}

	return transfer, nil
}

func (store *GormStore) createEntryTx(ctx context.Context, tx *gorm.DB, arg CreateEntryParams) (Entry, error) {
	entry := Entry{
		AccountID: arg.AccountID,
		Amount:    arg.Amount,
	}

	result := tx.WithContext(ctx).Create(&entry)
	if result.Error != nil {
		return Entry{}, result.Error
	}

	return entry, nil
}

func (store *GormStore) addAccountBalanceTx(ctx context.Context, tx *gorm.DB, arg AddAccountBalanceParams) (Account, error) {
	var account Account

	// First get the current account
	if err := tx.WithContext(ctx).Where("_id = ?", arg.ID).First(&account).Error; err != nil {
		return Account{}, err
	}

	// Update the balance
	result := tx.WithContext(ctx).Model(&account).Where("_id = ?", arg.ID).
		Update("balance", gorm.Expr("balance + ?", arg.Amount))

	if result.Error != nil {
		return Account{}, result.Error
	}

	// Get the updated account
	if err := tx.WithContext(ctx).Where("_id = ?", arg.ID).First(&account).Error; err != nil {
		return Account{}, err
	}

	return account, nil
}
