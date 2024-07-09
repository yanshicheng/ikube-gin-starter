package model

import "github.com/yanshicheng/ikube-gin-starter/global"

func Register(model ...interface{}) {
	global.M = append(global.M, model...)
}
