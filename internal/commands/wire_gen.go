// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package commands

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shencw/kratos-realworld/internal/commands/biz"
	"github.com/shencw/kratos-realworld/internal/commands/data"
	"github.com/shencw/kratos-realworld/internal/commands/pkg"
	"github.com/shencw/kratos-realworld/internal/commands/service"
	"github.com/shencw/kratos-realworld/internal/pkg/conf"
)

// Injectors from wire.go:

func LoadCommandService(ctx context.Context, logger log.Logger) ([]pkg.CommandService, func(), error) {
	confData, cleanup := conf.NewDataConfig(logger)
	client := data.NewRedisConn(confData, logger)
	connection := data.NewHiveConn(confData, logger)
	saramaClient := data.NewKafkaClient(confData, logger)
	server, cleanup2 := conf.NewServerConfig(logger)
	clientConn, err := data.NewRealWorldGrpcConn(ctx, server, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	dataData, cleanup3, err := data.NewData(confData, logger, client, connection, saramaClient, clientConn)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	accountRepo := data.NewAccountRepo(dataData, logger)
	accountUseCase := biz.NewAccountUseCase(accountRepo, logger)
	accountService := service.NewAccountService(ctx, logger, accountUseCase)
	v := service.GetCommandsCollect(accountService)
	return v, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}