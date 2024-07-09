package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanshicheng/ikube-gin-starter/global"
)

var (
	// 维护当前所有服务
	logicApps = map[string]LogicService{}
	// gin
	ginApps = map[string]GinService{}
)

func RegistryGinRouter(svc GinService) {
	// 服务实例注册到svcs map当中
	if _, ok := ginApps[svc.Name()]; ok {
		panic(fmt.Sprintf("gin service %s has registried", svc.Name()))
	}

	ginApps[svc.Name()] = svc
}

// 通过断言自动注册
func RegistryLogic(svc LogicService) {
	// 服务实例注册到svcs map当中
	if _, ok := logicApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}

	logicApps[svc.Name()] = svc
}

// LoadedGinApp 查询加载成功的服务
func LoadedGinApp() (apps []string) {
	for k := range ginApps {
		apps = append(apps, k)
	}
	return
}

// 返回一个对象, 任何类型都可以, 使用时, 由使用方进行断言
func GetLogic(name string) interface{} {

	for k, v := range logicApps {
		if k == name {
			return v
		}
	}

	return nil
}

// 用于初始化注册到 IOC容器中的所有服务
func InitImpl() {
	for _, v := range logicApps {
		v.Config()
	}
}

// 用于初始化注册到 IOC容器中的所有服务
func InitGin() *gin.Engine {
	// 初始化 logic
	for _, v := range logicApps {
		v.Config()
	}
	for _, v := range ginApps {
		v.Config()
		global.LSys.Info(fmt.Sprintf("服务注册成功: %s", v.Name()))
	}
	// 自动注册路由
	return BusinessRouter(ginApps)
}
