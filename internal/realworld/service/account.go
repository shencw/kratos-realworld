package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/shencw/kratos-realworld/api/realworld/v1"
	"github.com/shencw/kratos-realworld/internal/realworld/biz"
)

type AccountService struct {
	v1.UnimplementedAccountServer

	account *biz.AccountUseCase
	log     *log.Helper
}

func NewAccountService(account *biz.AccountUseCase, logger log.Logger) *AccountService {
	log.NewHelper(logger).Info("NewAccountService")

	return &AccountService{account: account, log: log.NewHelper(logger)}
}

func (s *AccountService) Accounts(ctx context.Context, in *v1.AccountsRequest) (*v1.AccountsReply, error) {
	s.log.WithContext(ctx).Infof("Get Accounts")
	accounts, err := s.account.GetAccounts(ctx, in.GetLimit())

	var accountReply []*v1.AccountsReply_Account
	for _, v := range accounts {
		accountReply = append(accountReply, &v1.AccountsReply_Account{
			ID:      v.ID,
			Uid:     v.Uid,
			Type:    v.Type,
			Balance: v.Balance,
			Tag:     v.Tag,
			DT:      v.DT,
			Ctime:   v.Ctime,
		})
	}

	return &v1.AccountsReply{Account: accountReply}, err
}
