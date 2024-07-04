package config_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/yanshicheng/ikube-gin-starter/pkg/config"
	"github.com/yanshicheng/ikube-gin-starter/pkg/types"
	"os"
	"testing"
)

func setupEnvironment(t *testing.T) {
	t.Helper() // 标记为辅助函数
	err := os.Setenv("IKUBEOPS_APP_HTTP_ADDR", "172.16.1.100")
	assert.NoError(t, err, "Setting IKUBEOPS_APP_HTTP_ADDR should not produce an error")
	err = os.Setenv("IKUBEOPS_APP_HTTP_PORT", "8080")
	assert.NoError(t, err, "Setting IKUBEOPS_APP_HTTP_PORT should not produce an error")
}

func cleanupEnvironment(t *testing.T) {
	t.Helper() // 标记为辅助函数
	err := os.Unsetenv("IKUBEOPS_APP_HTTP_PORT")
	assert.NoError(t, err, "Unsetting IKUBEOPS_APP_HTTP_PORT should not produce an error")
	err = os.Unsetenv("IKUBEOPS_APP_HTTP_ADDR")
	assert.NoError(t, err, "Unsetting IKUBEOPS_APP_HTTP_ADDR should not produce an error")
}

func TestLoadFileConfig(t *testing.T) {
	setupEnvironment(t)
	t.Cleanup(func() { cleanupEnvironment(t) })

	// 创建临时文件
	configFile := "./config_test.yaml"
	testConfigContent := []byte("app:\n  http_port: 9999\n  http_addr: \"0.0.0.0\"\n  tls: true")
	err := os.WriteFile(configFile, testConfigContent, 0644)
	assert.NoError(t, err, "Writing to config file should not produce an error")
	t.Cleanup(func() { os.Remove(configFile) }) // 清理文件

	destConfig := types.NewDefaultConfig()
	err = config.InitIkubeConfig(configFile, "IKUBEOPS_", destConfig)

	assert.NoError(t, err, "Loading file config should not produce an error")

	// 校验配置是否按预期加载
	assert.Equal(t, 8080, destConfig.App.HttpPort, "Loaded HttpPort should match the environment setting of 8080")
	assert.Equal(t, "172.16.1.100", destConfig.App.HttpAddr, "Loaded HttpAddr should match the environment setting of 172.16.1.100")
	assert.Equal(t, true, destConfig.App.Tls, "Loaded Tls should match the file setting of true")
}
