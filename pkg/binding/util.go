package binding

import (
	"database/sql/driver"
	"reflect"
	"strings"
	"time"
)

// WriteDefaultValueOnTag 按TAG写入默认值
func WriteDefaultValueOnTag(ptr interface{}, tagName string) error {
	t := reflect.TypeOf(ptr)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	values := reflect.Indirect(reflect.ValueOf(ptr))
	// 当不是 Struct 类型跳过
	if t.Kind() != reflect.Struct {
		return nil
	}

	return writeDefaultValue(t, values, tagName)
}

func writeDefaultValue(t reflect.Type, refValue reflect.Value, tagName string) error {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	values := reflect.Indirect(refValue)
	// 当不是 Struct 类型跳过
	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		defaultValue := ""
		field := t.Field(i)
		val := values.Field(i)
		// 当他是一个地址时回调
		if field.Type.Kind() == reflect.Ptr {
			if err := writeDefaultValue(field.Type, val, tagName); err != nil {
				return err
			}
			continue
		}
		// 取得默认值
		tag := field.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			defaultValue = field.Tag.Get("default")
		} else {
			tagSlice := strings.SplitN(tag, ",", 3)
			if len(tagSlice) == 2 {
				defaultValueSlice := strings.SplitN(tagSlice[1], "=", 2)
				if defaultValueSlice[0] == "default" {
					defaultValue = defaultValueSlice[1]
				}
			}
		}
		// 给数据库类型付值
		if defaultValue != "" {
			for keyType, vFn := range RegisterFormType.Get() {
				if field.Type == keyType {
					v, err := vFn([]string{defaultValue})
					if err != nil {
						return err
					}
					val.Set(v)
					continue
				}
			}
		}

		// 检查是否存在包中类型 并给数据库db类型付值
		if field.Type.Kind() == reflect.Struct {
			// 默认值为空跳过
			if defaultValue != "" {
				// 检查是否是时间类型并付值
				if _, isTime := val.Interface().(time.Time); isTime {
					if err := setTimeField(defaultValue, field, val); err != nil {
						return err
					}
					continue
				}
				continue
			}
			// 如果是其他的 Struct
			if err := writeDefaultValue(field.Type, val, tagName); err != nil {
				return err
			}
			continue
		}
		// 默认值为空跳过
		if defaultValue == "" || !val.CanSet() {
			continue
		}

		// 给其他类型付值
		if err := setWithProperType(field.Type.Kind(), defaultValue, val); err != nil {
			return err
		}

	}
	return nil
}
// DbHasValToMap db存在值就转换数据
func DbHasValToMap(ptr interface{}) map[string]interface{} {
	return HasValToMap(ptr,"db")
}

// HasValToMap 按TAG存在值就转换数据
func HasValToMap(ptr interface{}, tag string) (kv map[string]interface{}) {
	t := reflect.TypeOf(ptr)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	values := reflect.Indirect(reflect.ValueOf(ptr))
	// 当不是 Struct 类型跳过
	if t.Kind() != reflect.Struct {
		return
	}

	return hasValToMap(t, values, tag)
}

func hasValToMap(t reflect.Type, refValue reflect.Value, tagName string) (kv map[string]interface{}) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 当不是 Struct 类型跳过
	if t.Kind() != reflect.Struct {
		return nil
	}
	values := reflect.Indirect(refValue)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := values.Field(i)
		tag := field.Tag.Get(tagName)
		// 如查不存在 tag 名称检查转入下一层
		if tag == "" || tag == "-" {
			if field.Type.Kind() == reflect.Struct {
				for key, value := range hasValToMap(field.Type, val, tagName) {
					kv[key] = value
				}
			}
			continue
		}

		// 数据库类型
		if _, ok := val.Interface().(driver.Valuer); ok {
			kv[tag] = val.Interface()
			continue
		}

		if !val.CanSet() {
			continue
		}

		// 当为期它值时直接写入值
		kv[tag] = val.Interface()

	}
	return
}
