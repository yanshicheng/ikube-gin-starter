package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync/atomic"
	"unsafe"
)

var (
	_log unsafe.Pointer // 指向 coreLogger 的指针。通过 atomic.LoadPointer 访问。
)

type LogOption = zap.Option

// coreLogger 是日志核心结构体，包含了多个日志记录器和配置信息。
type coreLogger struct {
	logger       *Logger         // 基础日志记录器
	rootLogger   *zap.Logger     // 没有任何配置选项的根日志记录器
	webLogger    *Logger         // 用于 Web 日志记录的日志记录器
	globalLogger *zap.Logger     // 全局日志记录器
	atom         zap.AtomicLevel // 动态日志级别设置
}

// Logger 是包装了 zap.Logger 和 zap.SugaredLogger 的日志结构体。
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// storeLogger 存储日志记录器实例到 _log 中。
func storeLogger(l *coreLogger) {
	if old := loadLogger(); old != nil {
		old.rootLogger.Sync() // 同步旧的根日志记录器，确保日志被写入文件。
	}
	atomic.StorePointer(&_log, unsafe.Pointer(l))
}

// newLogger 创建一个新的 Logger 实例。
func newLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.
		WithOptions().
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

// newGinLogger 创建一个新的用于 Gin 框架的 Logger 实例。
func newGinLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.
		WithOptions().
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

// NewLogger 初始化日志记录器，设置全局日志记录器和 webLogger。
func NewLogger(e *IkubeLogger) error {
	atom := zap.NewAtomicLevel()           // 创建一个新的原子级别控制器
	logger, webLogger := e.EncoderConfig() // 获取编码器配置信息
	storeLogger(&coreLogger{
		rootLogger:   logger,
		logger:       newLogger(logger, ""),
		globalLogger: logger.WithOptions(),
		webLogger:    newGinLogger(webLogger, ""),
		atom:         atom,
	})
	return nil
}

// Named 返回一个添加了新路径段的日志记录器。
func (l *Logger) Named(name string) *Logger {
	logger := l.logger.Named(name)
	return &Logger{logger, logger.Sugar()}
}

// SetLevel 动态设置日志记录器的日志级别。
func (l *Logger) SetLevel(level string) {
	var zapLevel zap.AtomicLevel
	zapLevel.UnmarshalText([]byte(level))
	l.logger.Core().Enabled(zapLevel.Level())
}

// Option 是用于配置 zap.Config 的函数类型。
type Option func(*zap.Config)

// WithCaller 在日志输出中启用调用者字段。
func WithCaller(caller bool) Option {
	return func(config *zap.Config) {
		config.Development = !caller
		config.DisableCaller = !caller
	}
}

// Print 使用 fmt.Sprint 构造并记录一条消息。
func (l *Logger) Print(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Println 使用 fmt.Sprint 构造并记录一条消息。
func (l *Logger) Println(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debug 使用 fmt.Sprint 构造并记录一条调试级别的消息。
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Info 使用 fmt.Sprint 构造并记录一条信息级别的消息。
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Warn 使用 fmt.Sprint 构造并记录一条警告级别的消息。
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Error 使用 fmt.Sprint 构造并记录一条错误级别的消息。
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal 使用 fmt.Sprint 构造并记录一条致命错误级别的消息，然后调用 os.Exit(1)。
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Panic 使用 fmt.Sprint 构造并记录一条消息，然后 panic。
func (l *Logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

// DPanic 使用 fmt.Sprint 构造并记录一条消息。在开发模式下，日志记录器会 panic。
func (l *Logger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}

// IsDebug 检查日志记录器是否启用了调试级别。
func (l *Logger) IsDebug() bool {
	return l.logger.Check(zapcore.DebugLevel, "") != nil
}

// Printf 使用 fmt.Sprintf 记录一个格式化的消息。
func (l *Logger) Printf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Debugf 使用 fmt.Sprintf 记录一个格式化的调试级别的消息。
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Infof 使用 fmt.Sprintf 记录一个格式化的信息级别的消息。
func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warnf 使用 fmt.Sprintf 记录一个格式化的警告级别的消息。
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Errorf 使用 fmt.Sprintf 记录一个格式化的错误级别的消息。
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// Fatalf 使用 fmt.Sprintf 记录一个格式化的致命错误级别的消息，然后调用 os.Exit(1)。
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// Panicf 使用 fmt.Sprintf 记录一个格式化的消息，然后 panic。
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugar.Panicf(format, args...)
}

// DPanicf 使用 fmt.Sprintf 记录一个格式化的消息。在开发模式下，日志记录器会 panic。
func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.sugar.DPanicf(format, args...)
}

// Debugw 记录一个带有额外上下文的调试级别的消息。
func (l *Logger) Debugw(msg string, fields ...Field) {
	l.sugar.Debugw(msg, transfer(fields)...)
}

// Infow 记录一个带有额外上下文的信息级别的消息。
func (l *Logger) Infow(msg string, fields ...Field) {
	l.sugar.Infow(msg, transfer(fields)...)
}

// Warnw 记录一个带有额外上下文的警告级别的消息。
func (l *Logger) Warnw(msg string, fields ...Field) {
	l.sugar.Warnw(msg, transfer(fields)...)
}

// Errorw 记录一个带有额外上下文的错误级别的消息。
func (l *Logger) Errorw(msg string, fields ...Field) {
	l.sugar.Errorw(msg, transfer(fields)...)
}

// Fatalw 记录一个带有额外上下文的致命错误级别的消息，然后调用 os.Exit(1)。
func (l *Logger) Fatalw(msg string, fields ...Field) {
	l.sugar.Fatalw(msg, transfer(fields)...)
}

// Panicw 记录一个带有额外上下文的消息，然后 panic。
func (l *Logger) Panicw(msg string, fields ...Field) {
	l.sugar.Panicw(msg, transfer(fields)...)
}

// DPanicw 记录一个带有额外上下文的消息。在开发模式下，日志记录器会 panic。
func (l *Logger) DPanicw(msg string, fields ...Field) {
	l.sugar.DPanicw(msg, transfer(fields)...)
}

// Field 是键值对，用于传递额外的上下文信息。
type Field struct {
	Key   string      // 键
	Value interface{} // 值
}

// transfer 将 Field 转换为 zap.Any 类型的切片，用于日志记录。
func transfer(m []Field) (ma []interface{}) {
	for i := range m {
		ma = append(ma, zap.Any(m[i].Key, m[i].Value))
	}

	return
}

// globalLogger 返回全局日志记录器。
func globalLogger() *zap.Logger {
	return loadLogger().globalLogger
}

// loadLogger 加载当前的日志记录器实例。
func loadLogger() *coreLogger {
	p := atomic.LoadPointer(&_log)
	return (*coreLogger)(p)
}

// SetLevel 设置全局日志级别。
func SetLevel(lv Level) {
	loadLogger().atom.SetLevel(lv.zapLevel())
}

// L 返回基础日志记录器。
func L() *Logger {
	return loadLogger().logger
}

// W 返回 Web 日志记录器。
func W() *Logger {
	return loadLogger().webLogger
}

// Recover 停止一个 panic 的 goroutine，并记录一个错误级别的消息。
func (l *Logger) Recover(msg string) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("%s. Recovering, but please report this.", msg)
		globalLogger().WithOptions().
			Error(msg, zap.Any("panic", r), zap.Stack("stack"))
	}
}
