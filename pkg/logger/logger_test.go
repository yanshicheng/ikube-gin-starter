package logger_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yanshicheng/ikube-gin-starter/pkg/logger"
)

func TestIkubeLogger_LoadLogger(t *testing.T) {
	// 设置测试用的日志配置
	log, err := logger.InitIkubeLogger("file", "console", "info", true, true, "./test_logs/", 10, 3, 7)

	// 加载日志记录器

	// 断言检查
	assert.NoError(t, err, "加载日志记录器应该没有错误")

	// 进行日志输出测试
	log.Debug("这是一个调试日志")
	log.Info("这是一个信息日志")
	log.Warn("这是一个警告日志")
	log.Error("这是一个错误日志")
	log.Named("控制器").Error("这是一个错误日志")
}

func TestIkubeLogger_GetWriteSyncer(t *testing.T) {
	// 设置测试用的日志配置
	log, err := logger.InitIkubeLogger("file", "json", "debug", false, true, "./test_logs/", 1, 3, 7)

	// 断言检查
	assert.NoError(t, err, "加载日志记录器应该没有错误")

	// 进行日志输出测试
	log.Debug("这是一个调试日志")
	largeData := generateLargeData()
	log.Info("这是一个信息日志")
	log.Info(largeData)
	log.Info(largeData)
	log.Info(largeData)
	log.Warn("这是一个警告日志")
	log.Named("控制器").Error("这是一个错误日志")
	log.Named("控制器").Named("xxxx").Error("这是一个错误日志")

	// 等待日志文件滚动完成（预留时间超过日志文件最大保留时间）
	time.Sleep(20 * time.Second)

	// 清理测试文件夹
	err = os.RemoveAll("./test_logs")
	assert.NoError(t, err, "清理日志文件夹应该没有错误")
}

func generateLargeData() string {
	// 生成大数据（约 0.8MB）
	const eightHundredKB = 800 * 1024 // 800KB 数据大小
	largeData := make([]byte, eightHundredKB)
	for i := 0; i < eightHundredKB; i++ {
		largeData[i] = byte('A' + (i % 26))
	}
	return string(largeData)
}
