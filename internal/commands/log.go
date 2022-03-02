package commands

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shencw/kratos-realworld/pkg/rollinglog"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"os"
)

var (
	logPath string
	Name    string
	Version string
	id, _   = os.Hostname()
)

func addLogFlags(flags *pflag.FlagSet) {
	flags.StringVar(
		&logPath,
		"log-output",
		"../../logs/",
		"log output location",
	)
}

func newLogger() (log.Logger, func()) {
	KZapLogger := rollinglog.New(&rollinglog.Options{
		Level:            zap.DebugLevel.String(),
		Format:           "json",
		EnableColor:      false,
		DisableCaller:    true,
		OutputPaths:      []string{fmt.Sprintf("%s/commands_kratos.log", logPath)},
		ErrorOutputPaths: []string{fmt.Sprintf("%s/commands_error.log", logPath)},
		Rolling:          true,
		RollingMaxSize:   1,
	}, []interface{}{
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	}...)
	return KZapLogger.GetLogger(), func() {
		log.NewHelper(KZapLogger.GetLogger()).Info("Main Stop.")

		KZapLogger.Sync()
	}
}
