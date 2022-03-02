package conf

import (
	"flag"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"sync"
)

var (
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// ProviderSet is config providers.
var ProviderSet = wire.NewSet(
	NewServerConfig,
	NewDataConfig,
)

type fileConfig struct {
	fileConfig config.Config
	bc         Bootstrap
	insOnce    sync.Once
	closeOnce  sync.Once
}

var fileConfigIns fileConfig

func NewServerConfig(logger log.Logger) (*Server, func()) {
	ins := getFileConfigIns(logger)
	return ins.bc.GetServer(), func() {
		ins.Close(logger)
	}
}

func NewDataConfig(logger log.Logger) (*Data, func()) {
	ins := getFileConfigIns(logger)
	return ins.bc.GetData(), func() {
		ins.Close(logger)
	}
}

func getFileConfigIns(logger log.Logger) *fileConfig {
	fileConfigIns.insOnce.Do(func() {
		fileConfigIns.fileConfig = config.New(
			config.WithSource(
				file.NewSource(flagConf),
			),
			config.WithLogger(logger),
		)
		if err := fileConfigIns.fileConfig.Load(); err != nil {
			panic(err)
		}
		if err := fileConfigIns.fileConfig.Scan(&fileConfigIns.bc); err != nil {
			panic(err)
		}
	})

	return &fileConfigIns
}

func (c *fileConfig) Close(logger log.Logger) {
	c.closeOnce.Do(func() {
		if err := c.fileConfig.Close(); err != nil {
			log.NewHelper(logger).Errorf("FileConfig Close Error: %s", err)
		}
	})
}
