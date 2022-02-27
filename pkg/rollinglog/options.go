package rollinglog

import (
	"go.uber.org/zap/zapcore"
)

const (
	consoleFormat = "console"
	jsonFormat    = "json"
)

type Options struct {
	OutputPaths      []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level            string   `json:"level"              mapstructure:"level"`
	Format           string   `json:"format"             mapstructure:"format"`
	DisableCaller    bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	EnableColor      bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development      bool     `json:"development"        mapstructure:"development"`
	Name             string   `json:"name"               mapstructure:"name"`

	Rolling           bool `json:"rolling" mapstructure:"rolling"`
	RollingMaxSize    int  `json:"rolling-max-size" mapstructure:"rolling-max-size"`
	RollingMaxAge     int  `json:"rolling-max-age" mapstructure:"rolling-max-age"`
	RollingMaxBackups int  `json:"rolling-max-backups" mapstructure:"rolling-max-backups"`
	RollingLocalTime  bool `json:"rolling-local-time" mapstructure:"rolling-local-time"`
	RollingCompress   bool `json:"rolling-compress" mapstructure:"rolling-compress"`
}

// NewOptions creates an Options object with default parameters.
func NewOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel.String(),
		Format:           jsonFormat,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EnableColor:      false,
		DisableCaller:    true,
		Rolling:          true,
		RollingMaxSize:   1,
	}
}
