package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yanshicheng/ikube-gin-starter/common/middleware"
	docs "github.com/yanshicheng/ikube-gin-starter/docs" // 千万不要忘了导入把你上一步生成的docs
	"github.com/yanshicheng/ikube-gin-starter/global"
	"net/http"
)

// 业务路由
func BusinessRouter(ginApps map[string]GinService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 全局中间件配置
	router.Use(middleware.GinRecovery(true), middleware.GinLogger(), middleware.Cors())
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "10404",
			"data":    "",
			"message": "路由不存在",
		})
	})
	// 开放的路由组配置
	PublicRouterGroup := router.Group("")
	{
		registerSwagger(PublicRouterGroup)
	}

	// 鉴权路由
	AuthRouterGroup := router.Group("")
	// 鉴权中间件配置
	//AuthRouterGroup.Use()
	{
	}
	for _, ginApp := range ginApps {
		ginApp.PublicRegistry(PublicRouterGroup)
		ginApp.AuthRegistry(AuthRouterGroup)
	}
	return router
}

// 健康检查路由
func HealthRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "10404",
			"data":    "",
			"message": "路由不存在",
		})
	})
	router.GET("/healthy", healthStatus)
	return router
}

func healthStatus(c *gin.Context) {
	// Mysql 检查
	if global.C.Mysql.Enable {
		if err := global.DB.Ping(); err != nil {
			global.LSys.Error(fmt.Sprintf("Mysql 健康检查失败: %s", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"data":    "Mysql error",
				"message": err.Error(),
			})
			return
		}
	}
	// Redis 检查
	if global.C.Redis.Enable {
		if err := global.RDB.Ping(); err != nil {
			global.LSys.Error(fmt.Sprintf("Redis 健康检查失败: %s", err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"data":    "Redis error",
				"message": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    "ok",
		"message": "ok",
	})
	return
}

func registerSwagger(r gin.IRouter) {
	// API文档访问地址: http://host/swagger/index.html
	// 注解定义可参考 https://github.com/swaggo/swag#declarative-comments-format
	// 样例 https://github.com/swaggo/swag/blob/master/example/basic/api/api.go
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "管理后台接口"
	docs.SwaggerInfo.Description = "实现一个管理系统的后端API服务"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
