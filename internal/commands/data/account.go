package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/shencw/kratos-realworld/api/realworld/v1"
	"github.com/shencw/kratos-realworld/internal/commands/biz"
)

type accountRepo struct {
	data *Data
	log  *log.Helper
}

var _ biz.AccountRepo = (*accountRepo)(nil)

func NewAccountRepo(data *Data, logger log.Logger) biz.AccountRepo {
	return &accountRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accountRepo) GetAccountList(ctx context.Context, limit int32) ([]biz.Account, error) {
	client := v1.NewAccountClient(r.data.realWorldGrpcConn)
	reply, err := client.Accounts(ctx, &v1.AccountsRequest{Limit: &limit})
	if err != nil {
		r.log.Error("GetAccountList Grpc Error: ",err)
		return nil, err
	}
	var accounts []biz.Account
	for _, v := range reply.Account {
		accounts = append(accounts, biz.Account{
			ID:      v.ID,
			Uid:     v.Uid,
			Type:    v.Type,
			Balance: v.Balance,
			Tag:     v.Tag,
			DT:      v.DT,
			Ctime:   v.Ctime,
		})
	}
	return accounts, nil
}
