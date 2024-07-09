package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync/atomic"
	"unsafe"
)

var (
	_log unsafe.Pointer // Pointer to coreLogger, accessed via atomic.LoadPointer
)

type LogOption = zap.Option

// coreLogger is the core logging structure, containing multiple loggers and configuration information.
type coreLogger struct {
	logger       *Logger         // Basic logger
	rootLogger   *zap.Logger     // Root logger without any configuration options
	webLogger    *Logger         // Logger for web logging
	globalLogger *zap.Logger     // Global logger
	atom         zap.AtomicLevel // Dynamic log level setting
}

// Logger wraps zap.Logger and zap.SugaredLogger.
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// storeLogger stores the logger instance in _log.
func storeLogger(l *coreLogger) {
	if old := loadLogger(); old != nil {
		old.rootLogger.Sync() // Sync the old root logger to ensure logs are written to file.
	}
	atomic.StorePointer(&_log, unsafe.Pointer(l))
}

// newLogger creates a new Logger instance.
func newLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.
		WithOptions().
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

// newGinLogger creates a new Logger instance for the Gin framework.
func newGinLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.
		WithOptions().
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

// NewLogger initializes the logger, setting the global logger and webLogger.
func NewLogger(e *IkubeLogger) error {
	atom := zap.NewAtomicLevel() // Create a new atomic level controller
	logger := e.encoderConfig()  // Get the encoder configuration

	coreLoggerInstance := &coreLogger{
		rootLogger:   logger,
		logger:       newLogger(logger, ""),
		globalLogger: logger.WithOptions(),
		webLogger:    newGinLogger(logger, ""),
		atom:         atom,
	}

	storeLogger(coreLoggerInstance)
	return nil
}

// Named returns a logger with a new path segment.
func (l *Logger) Named(name string) *Logger {
	logger := l.logger.Named(name)
	return &Logger{logger, logger.Sugar()}
}

// SetLevel dynamically sets the logger's log level.
func (l *Logger) SetLevel(level string) {
	var zapLevel zap.AtomicLevel
	zapLevel.UnmarshalText([]byte(level))
	l.logger.Core().Enabled(zapLevel.Level())
}

// Option is a function type for configuring zap.Config.
type Option func(*zap.Config)

// WithCaller enables the caller field in log output.
func WithCaller(caller bool) Option {
	return func(config *zap.Config) {
		config.Development = !caller
		config.DisableCaller = !caller
	}
}

// Print logs a message using fmt.Sprint.
func (l *Logger) Print(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Println logs a message using fmt.Sprint.
func (l *Logger) Println(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debug logs a debug-level message using fmt.Sprint.
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Info logs an info-level message using fmt.Sprint.
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Warn logs a warning-level message using fmt.Sprint.
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Error logs an error-level message using fmt.Sprint.
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal logs a fatal error-level message using fmt.Sprint, then calls os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Panic logs a message using fmt.Sprint, then panics.
func (l *Logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

// DPanic logs a message using fmt.Sprint. In development mode, the logger panics.
func (l *Logger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}

// IsDebug checks if the logger is enabled for the debug level.
func (l *Logger) IsDebug() bool {
	return l.logger.Check(zapcore.DebugLevel, "") != nil
}

// Printf logs a formatted message using fmt.Sprintf.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Debugf logs a formatted debug-level message using fmt.Sprintf.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Infof logs a formatted info-level message using fmt.Sprintf.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warnf logs a formatted warning-level message using fmt.Sprintf.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Errorf logs a formatted error-level message using fmt.Sprintf.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// Fatalf logs a formatted fatal error-level message using fmt.Sprintf, then calls os.Exit(1).
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// Panicf logs a formatted message using fmt.Sprintf, then panics.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugar.Panicf(format, args...)
}

// DPanicf logs a formatted message using fmt.Sprintf. In development mode, the logger panics.
func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.sugar.DPanicf(format, args...)
}

// Debugw logs a debug-level message with additional context.
func (l *Logger) Debugw(msg string, fields ...Field) {
	l.sugar.Debugw(msg, transfer(fields)...)
}

// Infow logs an info-level message with additional context.
func (l *Logger) Infow(msg string, fields ...Field) {
	l.sugar.Infow(msg, transfer(fields)...)
}

// Warnw logs a warning-level message with additional context.
func (l *Logger) Warnw(msg string, fields ...Field) {
	l.sugar.Warnw(msg, transfer(fields)...)
}

// Errorw logs an error-level message with additional context.
func (l *Logger) Errorw(msg string, fields ...Field) {
	l.sugar.Errorw(msg, transfer(fields)...)
}

// Fatalw logs a fatal error-level message with additional context, then calls os.Exit(1).
func (l *Logger) Fatalw(msg string, fields ...Field) {
	l.sugar.Fatalw(msg, transfer(fields)...)
}

// Panicw logs a message with additional context, then panics.
func (l *Logger) Panicw(msg string, fields ...Field) {
	l.sugar.Panicw(msg, transfer(fields)...)
}

// DPanicw logs a message with additional context. In development mode, the logger panics.
func (l *Logger) DPanicw(msg string, fields ...Field) {
	l.sugar.DPanicw(msg, transfer(fields)...)
}

// Field is a key-value pair for passing additional context information.
type Field struct {
	Key   string      // Key
	Value interface{} // Value
}

// transfer converts Field to a slice of zap.Any for logging.
func transfer(m []Field) (ma []interface{}) {
	for i := range m {
		ma = append(ma, zap.Any(m[i].Key, m[i].Value))
	}
	return
}

// globalLogger returns the global logger.
func globalLogger() *zap.Logger {
	cl := loadLogger()
	if cl == nil {
		panic("global logger is not initialized")
	}
	return cl.globalLogger
}

// loadLogger loads the current logger instance.
func loadLogger() *coreLogger {
	p := atomic.LoadPointer(&_log)
	if p == nil {
		return nil
	}
	return (*coreLogger)(p)
}

// SetLevel sets the global log level.
func SetLevel(lv Level) {
	cl := loadLogger()
	if cl == nil {
		panic("logger is not initialized")
	}
	cl.atom.SetLevel(lv.zapLevel())
}

// L returns the basic logger.
func L() *Logger {
	cl := loadLogger()
	if cl == nil {
		panic("logger is not initialized")
	}
	return cl.logger
}

// W returns the web logger.
func W() *Logger {
	cl := loadLogger()
	if cl == nil {
		panic("logger is not initialized")
	}
	return cl.webLogger
}

// Recover stops a goroutine panic and logs an error-level message.
func (l *Logger) Recover(msg string) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("%s. Recovering, but please report this.", msg)
		globalLogger().WithOptions().
			Error(msg, zap.Any("panic", r), zap.Stack("stack"))
	}
}
