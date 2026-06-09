package repository

import (
	"esaccount"

	"gorm.io/gorm"
)

type Account interface {
	CreateAccount(int64) (*esaccount.Account, error)
	SaveAccount(user *esaccount.Account) (int64, error)
	GetAccount(id int64) (*esaccount.Account, error)
	FindAccount(query, param string) (esaccount.Account, error)
	GetAccounts(start, limit int) []int64
	DeleteAccount(id int64) error
	GetAccountsCount() int64
}

type Repository struct {
	Account
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Account: NewAccountSchema(db),
	}
}
