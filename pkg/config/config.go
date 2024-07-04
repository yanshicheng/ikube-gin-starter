package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/spf13/viper"
	"github.com/yanshicheng/ikube-gin-starter/pkg/types"
	"os"
	"strings"
)

type IkubeConfig struct {
	// 配置文件路径
	ConfigFile string
	EnvPrefix  string
	DestStruct *types.Config
}

func InitIkubeConfig(configFile, EnvPrefix string, destStruct *types.Config) error {
	ic := IkubeConfig{
		ConfigFile: configFile,
		EnvPrefix:  EnvPrefix,
		DestStruct: destStruct,
	}
	return ic.loadFileConfig()
}

func (i *IkubeConfig) loadFileConfig() error {
	// 检查文件是否存在
	if err := i.checkFileConfig(); err != nil {
		return err
	}
	viper.SetConfigFile(i.ConfigFile)
	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {             // 读取配置信息失败
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(&i.DestStruct); err != nil {
		return fmt.Errorf("unmarshal conf failed, err:%s \n", err)
	}
	if err := i.loadConfigFromEnv(); err != nil {
		return fmt.Errorf("load config from env failed, err:%s \n", err)
	}
	//// 监控配置文件变化
	//viper.WatchConfig()
	//// 注意！！！配置文件发生变化后要同步到全局变量Conf
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println("配置文件被修改啦,正在重载...")
	//	if err := viper.Unmarshal(&global.C); err != nil {
	//		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	//	}
	//})
	return nil
}

// 检查文件是否存在，并且是可读，还有是 yaml 或者 toml 文件
func (i *IkubeConfig) checkFileConfig() error {
	// 检查文件是否存在
	if _, err := os.Stat(i.ConfigFile); err != nil {
		return fmt.Errorf("config file does not exist: %s", i.ConfigFile)
	}
	// 检查文件是否可读
	if _, err := os.Open(i.ConfigFile); err != nil {
		return fmt.Errorf("config file is not readable: %s", i.ConfigFile)
	}
	// 检查文件是否是 yaml 或者 toml 文件
	// 转换为小写以忽略大小写差异
	lowerFilename := strings.ToLower(i.ConfigFile)

	if strings.HasSuffix(lowerFilename, ".yaml") ||
		strings.HasSuffix(lowerFilename, ".yml") ||
		strings.HasSuffix(lowerFilename, ".toml") {
		return nil
	}
	return errors.New("file is not YAML or TOML")
}

func (i *IkubeConfig) loadConfigFromEnv() error {
	opts := env.Options{
		Prefix: i.EnvPrefix, // 只加载以 T_ 开头的环境变量。
	}
	err := env.ParseWithOptions(i.DestStruct, opts)
	if err != nil {
		return fmt.Errorf("failed to load config from environment variables: %w", err)
	}
	return nil
}

//func LoadGlobalConfig(configType, configFile string) error {
//	// 配置加载
//	version.IkubeopsConfigFile = configFile
//	version.IkubeopsConfigType = configType
//	switch configType {
//	case "file":
//		return LoadFileConfig(configFile)
//	default:
//		return errors.New("unknown config type")
//	}
//}
