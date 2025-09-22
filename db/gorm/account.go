package gorm

import (
	"context"

	"gorm.io/gorm"
)

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

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

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
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

type UpdateAccountParams struct {
	ID      int64 `json:"_id"`
	Balance int64 `json:"balance"`
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

type AddAccountBalanceParams struct {
	ID     int64 `json:"_id"`
	Amount int64 `json:"amount"`
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
