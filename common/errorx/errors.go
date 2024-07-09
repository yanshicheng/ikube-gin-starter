package errorx

// ErrorCode 定义不同的错误类型
type ErrorCode int

const (
	ErrNormal  ErrorCode = 0
	ErrGeneric ErrorCode = 99999 // 正常
	// token 相关
	ErrTokenMissing     ErrorCode = 11000 // 缺少 Token
	ErrTokenExpired     ErrorCode = 11001 // Token 失效
	ErrTokenInvalid     ErrorCode = 11002 // Token 解析失败
	ErrTokenRefresh     ErrorCode = 11003 // Token 刷新失败
	ErrTokenBlacklisted ErrorCode = 11004 // Token 被列入黑名单
	ErrLoginExpired     ErrorCode = 10110 // 登录过期
	ErrLoginInvalid     ErrorCode = 10111 // 登录信息无效
	// 权限相关

	ErrPermissionDenied ErrorCode = 10130 // 权限不足
	ErrRoleNotFound     ErrorCode = 10131 // 角色未找到

	// 数据库相关
	ErrDatabase     ErrorCode = 18000 // 数据库相关错误
	ErrDataConflict ErrorCode = 18004 // 数据冲突

	// 业务相关
	ErrBusinessLogic ErrorCode = 10101 // 业务逻辑错误
	ErrParamParse    ErrorCode = 10102 // 参数解析失败
	ErrDataNotFound  ErrorCode = 10103 // 数据未找到
	ErrDataCreation  ErrorCode = 10104 // 数据创建失败
	ErrDataDeletion  ErrorCode = 10105 // 数据删除失败
	ErrToOperation   ErrorCode = 10106 // "to" 操作相关错误

	// 服务器相关错误
	ErrServerErr ErrorCode = 10500
)
