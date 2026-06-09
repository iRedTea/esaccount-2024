package service

import (
	"esaccount"
)

type Authorization interface {
	Authorize(header string) (*esaccount.AuthorizedUser, error)
	AuthorizeById(id int64) (*esaccount.AuthorizedUser, error)
	AuthorizeAndUpdatePicture(header string, picUrl string) (*esaccount.AuthorizedUser, error)
}

type Service struct {
	Authorization
}

func NewService() *Service {
	return &Service{
		Authorization: NewAuthService(),
	}
}
