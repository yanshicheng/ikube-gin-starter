package all

import (
	_ "github.com/yanshicheng/ikube-gin-starter/apps/portal/handler"
	//_ "github.com/yanshicheng/ikube-gin-starter/apps/book/handler"
	//_ "github.com/yanshicheng/ikube-gin-starter/apps/book/logic"
	//_ "github.com/yanshicheng/ikube-gin-starter/apps/book/model"
	_ "github.com/yanshicheng/ikube-gin-starter/apps/portal/logic"
	_ "github.com/yanshicheng/ikube-gin-starter/apps/portal/model"
	// 引入自定义验证器
	_ "github.com/yanshicheng/ikube-gin-starter/common/validator"
)
