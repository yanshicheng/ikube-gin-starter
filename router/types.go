package router

import (
	"github.com/gin-gonic/gin"
)

type GinService interface {
	AuthRegistry(r gin.IRouter)
	PublicRegistry(r gin.IRouter)
	Config()
	Name() string
}

type LogicService interface {
	Config()
	Name() string
}
