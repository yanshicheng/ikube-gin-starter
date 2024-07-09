package cmd

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/cobra"
	_ "github.com/yanshicheng/ikube-gin-starter/apps/all"
	"github.com/yanshicheng/ikube-gin-starter/common/validator"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"github.com/yanshicheng/ikube-gin-starter/pkg/config"
	"github.com/yanshicheng/ikube-gin-starter/pkg/http"
	"github.com/yanshicheng/ikube-gin-starter/pkg/logger"
	"github.com/yanshicheng/ikube-gin-starter/pkg/mysql"
	"github.com/yanshicheng/ikube-gin-starter/pkg/redis"
	"github.com/yanshicheng/ikube-gin-starter/pkg/version"
	"github.com/yanshicheng/ikube-gin-starter/router"
	"log"
)

// 注册所有服务

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "start",
	Short: fmt.Sprintf("%s API服务", version.IkubeopsProjectName),
	Long:  fmt.Sprintf("%s API服务", version.IkubeopsProjectName),
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
		}

		// 初始化 redis
		if global.C.Redis.Enable {
			global.RDB, err = redis.InitIkubeRedis(
				fmt.Sprintf("%s:%d", global.C.Redis.Host, global.C.Redis.Port),
				global.C.Redis.Password,
				global.C.Redis.Db,
				global.C.Redis.PoolSize,
			)
			if err != nil {
				global.LSys.Error(fmt.Sprintf("Reids 初始化数据库失败: %s", err))
				return err
			} else {
				if global.RDB.Ping() != nil {
					global.LSys.Error(fmt.Sprintf("Redis 数据库连接失败: %s", err))
				}
				global.LSys.Info("Redis 数据库初始化成功!")
			}
		}

		// 初始化Gin框架翻译器
		var uni *ut.UniversalTranslator
		if global.IkubeopsTrans, uni, err = validator.InitTrans(global.C.App.Language); err != nil {
			global.LSys.Error(fmt.Sprintf("初始化Gin框架翻译器失败: %s", err))
			return err
		} else {

		}
		global.LSys.Info("Gin框架翻译器初始化成功!")

		// 引入自定义验证器
		if err = validator.RegisterValidatorsAndTranslations(validator.ValidatorSlice, uni); err != nil {
			global.LSys.Error(fmt.Sprintf("注册验证器失败: %s", err))
			return err
		}
		global.LSys.Info("验证器加载成功!")
		// 启动服务
		// 获取gin app 实例
		businessRouter := router.InitGin()
		// 初始化路由
		healthRouter := router.HealthRouter()

		// 初始化http管理器
		serverManager := http.NewIkubeHttpManager(
			global.C.App.MaxHeaderSize,
			global.C.App.ReadTimeout,
			global.C.App.ReadHeaderTimeout,
			global.C.App.WriteTimeout,
			global.C.App.ShutdownTimeout,
			global.C.App.Tls,
			global.C.App.CertFile,
			global.C.App.KeyFile)
		serverManager.AddServer(fmt.Sprintf("%s:%d", global.C.App.HttpAddr, global.C.App.HttpPort), "business", businessRouter)
		serverManager.AddServer(fmt.Sprintf("%s:%d", global.C.App.HttpAddr, global.C.App.HealthPort), "healthy", healthRouter)
		serverManager.Run()
		return nil

	},
}

func init() {
	rootCommand.AddCommand(serviceCmd)
}
