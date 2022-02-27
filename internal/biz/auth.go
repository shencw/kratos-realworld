package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID           uint64
	Username     string
	Email        string
	Password     string
	Image        string
	PasswordHash string
}

type UserLogin struct {
	Email    string
	Username string
	Token    string
	Bio      string
	Image    string
}

type UserRepo interface {
	Login(ctx context.Context, user User) error
	CreateUser(ctx context.Context, user *User) (uint64, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

type AuthUseCase struct {
	user UserRepo
	log  *log.Helper
}

func NewAuthUseCase(user UserRepo, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{user: user, log: log.NewHelper(logger)}
}

// Register 用户注册
func (uc *AuthUseCase) Register(ctx context.Context, username, email, password string) (*UserLogin, error) {
	u, err := uc.user.GetUserByUsername(ctx, username)
	if err != nil {
		return &UserLogin{}, err
	}
	if u.ID > 0 {
		return &UserLogin{}, errors.New(500,"已经创建", "")
	}
	u = &User{
		Email:        email,
		Username:     username,
		PasswordHash: password,
	}
	id, err := uc.user.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	uc.log.Infof("Register User %s:%v", username, id)
	return &UserLogin{
		Email:    email,
		Username: username,
	}, nil
}
