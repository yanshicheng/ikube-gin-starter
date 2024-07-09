package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// Cors 函数定义了一个用于处理跨源资源共享（CORS）的中间件。
func Cors() gin.HandlerFunc {
	// 返回 cors.New 函数创建的一个新的 CORS 中间件处理器。
	// cors.New 接受一个 cors.Config 配置对象，用于定义 CORS 策略。
	cors.Default()
	return cors.New(cors.Config{
		// AllowAllOrigins 设置为 true，允许所有域名进行跨域请求。
		//AllowOrigins: []string{"http://127.0.0.1:9999"},
		AllowAllOrigins: true,

		// AllowMethods 定义了允许的 HTTP 方法，这里包括了常见的 CRUD 操作和 OPTIONS 方法。
		AllowMethods: []string{
			"GET",     // 允许 GET 请求
			"POST",    // 允许 POST 请求
			"PUT",     // 允许 PUT 请求
			"PATCH",   // 允许 PATCH 请求
			"DELETE",  // 允许 DELETE 请求
			"HEAD",    // 允许 HEAD 请求
			"OPTIONS", // 允许 OPTIONS 请求，用于预请求
		},

		// AllowHeaders 定义了允许在请求中发送的自定义头信息。
		AllowHeaders: []string{
			"Origin",         // 允许 Origin 头
			"Content-Length", // 允许 Content-Length 头
			"Content-Type",   // 允许 Content-Type 头
			"Authorization",  // 允许 Authorization 头，通常用于认证信息
		},

		// AllowCredentials 设置为 false，表示服务器不会将响应的 Cookies 包含在响应中。
		// 这通常用于公共 API，不涉及用户认证信息的传递。
		AllowCredentials: false,

		// MaxAge 定义了预请求（OPTIONS 请求）的有效时间，这里设置为 12 小时。
		MaxAge: 12 * time.Hour,

		// ExposeHeaders 定义了客户端可以访问的响应头，这里只允许 Content-Length。
		//ExposeHeaders: []string{
		//	"Content-Length", // 允许客户端访问响应的 Content-Length 头
		//},
	})
}
