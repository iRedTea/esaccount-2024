package repository

import (
	"esaccount"
	"fmt"

	"gorm.io/gorm"
)

type AccountSchema struct {
	db *gorm.DB
}

func NewAccountSchema(db *gorm.DB) *AccountSchema {
	return &AccountSchema{db: db}
}

func (r AccountSchema) CreateAccount(id int64) (*esaccount.Account, error) {
	account := &esaccount.Account{
		Id:          id,
		Description: "",
		DateOfBirth: "",
		Follows:     []int64{},
		Followers:   []int64{},
		Confidentials: map[string]esaccount.ConfidentialType{
			"username":      "ALL",
			"first_name":    "ALL",
			"last_name":     "ALL",
			"email":         "NOBODY",
			"description":   "ALL",
			"date_of_birth": "FRIENDS",
			"follows":       "ALL",
			"followers":     "ALL",
		},
	}
	err := r.db.Save(account).Error
	return account, err
}

func (r AccountSchema) SaveAccount(user *esaccount.Account) (int64, error) {
	err := r.db.Save(user).Error
	return user.Id, err
}

func (r AccountSchema) GetAccount(id int64) (*esaccount.Account, error) {
	user := esaccount.Account{}
	err := r.db.Find(&user, "id = ?", id).Error
	if user.Id == -1 {
		err = fmt.Errorf("unknown account %d", id)
	}
	return &user, err
}

func (r AccountSchema) FindAccount(query, param string) (esaccount.Account, error) {
	user := esaccount.Account{}
	err := r.db.Find(&user, query, param).Error
	return user, err
}

func (r AccountSchema) GetAccounts(start, limit int) []int64 {
	result := []esaccount.Account{}
	r.db.Select("id").Limit(limit).Offset(start).Find(&result)
	var ans = make([]int64, len(result))
	for i, res := range result {
		ans[i] = res.Id
	}
	return ans
}

func (r AccountSchema) GetAccountsCount() int64 {
	var count int64 = 0
	r.db.Model(&esaccount.Account{}).Count(&count)
	return count
}

func (r AccountSchema) DeleteAccount(id int64) error {
	user := esaccount.Account{}
	err := r.db.Find(&user, "id = ?", id).Error
	if user.Id == -1 {
		err = fmt.Errorf("unknown user %d", id)
	}
	r.db.Delete(user)
	return err
}
