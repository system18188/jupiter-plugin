package binding

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

type Int32 struct {
	Val   int32
	Valid bool // Valid is true if Int32 is not NULL
}

// Scan 扫描写入值 interface.
func (n *Int32) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = 0, false
		return err
	}
	err = convertAssign(&n.Val, value)
	n.Valid = err == nil
	return err
}

func (n Int32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n Int32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val)
	}
	return nullString, nil
}

func (n *Int32) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

type Int32s struct {
	Val   []int32
	Valid bool // Valid is true if Int32 is not NULL
}

// Scan 扫描写入值 interface.
func (n *Int32s) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = make([]int32, 0), false
		return err
	}
	// 取类型
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 绑断是不是数组类型
	if t.Kind() != reflect.Slice {
		return errors.New("scan on int32s: no slice")
	}
	// 付值
	v := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < t.Len(); i++ {
		var d int32
		err = convertAssign(&d, v.Index(i))
		if err != nil {
			continue
		}
		n.Val = append(n.Val, d)
	}

	n.Valid = err == nil
	return err
}

func (n Int32s) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n Int32s) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val)
	}
	return nullString, nil
}

func (n *Int32s) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}
