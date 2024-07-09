package response

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/yanshicheng/ikube-gin-starter/common/errorx"
	"github.com/yanshicheng/ikube-gin-starter/common/types"
	v "github.com/yanshicheng/ikube-gin-starter/common/validator"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func SuccessMap(c *gin.Context, data interface{}) {
	rd := &types.Data[string]{
		Code:     errorx.ErrNormal,
		Data:     data,
		Message:  "",
		DataType: types.DataTypeJson,
	}
	c.JSON(http.StatusOK, rd)
	return
}

func SuccessSlice(c *gin.Context, data interface{}) {
	rd := &types.Data[string]{
		Code:     errorx.ErrNormal,
		Data:     data,
		Message:  "",
		DataType: types.DataTypeSlice,
	}
	c.JSON(http.StatusOK, rd)
	return
}

func SuccessMapCode(c *gin.Context, code errorx.ErrorCode, data interface{}) {
	rd := &types.Data[string]{
		Code:     code,
		Data:     data,
		Message:  "",
		DataType: types.DataTypeJson,
	}
	c.JSON(http.StatusOK, rd)
	return
}

func FailedMap[T string | map[string]string](c *gin.Context, msg T) {

	rd := &types.Data[T]{
		Code:     errorx.ErrGeneric,
		Data:     "",
		Message:  msg,
		DataType: types.DataTypeJson,
	}
	c.JSON(http.StatusOK, rd)
	return
}
func FailedStr[T string | map[string]string](c *gin.Context, msg T) {

	rd := &types.Data[T]{
		Code:     errorx.ErrGeneric,
		Data:     "",
		Message:  msg,
		DataType: types.DataTypeString,
	}
	c.JSON(http.StatusOK, rd)
	return
}
func FailServerErr[T string | map[string]string](c *gin.Context, msg T) {
	rd := &types.Data[T]{
		Code:     errorx.ErrServerErr,
		Data:     "",
		Message:  msg,
		DataType: types.DataTypeString,
	}
	c.JSON(http.StatusInternalServerError, rd)
	return
}

func FailedCode[T string | map[string]string](c *gin.Context, code errorx.ErrorCode, msg T) {
	var dt types.DataType
	switch reflect.TypeOf(msg).Kind() {
	case reflect.String:
		dt = types.DataTypeString
	case reflect.Map:
		dt = types.DataTypeJson
	default:
		panic("unhandled default case")
	}
	rd := &types.Data[T]{
		Code:     code,
		Data:     "",
		Message:  msg,
		DataType: dt,
	}
	c.JSON(http.StatusOK, rd)
	return
}

func FailedParam(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		FailedCode[map[string]string](c, errorx.ErrParamParse, v.RemoveTopStruct(validationErrs))
		return
	} else {
		if err == io.EOF {
			FailedMap(c, "请求体为空或格式不正确")
			return
		} else {
			FailedMap(c, err.Error())
		}

	}
}
