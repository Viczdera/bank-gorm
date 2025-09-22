package gorm

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username          string    `gorm:"primaryKey;type:varchar;not null" json:"username"`
	PasswordHashed    string    `gorm:"column:password_hashed;type:varchar;not null" json:"password_hashed"`
	FullName          string    `gorm:"column:full_name;type:varchar;not null" json:"full_name"`
	Email             string    `gorm:"type:varchar;not null;uniqueIndex" json:"email"`
	PasswordChangedAt time.Time `gorm:"column:password_changed_at;type:timestamptz" json:"password_changed_at"`
	CreatedAt         time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()" json:"created_at"`

	Accounts []Account `gorm:"foreignKey:Owner;references:Username" json:"accounts,omitempty"`
}

type Account struct {
	ID        int64     `gorm:"primaryKey;column:_id;autoIncrement" json:"_id"`
	Balance   int64     `gorm:"type:bigint;not null;default:0" json:"balance"`
	Owner     string    `gorm:"type:varchar;not null;index" json:"owner"`
	Currency  string    `gorm:"type:varchar;not null" json:"currency"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`

	User          User       `gorm:"foreignKey:Owner;references:Username" json:"user,omitempty"`
	Entries       []Entry    `gorm:"foreignKey:AccountID;references:ID" json:"entries,omitempty"`
	FromTransfers []Transfer `gorm:"foreignKey:FromAccount;references:ID" json:"from_transfers,omitempty"`
	ToTransfers   []Transfer `gorm:"foreignKey:ToAccount;references:ID" json:"to_transfers,omitempty"`
}

type Entry struct {
	ID        int64     `gorm:"primaryKey;column:_id;autoIncrement" json:"_id"`
	AccountID int64     `gorm:"column:account_id;type:bigint;not null;index" json:"account_id"`
	Amount    int64     `gorm:"type:bigint;not null" json:"amount"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`

	Account Account `gorm:"foreignKey:AccountID;references:ID" json:"account"`
}

type Transfer struct {
	ID          int64     `gorm:"primaryKey;column:_id;autoIncrement" json:"_id"`
	FromAccount int64     `gorm:"column:from_account;type:bigint;not null;index" json:"from_account"`
	ToAccount   int64     `gorm:"column:to_account;type:bigint;not null;index" json:"to_account"`
	Amount      int64     `gorm:"type:bigint;not null;check:amount > 0" json:"amount"` // must be positive
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`

	// FromAccountRel Account `gorm:"foreignKey:FromAccount;references:ID" json:"from_account_rel,omitempty"`
	// ToAccountRel   Account `gorm:"foreignKey:ToAccount;references:ID" json:"to_account_rel,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (Account) TableName() string {
	return "accounts"
}

func (Entry) TableName() string {
	return "entries"
}

func (Transfer) TableName() string {
	return "transfers"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.PasswordChangedAt = time.Time{}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	return nil
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	return nil
}

func (e *Entry) BeforeCreate(tx *gorm.DB) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	return nil
}

func (t *Transfer) BeforeCreate(tx *gorm.DB) error {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return nil
}
