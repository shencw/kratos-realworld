//go:build wireinject
// +build wireinject

package commands

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/shencw/kratos-realworld/internal/commands/biz"
	"github.com/shencw/kratos-realworld/internal/commands/data"
	"github.com/shencw/kratos-realworld/internal/commands/pkg"
	"github.com/shencw/kratos-realworld/internal/commands/service"
	"github.com/shencw/kratos-realworld/internal/pkg/conf"
)

func LoadCommandService(ctx context.Context, logger log.Logger) ([]pkg.CommandService, func(), error) {
	panic(wire.Build(conf.ProviderSet, biz.ProviderSet, service.ProviderSet, data.ProviderSet))
}
