// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shencw/kratos-realworld/internal/pkg/conf"
	"github.com/shencw/kratos-realworld/internal/realworld/biz"
	"github.com/shencw/kratos-realworld/internal/realworld/data"
	"github.com/shencw/kratos-realworld/internal/realworld/server"
	"github.com/shencw/kratos-realworld/internal/realworld/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(logger log.Logger) (*kratos.App, func(), error) {
	confServer, cleanup := conf.NewServerConfig(logger)
	confData, cleanup2 := conf.NewDataConfig(logger)
	db := data.NewRealWorldDB(confServer, confData, logger)
	client := data.NewRedisConn(confData, logger)
	connection := data.NewHiveConn(confData, logger)
	saramaClient := data.NewKafkaClient(confData, logger)
	dataData, cleanup3, err := data.NewData(confData, logger, db, client, connection, saramaClient)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	userRepo := data.NewAuthRepo(dataData, logger)
	authUseCase := biz.NewAuthUseCase(userRepo, logger)
	realWorldService := service.NewRealWorldService(authUseCase, logger)
	accountRepo := data.NewAccountRepo(dataData, logger)
	accountUseCase := biz.NewAccountUseCase(accountRepo, logger)
	accountService := service.NewAccountService(accountUseCase, logger)
	httpServer := server.NewHTTPServer(confServer, realWorldService, accountService, logger)
	grpcServer := server.NewGRPCServer(confServer, realWorldService, accountService, logger)
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
