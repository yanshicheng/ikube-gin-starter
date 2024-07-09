package types

import "github.com/yanshicheng/ikube-gin-starter/common/errorx"

type DataType string

const (
	DataTypeString DataType = "string"
	DataTypeJson   DataType = "json"
	DataTypeSlice  DataType = "slice"
)

// 自定义泛型
type Data[T string | map[string]string] struct {
	Code     errorx.ErrorCode `json:"Code"`
	Data     interface{}      `json:"Data"`
	Message  T                `json:"Message" swaggertype:"string"`
	DataType DataType         `json:"DataType"`
}

type QueryResponse struct {
	Page       int         `json:"Page"`
	PageNumber int         `json:"PageNumber"`
	TotalPage  int         `json:"TotalPage"`
	Total      int         `json:"Total"`
	Data       interface{} `json:"Data"`
}
