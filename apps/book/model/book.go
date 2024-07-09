package model

import (
	"encoding/json"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title      string          `json:"Title" binding:"required,max=20" gorm:"type:varchar(20);unique_index;not null" example:"My Book" description:"书籍名称，必须传递"`
	PageNumber int             `json:"PageNumber" binding:"required,number" gorm:"type:int;not null"`
	Desc       string          `json:"Desc" gorm:"type:TEXT"`
	Meta       json.RawMessage `json:"Meta" gorm:"type:json;serializer:json" swaggertype:"object"` // 使用 json.RawMessage 来存储未解析的 JSON 数据
}

type GetQuery struct {
	Title string `json:"Title"`
	types.Pagination
}

func init() {
	global.M = append(global.M, &Book{})
}
