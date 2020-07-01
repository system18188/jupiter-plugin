package binding

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func mapUri(ptr interface{}, m map[string][]string) error {
	return mapFormByTag(ptr, m, "uri")
}

func mapForm(ptr interface{}, form map[string][]string) error {
	return mapFormByTag(ptr, form, "form")
}

func mapArray(ptr interface{}, array map[string]interface{}) error {
	return mapArrayByTag(ptr, array, "array")
}

// mapFormByTag 分配值 （付值struct，传入form表单值）
func mapFormByTag(ptr interface{}, form map[string][]string, tag string) error {
	typ := reflect.TypeOf(ptr)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	val := reflect.Indirect(reflect.ValueOf(ptr))
	return mapFormByTagStruct(typ, val, form, tag)
}

func mapFormByTagStruct(typ reflect.Type, val reflect.Value, form map[string][]string, tag string) error {
	if typ.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		// 类型为上传的文件跳过
		if typeField.Type == tFile || typeField.Type == tFiles {
			continue
		}
		// 取得绑定名称 inputFieldName
		inputFieldName := typeField.Tag.Get(tag)
		inputFieldNameList := strings.Split(inputFieldName, ",")
		inputFieldName = inputFieldNameList[0]

		structFieldKind := structField.Kind()
		// 取得默认值写入 defaultValue
		var defaultValue string
		if len(inputFieldNameList) > 1 {
			defaultList := strings.SplitN(inputFieldNameList[1], "=", 2)
			if defaultList[0] == "default" {
				defaultValue = defaultList[1]
			}
		}
		// 绑定名为空时检查是否是指针与Struct
		if structFieldKind == reflect.Ptr {
			if !structField.Elem().IsValid() {
				structField.Set(reflect.New(structField.Type().Elem()))
			}
			structField = structField.Elem()
			structFieldKind = structField.Kind()
		}
		// 没有找个绑定名称 如果是struct就转到下一个struct
		if inputFieldName == "" {
			if structFieldKind == reflect.Struct {
				// 其他 struct
				err := mapFormByTagStruct(typeField.Type, structField, form, tag)
				if err != nil {
					return err
				}
			}
			continue
		}
		// 按tag取得form表单值
		inputValue, exists := form[inputFieldName]
		if !exists {
			if defaultValue == "" {
				continue
			}
			inputValue = make([]string, 1)
			inputValue[0] = defaultValue
		}
		if inputValue[0] == "" {
			continue
		}
		// 给数据库类型付值
		for keyType, vFn := range RegisterFormType.Get() {
			if typeField.Type == keyType {
				v, err := vFn(inputValue)
				if err != nil {
					return err
				}
				structField.Set(v)
				continue
			}
		}

		// 给Slice付值
		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
					return err
				}
			}
			structField.Set(slice)
			continue
		}

		// 检查是否是时间类型并付值
		if _, isTime := structField.Interface().(time.Time); isTime {
			if err := setTimeField(inputValue[0], typeField, structField); err != nil {
				return err
			}
			continue
		}

		if structFieldKind == reflect.Struct {
			continue
		}

		// 检查值是否可以更改
		if !structField.CanSet() {
			continue
		}

		// 给其他类型付值
		if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
			return err
		}
	}
	return nil
}

func mapArrayByTag(ptr interface{}, array map[string]interface{}, tag string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind() // 取struct类型
		inputFieldName := typeField.Tag.Get(tag)
		inputFieldNameList := strings.Split(inputFieldName, ",") // tag 内容组
		inputFieldName = inputFieldNameList[0]                   // tag 绑定名称
		var defaultValue string
		// 检查是否有设置默认值
		if len(inputFieldNameList) > 1 {
			defaultList := strings.SplitN(inputFieldNameList[1], "=", 2)
			if defaultList[0] == "default" {
				defaultValue = defaultList[1]
			}
		}
		if inputFieldName == "-" {
			continue
		}
		// 如果没有取到绑定名称取 取struct 定义名
		if inputFieldName == "" {
			inputFieldName = typeField.Name
			if structFieldKind == reflect.Ptr {
				if !structField.Elem().IsValid() {
					structField.Set(reflect.New(structField.Type().Elem()))
				}
				structField = structField.Elem()
				structFieldKind = structField.Kind()
			}
			if structFieldKind == reflect.Struct {
				err := mapArray(structField.Addr().Interface(), array)
				if err != nil {
					return err
				}
				continue
			}
		}
		// 当类型一样时直接付值
		arrVal := reflect.ValueOf(array[inputFieldName])
		if structField.Type() == arrVal.Type() {
			structField.Set(arrVal)
			continue
		}

		// 当 float 转 int
		if structFieldKind == reflect.Int32 && arrVal.Kind() == reflect.Float32 {
			structField.SetInt(int64(arrVal.Float()))
			continue
		}
		if structFieldKind == reflect.Int64 && arrVal.Kind() == reflect.Float64 {
			structField.SetInt(int64(arrVal.Float()))
			continue
		}
		// 把类型值转为 string
		inputValue, ok := array[inputFieldName].(string)
		if !ok {
			if defaultValue == "" {
				continue
			}
			inputValue = defaultValue
		}
		// 写入时间
		if _, isTime := structField.Interface().(time.Time); isTime {
			if err := setTimeField(inputValue, typeField, structField); err != nil {
				return err
			}
			continue
		}
		// 写入值
		if err := setWithProperType(typeField.Type.Kind(), inputValue, structField); err != nil {
			return err
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	case reflect.Ptr:
		if !structField.Elem().IsValid() {
			structField.Set(reflect.New(structField.Type().Elem()))
		}
		structFieldElem := structField.Elem()
		return setWithProperType(structFieldElem.Kind(), val, structFieldElem)
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}
