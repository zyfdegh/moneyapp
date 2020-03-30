package models

import (
	"errors"
)

type User struct {
	Username  string
	Password  string
	Realname  string
	Cellphone string
}

func (u *User) VerifyNewAccount() error {
	if u == nil {
		return errors.New("空用户")
	}
	if u.Username == "" {
		return errors.New("用户名为空")
	}
	if len(u.Password) < 6 {
		return errors.New("密码少于 6 位")
	}
	return nil
}
