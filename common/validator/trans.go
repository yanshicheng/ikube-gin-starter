package validator

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/yanshicheng/ikube-gin-starter/global"
	"reflect"
	"strings"
)

// InitTrans 初始化翻译器
func InitTrans(language string) (translator ut.Translator, uni *ut.UniversalTranslator, err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		//// 为SignUpparam注册自定义校验方法
		//for _, val := range global.SOG_VALIDATORREG {
		//	v.RegisterStructValidation(val.Func, val.Type)
		//}

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni = ut.New(enT, zhT, enT)

		// locale 通常取决于 handler 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		translator, ok = uni.GetTranslator(global.C.App.Language)
		if !ok {
			return nil, nil, fmt.Errorf("uni.GetTranslator(%s) failed", global.C.App.Language)
		}

		// 注册翻译器
		switch language {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, translator)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, translator)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, translator)
		}
		return
	}
	return
}

func RemoveTopStruct(err validator.ValidationErrors) map[string]string {
	res := map[string]string{}
	//var validationErrs validator.ValidationErrors
	//ok := errors.As(err, &validationErrs)
	//if !ok {
	//	res["error"] = fmt.Sprintf("validator 错误解析失败: %v", err)
	//	return res
	//}
	fields := err.Translate(global.IkubeopsTrans)
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}

	return res
}
