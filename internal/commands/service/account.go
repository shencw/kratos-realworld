package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shencw/kratos-realworld/internal/commands/biz"
	"github.com/shencw/kratos-realworld/internal/commands/pkg"
	"github.com/spf13/cobra"
)

type AccountService struct {
	account *biz.AccountUseCase
	log     *log.Helper
	command *cobra.Command
}

var _ pkg.CommandService = (*AccountService)(nil)

func (s *AccountService) GetCommands() *cobra.Command {
	return s.command
}

func NewAccountService(ctx context.Context, logger log.Logger, account *biz.AccountUseCase) *AccountService {
	log.NewHelper(logger).Info("NewAccountService")

	service := &AccountService{
		account: account,
		log:     log.NewHelper(logger),
		command: &cobra.Command{
			Use:                   "AccountService",
			DisableFlagsInUseLine: true,
			Short:                 "AccountService Short",
			Long:                  "AccountService Long",
		},
	}

	service.command.AddCommand(service.Accounts(ctx))

	return service
}

func (s *AccountService) Accounts(ctx context.Context) *cobra.Command {
	var limit int32

	cmd := &cobra.Command{
		Use:                   "accounts",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "获取10条数据",
		TraverseChildren:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := s.account.GetAccounts(ctx, limit)
			if err != nil {
				return err
			}
			// TODO: 格式化处理
			fmt.Println(list)
			return nil
		},
		SuggestFor: []string{},
	}

	cmd.Flags().Int32VarP(&limit, "limit", "l", 10, "Specify the amount records to be returned.")

	return cmd
}
