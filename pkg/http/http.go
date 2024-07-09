package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"github.com/yanshicheng/ikube-gin-starter/pkg/version"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type NamedServer struct {
	Name   string
	Server *http.Server
}

// IkubeopsServerManager 结构体管理多个HTTP服务器
type IkubeopsServerManager struct {
	servers []*NamedServer // 存储HTTP服务器的数组
	g       errgroup.Group // 等待所有服务器关闭的等待组

	// 通用参数
	MaxHeaderSize     int    `mapstructure:"max_header_size" json:"max_header_size" yaml:"max_header_size" env:"APP_MAX_HEADER_SIZE"`
	ReadTimeout       int    `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout" env:"APP_READ_TIMEOUT"`
	ReadHeaderTimeout int    `mapstructure:"read_header_timeout" json:"read_header_timeout" yaml:"read_header_timeout" env:"APP_READ_HEADER_TIMEOUT"`
	WriteTimeout      int    `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout" env:"APP_WRITE_TIMEOUT"`
	Tls               bool   `mapstructure:"tls" json:"tls" yaml:"tls" env:"APP_TLS"`
	CertFile          string `mapstructure:"cert_file" json:"cert_file" yaml:"cert_file" env:"APP_CERT_FILE"`
	KeyFile           string `mapstructure:"key_file" json:"key_file" yaml:"key_file" env:"APP_KEY_FILE"`
	ShutdownTimeout   int    `mapstructure:"shutdown_timeout" json:"shutdown_timeout" yaml:"shutdown_timeout" env:"APP_SHUTDOWN_TIMEOUT"`
}

func NewIkubeHttpManager(maxHeaderSize, readTimeout, readHeaderTimeout, writeTimeout, shutdownTimeout int, tls bool, certFile, keyFile string) *IkubeopsServerManager {
	return &IkubeopsServerManager{
		MaxHeaderSize:     maxHeaderSize,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		ShutdownTimeout:   shutdownTimeout,
		Tls:               tls,
		CertFile:          certFile,
		KeyFile:           keyFile,
	}
}

// AddServer 添加一个HTTP服务器到管理器中
func (ism *IkubeopsServerManager) AddServer(addr, name string, router *gin.Engine) {
	server := &http.Server{
		Addr:              addr,                                               // 服务器地址
		Handler:           router,                                             // 路由处理器
		ReadTimeout:       time.Duration(ism.ReadTimeout) * time.Second,       // 读取超时时间
		WriteTimeout:      time.Duration(ism.WriteTimeout) * time.Second,      // 写入超时时间
		ReadHeaderTimeout: time.Duration(ism.ReadHeaderTimeout) * time.Second, // 读取请求头超时时间
		MaxHeaderBytes:    ism.MaxHeaderSize * 1024 * 1024,                    // 最大请求头字节限制
	}
	ism.servers = append(ism.servers, &NamedServer{
		Server: server,
		Name:   name,
	})
}

// Run 启动所有添加的HTTP服务器
func (ism *IkubeopsServerManager) Run() {
	for _, nameServer := range ism.servers {
		if ism.Tls && nameServer.Name == "business" {
			ism.g.Go(func() error {
				return nameServer.Server.ListenAndServeTLS(ism.CertFile, ism.KeyFile)
			})
		} else {
			ism.g.Go(func() error {
				return nameServer.Server.ListenAndServe()
			})
		}
		if nameServer.Name == "business" {
			fmt.Printf(`
欢迎使用: %s
当前版本: %s
配置文件: %s
演示地址: www.ikubeops.com
代码地址: %s
运行地址: %s

`, version.IkubeopsProjectName, version.ShortTagVersion(), version.GetConfig(), version.IkubeopsUrl, version.GetWebUrl(nameServer.Server.Addr, ism.Tls))
		}
	}
	// 优雅关闭监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号接收
	sig := <-quit

	global.LSys.Info(fmt.Sprintf("接收到退出信号: %s", sig))
	now := time.Now()
	// 设置一个超时时间的 ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ism.ShutdownTimeout)*time.Second)
	defer cancel()
	for _, nameServer := range ism.servers {
		ism.g.Go(func() error {
			return nameServer.Server.Shutdown(ctx)
		})
	}
	// 等待所有HTTP服务器关闭
	if err := ism.g.Wait(); err != nil {
		// 打印类型 err 类型
		// 检查错误类型是否为 *net.OpError
		if err.Error() == "http: Server closed" {
		} else {
			global.LSys.Error(fmt.Sprintf("关闭服务器异常: %s", err))
		}
	}
	// 关闭数据库链接
	err := global.DB.Close()
	if err != nil {
		global.LSys.Error(fmt.Sprintf("关闭数据库链接异常: %s", err))
	}
	// 关闭redis
	err = global.RDB.Close()
	if err != nil {
		global.LSys.Error(fmt.Sprintf("关闭redis链接异常: %s", err))
	}
	global.LSys.Info(fmt.Sprintf("退出耗时: %s", time.Since(now)))
}
