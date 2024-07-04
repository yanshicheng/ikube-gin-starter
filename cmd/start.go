package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"github.com/yanshicheng/ikube-gin-starter/pkg/config"
	"github.com/yanshicheng/ikube-gin-starter/pkg/logger"
	"github.com/yanshicheng/ikube-gin-starter/pkg/mysql"
	"github.com/yanshicheng/ikube-gin-starter/pkg/version"
	"log"
)

// 注册所有服务

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "start",
	Short: fmt.Sprintf("%s API服务", version.IkubeopsProjectName),
	Long:  fmt.Sprintf("%s API服务", version.IkubeopsProjectName),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 初始化全局变量
		err = config.InitIkubeConfig(confType, confFile, global.C)
		if err != nil {
			log.Printf("初始化配置文件失败: %s", err)
			return err
		}

		// 初始化日志
		global.L, err = logger.InitIkubeLogger(
			global.C.Logger.Output,
			global.C.Logger.Output,
			global.C.Logger.Level,
			global.C.Logger.MaxFile,
			global.C.Logger.Dev,
			global.C.Logger.FilePath,
			global.C.Logger.MaxSize,
			global.C.Logger.MaxAge,
			global.C.Logger.MaxBackups)
		if err != nil {
			log.Printf("初始化日志失败: %s", err)
		}
		global.LSys = global.L.Named("system")
		global.LSys.Info("日志初始化成功")

		// 初始化数据库
		if global.C.Mysql.Enable {
			global.DB, err = mysql.InitIkubeGorm(
				fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", global.C.Mysql.User, global.C.Mysql.Password, global.C.Mysql.Host, global.C.Mysql.Port, global.C.Mysql.DbName, global.C.Mysql.Opts),
				global.C.Mysql.MaxIdleConns,
				global.C.Mysql.MaxOpenConns,
				global.C.Mysql.LogToFile,
				global.C.Mysql.Level,
			)
			if err != nil {
				global.LSys.Error(fmt.Sprintf("初始化数据库失败: %s", err))
				return err
			}
		}
		return nil

	},
}

func init() {
	rootCommand.AddCommand(serviceCmd)
}
