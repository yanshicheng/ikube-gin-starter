package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanshicheng/ikube-gin-starter/pkg/version"
)

var (
	// pusher service config option
	confType string
	confFile string
	confETCD string
)

var vers bool

// RootCmd 表示在没有任何子命令调用时的基本命令
var rootCommand = &cobra.Command{
	Use:   version.IkubeopsProjectName,
	Short: fmt.Sprintf("%s 官网地址: www.ikubeops.com", version.IkubeopsProjectName),
	Long:  fmt.Sprintf("%s 官网地址: www.ikubeops.com", version.IkubeopsProjectName),
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullTagVersion())
			return nil
		}
		return cmd.Help()
	},
}

// Execute 将所有子命令添加到根命令并设置标志位。这由 main.main() 调用。只需要对 rootCmd 进行一次处理。
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	rootCommand.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "配置文件类型 [file/etcd]")
	rootCommand.PersistentFlags().StringVarP(&confFile, "config-file", "f", "config/config.yaml", "配置文件路径")
	rootCommand.PersistentFlags().StringVarP(&confETCD, "config-etcd", "e", "127.0.0.1:2379", "etcd 配置")
	rootCommand.PersistentFlags().BoolVarP(&vers, "version", "v", false, "app 版本信息")
}
