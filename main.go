package main

import (
	"github.com/yanshicheng/ikube-gin-starter/cmd"
)

// @title IkubeOps OpenApi API
// @version 0.0.1
// @description gin 脚手架
// @termsOfService  http://swagger.io/terms/
// @contact.name   官网地址
// @contact.url    http://www.ikubeops.com
// @contact.email  ikubeops@gmail.com
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host      localhost:8080
// @BasePath /api/v1
func main() {
	cmd.Execute()
}

// swag init --parseDependency --parseInternal --parseGoList=false --parseDepth=1
