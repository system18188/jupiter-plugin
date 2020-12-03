// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"sync"
)

type defaultValidator struct {
	once     sync.Once
	Trans    ut.Translator
	uni      *ut.UniversalTranslator
	validate *validator.Validate
}

var _ StructValidator = &defaultValidator{}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		en := en.New()
		v.uni = ut.New(en, en)
		v.Trans, _ = v.uni.GetTranslator("en")
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		en_translations.RegisterDefaultTranslations(v.validate, v.Trans)
		// 注册string数组类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Strings); ok {
				return val.Val
			}
			return nil
		}, Strings{})
		// 注册int32数组类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Int32s); ok {
				return val.Val
			}
			return nil
		}, Int32s{})
		// 注册int64数组类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Int64s); ok {
				return val.Val
			}
			return nil
		}, Int64s{})
		// 注册float32数组类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Float32s); ok {
				return val.Val
			}
			return nil
		}, Float32s{})
		// 注册float64数组类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Float64s); ok {
				return val.Val
			}
			return nil
		}, Float64s{})
		// 注册stirng类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(String); ok {
				return val.Val
			}
			return nil
		}, String{})
		// 注册int32类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Int32); ok {
				return val.Val
			}
			return nil
		}, Int32{})
		// 注册int64类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Int64); ok {
				return val.Val
			}
			return nil
		}, Int64{})
		// 注册flat32类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Float32); ok {
				return val.Val
			}
			return nil
		}, Float32{})
		// 注册float64类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Float64); ok {
				return val.Val
			}
			return nil
		}, Float64{})
		// 注册时间类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Time); ok {
				return val.String()
			}
			return nil
		}, Time{})
		// 注册date类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Date); ok {
				return val.String()
			}
			return nil
		}, Date{})
		// 注册bool类型
		v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if val, ok := field.Interface().(Bool); ok {
				return val.Val
			}
			return nil
		}, Bool{})
		// 添加验证HTML
		v.validate.RegisterValidation("nohtml", IsNoHTML)
		v.validate.RegisterTranslation("nohtml", v.Trans, func(ut ut.Translator) error {
			return ut.Add("nohtml", "{0} You cannot include HTML tags!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("nohtml", fe.Field())
			return t
		})

		// 添加验证数据库字段名格式
		v.validate.RegisterValidation("column", IsColumn)

		// 添加验证日期 0000-00-00
		v.validate.RegisterValidation("isDate", func(fl validator.FieldLevel) bool {
			return dateRegex.MatchString(fl.Field().String())
		})
		v.validate.RegisterTranslation("isDate", v.Trans, func(ut ut.Translator) error {
			return ut.Add("isDate", "{0} Not a date format (2000-01-01)!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("isDate", fe.Field())
			return t
		})

		// 验证日期时间格式： 0000-00-00 00:00:00
		v.validate.RegisterValidation("isTime", func(fl validator.FieldLevel) bool {
			return timeRegex.MatchString(fl.Field().String())
		})
		v.validate.RegisterTranslation("isTime", v.Trans, func(ut ut.Translator) error {
			return ut.Add("isTime", "{0} Not a date format (2006-01-02 03:04:05)!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("isTime", fe.Field())
			return t
		})
	})
}


func (v *defaultValidator) AddValidation(tag string, fn func(fl validator.FieldLevel) bool) error {
	v.lazyinit()
	return v.validate.RegisterValidation(tag, fn)
}

func (v *defaultValidator) AddTranslation(tag string, addfn func(ut ut.Translator) error, translatorFn func(ut ut.Translator, fe validator.FieldError) string) error {
	v.lazyinit()
	return v.validate.RegisterTranslation(tag, v.Trans, addfn, translatorFn)
}

// AddValidateType implements validator.CustomTypeFunc
func (v *defaultValidator) AddValidateType(CustomTypeFunc func(field reflect.Value) interface{}, types ...interface{}) {
	v.lazyinit()
	v.validate.RegisterCustomTypeFunc(CustomTypeFunc, types...)
}

