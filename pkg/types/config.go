package types

import "github.com/yanshicheng/ikube-gin-starter/pkg/logger"

type AppConfig struct {
	HttpAddr          string `mapstructure:"http_addr" json:"http_addr" yaml:"http_addr" env:"APP_HTTP_ADDR"`
	HttpPort          int    `mapstructure:"http_port" json:"http_port" http_port:"http_port" env:"APP_HTTP_PORT"`
	Language          string `mapstructure:"language" json:"language" yaml:"language" env:"APP_LANGUAGE"`
	MaxHeaderSize     int    `mapstructure:"max_header_size" json:"max_header_size" yaml:"max_header_size" env:"APP_MAX_HEADER_SIZE"`
	ReadTimeout       int    `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout" env:"APP_READ_TIMEOUT"`
	ReadHeaderTimeout int    `mapstructure:"read_header_timeout" json:"read_header_timeout" yaml:"read_header_timeout" env:"APP_READ_HEADER_TIMEOUT"`
	WriteTimeout      int    `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout" env:"APP_WRITE_TIMEOUT"`
	Tls               bool   `mapstructure:"tls" json:"tls" yaml:"tls" env:"APP_TLS"`
	CertFile          string `mapstructure:"cert_file" json:"cert_file" yaml:"cert_file" env:"APP_CERT_FILE"`
	KeyFile           string `mapstructure:"key_file" json:"key_file" yaml:"key_file" env:"APP_KEY_FILE"`
	ShutdownTimeout   int    `mapstructure:"shutdown_timeout" json:"shutdown_timeout" yaml:"shutdown_timeout" env:"APP_SHUTDOWN_TIMEOUT"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host" env:"MYSQL_HOST" `
	Port         int    `mapstructure:"port" json:"port" yaml:"port"  env:"MYSQL_PORT"`
	User         string `mapstructure:"user" json:"user" yaml:"user"  env:"MYSQL_USER"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"  env:"MYSQL_PASSWORD"`
	DbName       string `mapstructure:"db_name" json:"db_name" yaml:"db_name"  env:"MYSQL_DB_NAME"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"  env:"MYSQL_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"  env:"MYSQL_MAX_IDLE_CONNS"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"  env:"MYSQL_CONFIG"`        // 高级配置
	LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"  env:"MYSQL_LOG_MODE"` // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"  env:"MYSQL_LOG_ZAP"`     // 是否通过zap写入日志文件
	Enable       bool   `mapstructure:"enable" json:"enable" yaml:"enable" env:"MYSQL_ENABLE"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host" env:"REDIS_HOST"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port" env:"REDIS_PORT"`
	Db       int    `mapstructure:"db" json:"db" yaml:"db" env:"REDIS_DB"`
	Password string `mapstructure:"password" json:"password" yaml:"password" env:"REDIS_PASSWORD"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size" env:"REDIS_POLL_SIZE"`
	Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable" env:"REDIS_ENABLE"`
}

type Config struct {
	App    AppConfig           `mapstructure:"app" json:"app" yaml:"app" env:"IKUBEOPS"`
	Logger *logger.IkubeLogger `mapstructure:"logger" json:"logger" yaml:"logger" env:"IKUBEOPS"`
	Mysql  MysqlConfig         `mapstructure:"mysql" json:"mysql" yaml:"mysql" env:"IKUBEOPS"`
	Redis  RedisConfig         `mapstructure:"redis" json:"redis" yaml:"redis" env:"IKUBEOPS"`
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		HttpPort:          9900,
		HttpAddr:          "0.0.0.0",
		Language:          "zh",
		MaxHeaderSize:     1,
		ReadTimeout:       60,
		ReadHeaderTimeout: 60,
		WriteTimeout:      60,
		Tls:               false,
		KeyFile:           "",
		CertFile:          "",
		ShutdownTimeout:   60,
	}
}

func NewLoggerConfig() *logger.IkubeLogger {
	return &logger.IkubeLogger{
		Output:     "console",
		Format:     "console",
		Level:      "debug",
		MaxFile:    false,
		Dev:        true,
		FilePath:   "./logs",
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 100,
	}
}
func NewMysqlConfig() *MysqlConfig {
	return &MysqlConfig{
		Host:         "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "123456",
		DbName:       "test",
		MaxOpenConns: 100,
		MaxIdleConns: 20,
		Config:       "charset=utf8mb4&parseTime=True&loc=Local",
		LogMode:      "Info",
		LogZap:       false,
		Enable:       false,
	}
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Db:       0,
		Password: "12345678",
		PoolSize: 100,
		Enable:   false,
	}
}

func NewDefaultConfig() *Config {
	return &Config{
		App:    *NewAppConfig(),
		Logger: NewLoggerConfig(),
		Mysql:  *NewMysqlConfig(),
		Redis:  *NewRedisConfig(),
	}
}
