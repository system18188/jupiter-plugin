package xrestful

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
	restful "github.com/emicklei/go-restful"
	"github.com/pkg/errors"
)

//ModName named a mod
const ModName = "server.go-restful"

// Config HTTP config
type Config struct {
	// 绑定地址
	Host string `json:"host" toml:"host"`
	// 绑定端口
	Port       int    `json:"port" toml:"port"`
	Deployment string `json:"deployment" toml:"deployment"`
	Debug      bool   `json:"debug" toml:"debug"`
	// 测量请求响应时间
	DisableMetric bool `json:"disableMetric" toml:"disableMetric"`
	// 跟踪
	DisableTrace bool `json:"disableTrace" toml:"disableTrace"`
	// 开启gzip 压缩
	EnableGzip bool `json:"enableGzip" toml:"enableGzip"`
	// ServiceAddress service address in registry info, default to 'Host:Port'
	ServiceAddress string `json:"serviceAddress" toml:"serviceAddress"`

	SlowQueryThresholdInMilli int64 `json:"slowQueryThresholdInMilli" toml:"slowQueryThresholdInMilli"`

	logger *xlog.Logger
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                      flag.String("host"),
		Port:                      9091,
		Debug:                     false,
		Deployment:                constant.DefaultDeployment,
		SlowQueryThresholdInMilli: 500, // 500ms
		logger:                    xlog.JupiterLogger.With(xlog.FieldMod(ModName)),
	}
}

// StdConfig Jupiter Standard HTTP Server config
func StdConfig(name string) *Config {
	return RawConfig("jupiter.server." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil &&
		errors.Cause(err) != conf.ErrInvalidKey {
		config.logger.Panic("http server parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(key), xlog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *xlog.Logger) *Config {
	config.logger = logger
	return config
}

// WithHost ...
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort ...
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	server := newServer(config)
	restful.DefaultContainer.Filter(recoverMiddleware(config.logger, config.SlowQueryThresholdInMilli))

	if !config.DisableMetric {
		restful.DefaultContainer.Filter(metricServerInterceptor())
	}

	if !config.DisableTrace {
		restful.DefaultContainer.Filter(traceServerInterceptor())
	}

	if config.EnableGzip {
		restful.DefaultContainer.EnableContentEncoding(true)
	}
	return server
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
