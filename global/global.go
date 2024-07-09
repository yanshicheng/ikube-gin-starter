package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/yanshicheng/ikube-gin-starter/pkg/mysql"
	"github.com/yanshicheng/ikube-gin-starter/pkg/redis"
	"github.com/yanshicheng/ikube-gin-starter/pkg/types"
	"go.uber.org/zap"
)

var (
	IkubeopsTrans ut.Translator
	C             *types.Config = types.NewDefaultConfig()
	L             *zap.Logger
	LSys          *zap.Logger
	DB            *mysql.IkubeGorm
	RDB           *redis.IkubeRedis
	M             []interface{}
)
