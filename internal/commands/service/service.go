package service

import (
	"github.com/google/wire"
	"github.com/shencw/kratos-realworld/internal/commands/pkg"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(GetCommandsCollect, NewAccountService)

func GetCommandsCollect(account *AccountService) []pkg.CommandService {
	var collect []pkg.CommandService
	var accountService pkg.CommandService = account

	return append(collect, accountService)
}
