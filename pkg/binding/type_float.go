package binding

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

type Float32 struct {
	Val   float32
	Valid bool
}

func (n *Float32) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = 0, false
		return err
	}
	err = convertAssign(&n.Val, value)
	n.Valid = err == nil
	return err
}

func (n Float32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n Float32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val)
	}
	return nullString, nil
}

func (n *Float32) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

type Float32s struct {
	Val   []float32
	Valid bool
}

func (n *Float32s) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = make([]float32, 0), false
		return err
	}
	// 取类型
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 绑断是不是数组类型
	if t.Kind() != reflect.Slice {
		return errors.New("scan on Float32s: no slice")
	}
	// 付值
	v := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < t.Len(); i++ {
		var d float32
		err = convertAssign(&d, v.Index(i))
		if err != nil {
			continue
		}
		n.Val = append(n.Val, d)
	}
	n.Valid = err == nil
	return err
}

func (n Float32s) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n Float32s) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val)
	}
	return nullString, nil
}

func (n *Float32s) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}
