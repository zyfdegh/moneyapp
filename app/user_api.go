package app

import (
	"github.com/zyfdegh/moneyapp/models"
)

type UserAPI interface {
	Login(username, password string) (*models.Session, error)
	Register(u *models.User) error
}
