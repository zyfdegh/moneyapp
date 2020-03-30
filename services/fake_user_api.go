package services

import (
	"errors"
	"math/rand"
	"time"

	"github.com/zyfdegh/moneyapp/models"
)

var (
	ErrAuthfailed = errors.New("用户名或密码错误")
)

type FakeUserAPI struct {
	users []*models.User
}

func NewFakeUserAPI() *FakeUserAPI {
	api := &FakeUserAPI{}
	api.genFakeData()
	return api
}

func (this *FakeUserAPI) genFakeData() {
	this.users = append(this.users, &models.User{Username: "1", Password: "1"})
	this.users = append(this.users, &models.User{Username: "admin", Password: "admin"})
	this.users = append(this.users, &models.User{Username: "zyf", Password: "111111"})
	this.users = append(this.users, &models.User{Username: "qhw", Password: "111111"})
}

func (this *FakeUserAPI) Login(username, password string) (*models.Session, error) {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	for _, v := range this.users {
		if v.Username == username && v.Password == password {
			return &models.Session{
				Username:  username,
				ExpiredAt: time.Now().Unix() + int64(12*time.Hour),
			}, nil
		}
	}
	return nil, ErrAuthfailed
}

func (this *FakeUserAPI) Register(u *models.User) error {
	if err := u.VerifyNewAccount(); err != nil {
		return err
	}
	for _, v := range this.users {
		if v.Username == u.Username {
			return errors.New("相同用户名已存在")
		}
	}
	this.users = append(this.users, u)
	return nil
}

func (this *FakeUserAPI) Logout(username, password string) {

}
