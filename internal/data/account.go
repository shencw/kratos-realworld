package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mitchellh/mapstructure"
	"github.com/shencw/kratos-realworld/internal/biz"
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

func (r *accountRepo) GetAccountList(ctx context.Context) ([]biz.Account, error) {
	cursor := r.data.hiveConn.Cursor()
	defer cursor.Close()
	cursor.Exec(ctx, "select id,uid,type,balance,tag,ctime,dt from ods_exchange_account_da limit 10")
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	var accountData []biz.Account
	for cursor.HasMore(ctx) {
		data := cursor.RowMap(ctx)
		if cursor.Error() != nil {
			if cursor.Error().Error() == "Context is done" {
				r.log.Info("data:%v,err:%s", data, cursor.Error())
			} else {
				r.log.Error("data:%v,err:%s", data, cursor.Error())
			}
			continue
		}
		var account biz.Account
		if err := mapstructure.Decode(data, &account); err != nil {
			r.log.Error("mapstructure.Decode data:%v Error:%s", data, err)
			continue
		}
		accountData = append(accountData, account)
	}

	return accountData, nil
}
