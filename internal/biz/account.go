package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Account struct {
	ID      int32
	Uid     int32
	Type    int32
	Balance string
	Tag     string
	DT      string
	Ctime   string
}

type AccountRepo interface {
	GetAccountList(context.Context) ([]Account, error)
}

type AccountUseCase struct {
	account AccountRepo
	log     *log.Helper
}

func NewAccountUseCase(account AccountRepo, logger log.Logger) *AccountUseCase {
	return &AccountUseCase{account: account, log: log.NewHelper(logger)}
}

func (uc *AccountUseCase) GetAccounts(ctx context.Context) ([]Account, error) {
	return uc.account.GetAccountList(ctx)
}
