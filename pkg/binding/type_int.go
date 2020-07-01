package binding

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

type Int struct {
	Val   int
	Valid bool // Valid is true if Int is not NULL
}

// Scan 扫描写入值 interface.
func (n *Int) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = 0, false
		return err
	}
	err = convertAssign(&n.Val, value)
	n.Valid = err == nil
	return err
}

func (n Int) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n Int) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val)
	}
	return nullString, nil
}

func (n *Int) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}

type Ints struct {
	Val   []int
	Valid bool // Valid is true if Int is not NULL
}

// Scan 扫描写入值 interface.
func (n *Ints) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = make([]int, 0), false
		return err
	}
	// 取类型
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 绑断是不是数组类型
	if t.Kind() != reflect.Slice {
		return errors.New("scan on ints: no slice")
	}
	// 付值
	v := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < t.Len(); i++ {
		var d int
		err = convertAssign(&d, v.Index(i))
		if err != nil {
			continue
		}
		n.Val = append(n.Val, d)
	}

	n.Valid = err == nil
	return err
}

func (n Ints) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n Ints) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val)
	}
	return nullString, nil
}

func (n *Ints) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}
