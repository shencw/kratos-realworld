package main

import (
	"flag"
	"fmt"
	"github.com/shencw/kratos-realworld/pkg/rollinglog"
	"go.uber.org/zap"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// logPath 日志地址
	logPath string

	id, _ = os.Hostname()
)

func init() {
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

	app, cleanup, err := initApp(logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
