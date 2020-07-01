package dbr

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
	"time"
)

// StdConfig 标准配置，规范配置文件头
func StdConfig(name string) *Config {
	return RawConfig("jupiter.dbr." + name)
}

// RawConfig 传入mapstructure格式的配置
// example: RawConfig("jupiter.dbr.stt_config")
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config, conf.TagName("toml")); err != nil {
		xlog.Panic("unmarshal key", xlog.FieldMod("dbr"), xlog.FieldErr(err), xlog.FieldKey(key))
	}

	return config
}

// config options
type Config struct {
	// 数据库驱动 mysql postgres or pgx sqlite3
	Drive string `json:"drive" toml:"drive"`
	// 数据库连接
	DSN string `json:"dsn" toml:"dsn"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" toml:"maxOpenConns"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" toml:"maxIdleConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" toml:"connMaxLifetime"`

	// 记录错误sql时,是否打印包含参数的完整sql语句
	// select * from aid = ?;
	// select * from aid = 288016;
	DetailSQL bool `json:"detailSql" toml:"detailSql"`

	logger       *xlog.Logger
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DSN:             "",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: xtime.Duration("300s"),
		logger:          xlog.DefaultLogger,
	}
}

// WithLogger ...
func (config *Config) WithLogger(log *xlog.Logger) *Config {
	config.logger = log
	return config
}

// Build ...
func (config *Config) Build() *Connection {
	conn, err := Open(config.Drive, config.DSN, JupiterReceiver)
	if err != nil {
		config.logger.Panic(fmt.Sprintf("open ",config.Drive), xlog.FieldMod("dbr"), xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldValueAny(config))
	}
	conn.SetMaxIdleConns(config.MaxIdleConns)
	conn.SetMaxOpenConns(config.MaxOpenConns)
	conn.SetConnMaxLifetime(config.ConnMaxLifetime * time.Second)
	return conn
}
