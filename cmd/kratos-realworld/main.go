package main

import (
	"flag"
	"fmt"
	"github.com/shencw/kratos-realworld/pkg/rollinglog"
	"go.uber.org/zap"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/shencw/kratos-realworld/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagConf is the config flag.
	flagConf string
	// logPath 日志地址
	logPath string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&logPath, "log", "../../logs/", "log path, eg: -log ../../logs/kratos.log")
	flag.Parse()
}

func newApp(logger log.Logger, httpServer *http.Server, gRPCServer *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			httpServer,
			//gRPCServer,
		),
	)
}

func initLogger() *rollinglog.KZapLogger2 {
	kv := []interface{}{
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	}
	return rollinglog.New(&rollinglog.Options{
		Level:            zap.DebugLevel.String(),
		Format:           "json",
		EnableColor:      false,
		DisableCaller:    true,
		OutputPaths:      []string{fmt.Sprintf("%s/test.log", logPath), "stdout"},
		ErrorOutputPaths: []string{fmt.Sprintf("%s/error.log", logPath)},
		Rolling:          true,
		RollingMaxSize:   1,
	}, kv...)
}

func main() {
	zapLogger := initLogger()
	logger := zapLogger.GetLogger()
	defer func () {
		log.NewHelper(logger).Info("Main Stop.")
		zapLogger.Sync()
	}()

	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
		config.WithLogger(logger),
	)

	defer func(c config.Config, logger log.Logger) {
		if err := c.Close(); err != nil {
			log.NewHelper(logger).Error("initFileConfig Close Error: %s", err.Error())
		}
	}(c, logger)

	var bc conf.Bootstrap

	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
