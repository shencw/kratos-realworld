// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/shencw/kratos-realworld/internal/pkg/conf"
	"github.com/shencw/kratos-realworld/internal/realworld/biz"
	"github.com/shencw/kratos-realworld/internal/realworld/data"
	"github.com/shencw/kratos-realworld/internal/realworld/server"
	"github.com/shencw/kratos-realworld/internal/realworld/service"
)

// initApp init kratos application.
func initApp(log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(conf.ProviderSet, server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
