package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-starter/common/response"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"go.uber.org/zap"
)

// Recover 是一个全局的恢复中间件，用于捕获 panic，并返回友好的错误响应。
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// 恢复 panic，并返回友好的错误响应
			c.JSON(http.StatusOK, gin.H{
				"code": "1",
				"msg":  errorToString(r),
				"data": nil,
			})
			c.Abort() // 中止请求处理
		}
	}()
	c.Next() // 继续处理请求
}

// errorToString 将 panic 中的错误转换为字符串
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}

// GinRecovery 是一个恢复中间件，用于捕获可能出现的 panic，并使用 zap 记录相关日志。
// 如果设置了 stack 为 true，则记录完整的堆栈信息。
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是断开的连接错误，这种情况不需要 panic 堆栈信息。
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取当前请求的 HTTP 请求内容
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					// 如果是断开的连接错误，只记录相关错误信息，不记录完整的堆栈信息。
					global.LSys.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 返回一个友好的错误响应给客户端
					response.FailedCode(c, 10500, "服务器端异常1，请联系管理员！")
					c.Abort() // 中止请求处理
					return
				}

				// 根据是否记录堆栈信息，选择性记录日志
				if stack {
					if global.C.Logger.Format == "json" {
						global.LSys.Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
							zap.String("stack", string(debug.Stack())),
						)
					} else {
						msg := fmt.Sprintf("[Recovery from panic] %s\n%s\n", err, string(debug.Stack()))
						global.LSys.Error(msg)
					}

				} else {
					global.LSys.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				// 返回一个友好的错误响应给客户端
				response.FailedCode(c, 10500, "服务器端异常2，请联系管理员！")
				c.Abort() // 中止请求处理
			}
		}()
		c.Next() // 继续处理请求
	}
}
