package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// IkubeLogger 日志配置结构体
type IkubeLogger struct {
	Output     string `json:"output" yaml:"output" mapstructure:"output" env:"LOG_OUTPUT"`                      // 日志输出位置，支持 "file" 和 "console"
	Format     string `json:"format" yaml:"format" mapstructure:"format"  env:"LOG_FORMAT"`                     // 日志格式，支持 "json" 和 "console"
	Level      string `json:"level" yaml:"level" mapstructure:"level"  env:"LOG_LEVEL"`                         // 日志级别，支持 "debug", "info", "warn", "error"
	MaxFile    bool   `json:"max_file" yaml:"max_file" mapstructure:"max_file"  env:"LOG_MAX_FILE"`             // 日志是否拆分文件
	Dev        bool   `json:"dev" yaml:"dev" mapstructure:"dev"  env:"LOG_DEV"`                                 // 是否开启开发模式
	FilePath   string `json:"file_path" yaml:"file_path" mapstructure:"file_path"  env:"LOG_FILE_PATH"`         // 日志文件路径
	MaxSize    int    `json:"max_size" yaml:"max_size" mapstructure:"max_size"  env:"LOG_MAX_SIZE"`             // 日志文件大小限制，单位 MB
	MaxAge     int    `json:"max_age" yaml:"max_age" mapstructure:"max_age"  env:"LOG_MAX_AGE"`                 // 日志文件保留天数
	MaxBackups int    `json:"max_backups" yaml:"max_backups" mapstructure:"max_backups"  env:"LOG_MAX_BACKUPS"` // 日志文件保留数量
}

func InitIkubeLogger(output, format, level string, maxFile bool, dev bool, filePath string, maxSize, maxAge, maxBackups int) (*zap.Logger, error) {
	il := IkubeLogger{
		Output:     output,
		Format:     format,
		Level:      level,
		MaxFile:    maxFile,
		Dev:        dev,
		FilePath:   filePath,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
	}
	// 初始化 logger

	logger := il.encoderConfig()
	return logger, nil
}
func (i *IkubeLogger) encoderConfig() *zap.Logger {
	logger, err := i.loadLogger()
	if err != nil {
		panic(err)
	}
	return logger
}

// LoadLogger 初始化日志记录器
// LoadLogger 方法根据配置信息初始化和返回一个 Zap 日志记录器实例
func (l *IkubeLogger) loadLogger() (*zap.Logger, error) {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch l.Level {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	// 配置日志输出位置和格式
	var cores []zapcore.Core

	if l.MaxFile {
		// 创建文件日志的各个级别核心
		debugCore := zapcore.NewCore(l.getEncoder(), l.GetWriteSyncer(l.FilePath+"debug.log"), zap.NewAtomicLevelAt(zapcore.DebugLevel))
		infoCore := zapcore.NewCore(l.getEncoder(), l.GetWriteSyncer(l.FilePath+"info.log"), zap.NewAtomicLevelAt(zapcore.InfoLevel))
		warnCore := zapcore.NewCore(l.getEncoder(), l.GetWriteSyncer(l.FilePath+"warn.log"), zap.NewAtomicLevelAt(zapcore.WarnLevel))
		errorCore := zapcore.NewCore(l.getEncoder(), l.GetWriteSyncer(l.FilePath+"error.log"), zap.NewAtomicLevelAt(zapcore.ErrorLevel))

		// 控制台日志核心
		consoleCore := zapcore.NewCore(l.getEncoder(), zapcore.Lock(os.Stdout), atomicLevel)
		cores = append(cores, debugCore, infoCore, warnCore, errorCore, consoleCore)
	} else {
		// 创建文件日志的统一核心
		unifiedCore := zapcore.NewCore(l.getEncoder(), l.GetWriteSyncer(l.FilePath+"ikubeops.log"), atomicLevel)

		// 控制台日志核心
		if l.Output == "console" && l.Dev {
			consoleCore := zapcore.NewCore(l.getEncoder(), zapcore.Lock(os.Stdout), atomicLevel)
			cores = append(cores, unifiedCore, consoleCore)
		} else if l.Output == "console" {
			consoleCore := zapcore.NewCore(l.getEncoder(), zapcore.Lock(os.Stdout), atomicLevel)
			cores = append(cores, consoleCore)
		} else {
			cores = append(cores, unifiedCore)
		}
	}

	// 创建Logger实例
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())

	// 设置全局Logger实例，以便在其他包中可以直接使用 zap.L() 调用
	zap.ReplaceGlobals(logger)
	// 创建并存储 coreLogger 实例
	coreLoggerInstance := &coreLogger{
		logger:       newLogger(logger, ""),
		rootLogger:   logger,
		globalLogger: logger.WithOptions(),
		webLogger:    newGinLogger(logger, ""),
		atom:         atomicLevel,
	}
	storeLogger(coreLoggerInstance)

	return logger, nil
}

// LoadGormLogger 初始化并返回一个用于 GORM 的 zap 日志记录器
func (l *IkubeLogger) LoadGormLogger() (*zap.Logger, error) {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch l.Level {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	// 创建文件日志核心
	sqlCore := zapcore.NewCore(l.getEncoder(), l.GetWriteSyncer(l.FilePath+"sql.log"), atomicLevel)

	// 创建Logger实例
	logger := zap.New(sqlCore, zap.AddCaller())

	return logger, nil
}

// getEncoder 获取日志编码器
// getEncoder 方法根据 Dev 字段设置选择性返回开发环境或生产环境的日志编码器
func (l *IkubeLogger) getEncoder() zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	if l.Dev {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         // 时间格式为ISO8601
	encoderConfig.EncodeTime = customTimeEncoder                  // 时间格式为ISO8601
	encoderConfig.TimeKey = "time"                                // 日志时间字段名
	encoderConfig.MessageKey = "message"                          // 日志消息字段名
	encoderConfig.CallerKey = "caller"                            // 日志调用者字段名
	encoderConfig.LevelKey = "level"                              // 日志级别字段名
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       // 日志级别大写
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder // 记录持续时间
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       // 短路径调用者编码器

	// 根据 Format 字段来选择编码器
	switch l.Format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		// 默认使用 JSON 编码器
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// customTimeEncoder 自定义时间编码器，格式为 YYYY-MM-DD HH:mm:ss
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// GetWriteSyncer 方法根据文件路径参数创建并返回一个日志写入同步器
func (l *IkubeLogger) GetWriteSyncer(file string) zapcore.WriteSyncer {
	// 创建 lumberjack.Logger 实例，用于日志文件滚动
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,         // 日志文件名
		MaxSize:    l.MaxSize,    // 单个日志文件大小上限（MB）
		MaxBackups: l.MaxBackups, // 日志文件最多保留备份数量
		MaxAge:     l.MaxAge,     // 日志文件保留天数
		Compress:   true,         // 是否压缩旧日志文件
	}
	// 返回写入同步器，将 lumberjack.Logger 转换为 zapcore.WriteSyncer
	return zapcore.AddSync(lumberJackLogger)
}
