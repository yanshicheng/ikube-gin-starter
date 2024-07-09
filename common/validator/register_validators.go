package validator

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// ValidatorSlice 定义了 ValidatorTranslation 的切片
var ValidatorSlice []*ValidatorTranslation

// LocaleType 定义了一种新的类型，用于限定 Locale 的值
type localeType string

// 定义 LocaleType 可以接受的常量值
const (
	LocaleEN localeType = "en"
	LocaleZH localeType = "zh"
)

// ValidatorTranslation 自定义验证器和翻译的结构定义
type ValidatorTranslation struct {
	Tag            string              // 验证器标签
	ValidationFunc validator.Func      // 验证函数
	Translations   []TranslationDetail // 翻译详情列表
}

type TranslationDetail struct {
	Locale          localeType                // 语言环境，如'en'或'zh'
	TranslationMsg  string                    // 翻译文本
	TranslationFunc validator.TranslationFunc // 翻译函数
}

// RegisterValidatorsAndTranslations 注册验证器和翻译
func RegisterValidatorsAndTranslations(validators []*ValidatorTranslation, uni *ut.UniversalTranslator) error {
	transEn, _ := uni.GetTranslator("en")
	transZh, _ := uni.GetTranslator("zh")
	for _, vt := range validators {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			if err := v.RegisterValidation(vt.Tag, vt.ValidationFunc); err != nil {
				return fmt.Errorf("failed to register validation for %s: %v", vt.Tag, err)
			}
			for _, td := range vt.Translations {
				var trans ut.Translator
				// 根据 LocaleType 获取对应的 Translator
				if td.Locale == "en" {
					trans = transEn
				} else if td.Locale == "zh" {
					trans = transZh
				}
				// 默认翻译函数
				if td.TranslationFunc == nil {
					td.TranslationFunc = defaultTranslationFunc
				}
				// 注册翻译
				err := v.RegisterTranslation(vt.Tag, trans, func(ut ut.Translator) error {
					return ut.Add(vt.Tag, td.TranslationMsg, true)
				}, td.TranslationFunc)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func defaultTranslationFunc(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(fe.Tag(), fe.Field())
	return t
}

// RegistryValidator 添加验证器
func RegistryValidator(validator *ValidatorTranslation) {
	ValidatorSlice = append(ValidatorSlice, validator)
}
