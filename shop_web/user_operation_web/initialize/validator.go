package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"shop_web/user_operation_web/global"
	"strings"
)

func InitTrans(locale string)(err error)  {
	if v, ok := binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterTagNameFunc(func(fld reflect.StructField) string{
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-"{
				return ""
			}
			return name
		})

		zhT := zh.New()//中文翻译器
		enT := en.New()//英文翻译器
		uni := ut.New(enT, zhT)
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok{
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		switch locale{
		case "en":
			_ = enTranslations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			_ = zhTranslations.RegisterDefaultTranslations(v, global.Trans)
		default:
			_ = enTranslations.RegisterDefaultTranslations(v, global.Trans)
		}
		return
	}
	return
}
