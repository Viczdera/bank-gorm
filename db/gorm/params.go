package gorm

// user params
type CreateUserParams struct {
	Username       string `json:"username"`
	PasswordHashed string `json:"password_hashed"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

// acount params
type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type UpdateAccountParams struct {
	ID      int64 `json:"_id"`
	Balance int64 `json:"balance"`
}

type AddAccountBalanceParams struct {
	ID     int64 `json:"_id"`
	Amount int64 `json:"amount"`
}

/*
entry params
*/
type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type ListEntriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

/*
transfer params
*/
type CreateTransferParams struct {
	FromAccount int64 `json:"from_account"`
	ToAccount   int64 `json:"to_account"`
	Amount      int64 `json:"amount"`
}

/*
transfer transaction params
*/
type TransferTxParams struct {
	FromAccount int64 `json:"from_account"`
	ToAccount   int64 `json:"to_account"`
	Amount      int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
}
