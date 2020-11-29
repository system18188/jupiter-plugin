// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

type defaultValidator struct {
	once     sync.Once
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
		v.validate = validator.New()
		v.validate.SetTagName("binding")
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
	})
}
