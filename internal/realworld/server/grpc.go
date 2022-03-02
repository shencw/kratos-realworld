package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "github.com/shencw/kratos-realworld/api/realworld/v1"
	"github.com/shencw/kratos-realworld/internal/pkg/conf"
	service2 "github.com/shencw/kratos-realworld/internal/realworld/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server,
	realWorldServer *service2.RealWorldService,
	accountServer *service2.AccountService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
		grpc.Logger(logger),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterRealWorldServer(srv, realWorldServer)
	v1.RegisterAccountServer(srv, accountServer)
	return srv
}
