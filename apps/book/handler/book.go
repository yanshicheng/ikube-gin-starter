package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-starter/apps/book"
	"github.com/yanshicheng/ikube-gin-starter/apps/book/logic"
	"github.com/yanshicheng/ikube-gin-starter/apps/book/model"
	"github.com/yanshicheng/ikube-gin-starter/common/response"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"github.com/yanshicheng/ikube-gin-starter/router"
	"go.uber.org/zap"
)

var handler = &BookHandler{}

type BookHandler struct {
	l   *zap.Logger
	svc *logic.BookService
}

func (h *BookHandler) Name() string {
	return book.AppBook
}

// Config 配置函数，在这里注入依赖，并且初始化实例，供其他函数使用。
func (h *BookHandler) Config() {
	h.l = global.L.Named(book.AppBook).Named("handler")
	h.svc = router.GetLogic(book.AppBook).(*logic.BookService)
}

// PublicRegistry 注册公开接口
func (h *BookHandler) PublicRegistry(r gin.IRouter) {

}

// AuthRegistry 注册认证接口
func (h *BookHandler) AuthRegistry(r gin.IRouter) {
	// 分组路由
	group := r.Group("v1/book-shelf")
	{
		// group.GET("/list", h.List)
		group.POST("/book", h.create)
		group.GET("/book", h.list)
		group.GET("/book/:id", h.get)
		group.DELETE("/book/:id", h.delete)
		group.PUT("/book/:id", h.put)
	}
}

// GetPostListHandler2 书籍管理接口
// @Summary 创建书籍接口
// @Description 创建书籍接口
// @Tags 书籍管理
// @Accept application/json
// @Produce application/json
// @Param req body  model.Book true "书籍名称" example "My Book"
// @Example request:
//
//	{
//	  "Title": "My Book",
//	  "PageNumber": 300,
//	  "Desc": "This is a great book",
//	  "Meta": {
//	    "details": "additional details here"
//	  }
//	}
//
// @success 200 {object} types.Data{data=model.Book} "desc"
// @Router /book-shelf/book [post]
func (h *BookHandler) create(c *gin.Context) {
	// @Security ApiKeyAuth

	var b model.Book
	// 绑定参数
	if err := c.ShouldBindJSON(&b); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	h.l.Debug(fmt.Sprintf("绑定数据: %+v", b))
	// 调用业务逻辑
	if err := h.svc.Create(&b); err != nil {
		response.FailedStr(c, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("创建成功: %+v", b))
		response.SuccessMap(c, b)
	}

}

func (h *BookHandler) list(c *gin.Context) {
	var query model.GetQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	h.l.Debug(fmt.Sprintf("query: %+v", query))
	if s, err := h.svc.List(&query); err != nil {
		response.FailedStr(c, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("查询成功: %+v", s))
		response.SuccessSlice(c, s)
	}
}

// @Summary 获取用户信息
// @Description 获取当前用户的详细信息
// @Accept  json
// @Produce  json
// @Security types.SearchId
// @Success 200 {object} model.Book
// @Router /api/user/info [get]
func (h *BookHandler) delete(c *gin.Context) {
	var id types.SearchId
	if err := c.ShouldBindUri(&id); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(c, err)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.FailedStr(c, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("删除成功: %+v", id))
		response.SuccessMap(c, "")
	}

}

func (h *BookHandler) get(ctx *gin.Context) {
	var id types.SearchId
	if err := ctx.ShouldBindUri(&id); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(ctx, err)
		return
	}
	if s, err := h.svc.Get(id); err != nil {
		response.FailedStr(ctx, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("查询成功: %+v", s))
		response.SuccessMap(ctx, s)
	}
}

func (h *BookHandler) put(ctx *gin.Context) {
	var id types.SearchId
	if err := ctx.ShouldBindUri(&id); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(ctx, err)
		return
	}
	var b model.Book
	if err := ctx.ShouldBindJSON(&b); err != nil {
		h.l.Error(fmt.Sprintf("数据绑定失败: %s", err))
		response.FailedParam(ctx, err)
		return
	}
	if err := h.svc.Update(id, &b); err != nil {
		response.FailedStr(ctx, err.Error())
	} else {
		h.l.Debug(fmt.Sprintf("更新成功: %+v", b))
		response.SuccessMap(ctx, b)
	}
}

func init() {
	router.RegistryGinRouter(handler)
}
