package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/shencw/kratos-realworld/api/realworld/v1"
	"github.com/shencw/kratos-realworld/internal/pkg/conf"
	"github.com/shencw/kratos-realworld/internal/realworld/server/health"
	service2 "github.com/shencw/kratos-realworld/internal/realworld/service"
)

// NewHTTPServer new HTTP server.
func NewHTTPServer(
	c *conf.Server,
	realWorldServer *service2.RealWorldService,
	accountServer *service2.AccountService,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
		http.Logger(logger),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	srv.Handle("/health", health.NewHttpHealthHandler())
	v1.RegisterRealWorldHTTPServer(srv, realWorldServer)
	v1.RegisterAccountHTTPServer(srv, accountServer)
	return srv
}
