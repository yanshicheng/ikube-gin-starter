package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	_ "github.com/yanshicheng/ikube-gin-starter/apps/all"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"github.com/yanshicheng/ikube-gin-starter/pkg/config"
	"github.com/yanshicheng/ikube-gin-starter/pkg/logger"
	"github.com/yanshicheng/ikube-gin-starter/pkg/mysql"
	"log"
)

var (
	db      string
	migrate bool
)

var dbCommand = &cobra.Command{
	Use:   "db",
	Short: "db console",
	Long:  "db console",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		// 初始化全局变量
		err = config.InitIkubeConfig(confFile, confType, global.C)
		if err != nil {
			log.Printf("初始化配置文件失败: %s", err)
			return err
		}
		// 检查 HttpPort 和 HealthPort 是否一致，一致则报错
		if global.C.App.HttpPort == global.C.App.HealthPort {
			log.Printf("HttpPort 和 HealthPort 不能一致")
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
		global.LSys.Info("日志初始化成功!")

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
			} else {
				if global.DB.Ping() != nil {
					global.LSys.Error(fmt.Sprintf("Mysql 数据库连接失败: %s", err))
					return err
				}
				global.LSys.Info("Mysql 数据库初始化成功!")
			}
		} else {
			global.LSys.Error(fmt.Sprintf("数据迁移失败，未启用mysql配置，如需迁移数据库，请先启用mysql配置"))
			return nil
		}

		global.LSys.Info("开始数据迁移...")
		defer func(DB *mysql.IkubeGorm) {
			err := DB.Close()
			if err != nil {
				global.LSys.Error(fmt.Sprintf("关闭数据库失败: %s", err))
				return
			}
		}(global.DB)
		if global.DB != nil && len(global.M) > 0 {
			db := global.DB.GetDb()
			if err := db.AutoMigrate(global.M...); err != nil {
				global.LSys.Error(fmt.Sprintf("数据库迁移失败: %s\n", err.Error()))
			} else {
				global.LSys.Info("数据库迁移成功")
				return nil
			}
		} else {
			global.LSys.Info("未检测到 model")
		}
		return nil
	},
}

func init() {
	rootCommand.AddCommand(dbCommand)
	dbCommand.Flags().StringVarP(&db, "database", "d", "default", "database")
	dbCommand.Flags().BoolVarP(&migrate, "migrate", "m", false, "force syncdb")
}
