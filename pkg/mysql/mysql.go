package mysql

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

// IkubeGorm 结构体用于管理 MySQL 连接
type IkubeGorm struct {
	db           *gorm.DB // GORM 数据库实例
	dsn          string   // 数据源名称
	level        string
	maxIdleConns int  // 连接池中空闲连接的最大数量
	maxOpenConns int  // 连接池中最大打开连接数
	logToFile    bool // 是否将日志写入文件
}

// InitIkubeGorm 初始化一个新的 IkubeGorm 实例
func InitIkubeGorm(dsn string, maxIdleConns, maxOpenConns int, logToFile bool, level string) (*IkubeGorm, error) {
	ikube := &IkubeGorm{
		dsn:          dsn,
		maxIdleConns: maxIdleConns,
		maxOpenConns: maxOpenConns,
		logToFile:    logToFile,
		level:        level,
	}

	// 初始化并加载 MySQL 连接
	if err := ikube.load(); err != nil {
		return nil, err
	}

	return ikube, nil
}

// load 初始化 MySQL 连接并设置数据库
func (ikube *IkubeGorm) load() error {
	// 配置 MySQL 连接参数
	mysqlConfig := mysql.Config{
		DSN:                       ikube.dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.New(mysqlConfig), ikube.gormConfig())
	if err != nil {
		zap.L().Error("MySQL 启动异常", zap.Error(err))
		return err
	}

	// 设置连接池参数
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(ikube.maxIdleConns)
	sqlDB.SetMaxOpenConns(ikube.maxOpenConns)

	ikube.db = db
	return nil
}

// gormConfig 是一个示例函数，用于返回 GORM 的配置
func (ikube *IkubeGorm) gormConfig() *gorm.Config {
	var gormLogger logger.Interface
	var loggerLevel logger.LogLevel
	switch strings.ToLower(ikube.level) {
	case "debug":
		loggerLevel = logger.Info
	case "info":
		loggerLevel = logger.Info
	case "warn":
		loggerLevel = logger.Warn
	case "error":
		loggerLevel = logger.Error
	case "silent":
		loggerLevel = logger.Silent
	default:
		loggerLevel = logger.Info
	}
	if ikube.logToFile {
		// 日志写入文件
		file, err := os.OpenFile("gorm.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("无法创建日志文件: %v", err)
			os.Exit(1)
		}
		gormLogger = logger.New(
			log.New(file, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // 慢速 SQL 阈值
				LogLevel:                  loggerLevel, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 错误
				ParameterizedQueries:      true,        // 参数化查询
				Colorful:                  false,       // 禁用颜色
			},
		)
	} else {
		// 日志输出到控制台
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // 慢速 SQL 阈值
				LogLevel:                  loggerLevel, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 错误
				ParameterizedQueries:      true,        // 参数化查询
				Colorful:                  true,        // 启用颜色
			},
		)
	}

	return &gorm.Config{
		Logger: gormLogger,
	}
}

// GetDb 返回 GORM 数据库实例
func (ikube *IkubeGorm) GetDb() *gorm.DB {
	return ikube.db
}

// Close 关闭数据库连接
func (ikube *IkubeGorm) Close() error {
	sqlDB, _ := ikube.db.DB()
	return sqlDB.Close()
}

// Ping 检查数据库连接是否存活
func (ikube *IkubeGorm) Ping() error {
	sqlDB, _ := ikube.db.DB()
	return sqlDB.Ping()
}
