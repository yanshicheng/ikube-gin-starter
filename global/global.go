package global

import (
	"github.com/yanshicheng/ikube-gin-starter/pkg/mysql"
	"github.com/yanshicheng/ikube-gin-starter/pkg/types"
	"go.uber.org/zap"
)

var (
	C    *types.Config = types.NewDefaultConfig()
	L    *zap.Logger
	LSys *zap.Logger
	DB   *mysql.IkubeGorm
	M    []interface{}
)
