package binding

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

type String struct {
	Val   string
	Valid bool
}

// Scan 扫描写入值 interface.
func (n *String) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = "", false
		return err
	}
	err = convertAssign(&n.Val, value)
	n.Valid = err == nil
	return err
}

func (n String) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (ns String) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Val)
	}
	return nullString, nil
}

func (ns *String) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return ns.Scan(s)
}

type Strings struct {
	Val   []string
	Valid bool
}

// Scan 扫描写入值 interface.
func (n *Strings) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = make([]string, 0), false
		return err
	}
	// 取类型
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 绑断是不是数组类型
	if t.Kind() != reflect.Slice {
		return errors.New("scan on Strings: no slice")
	}
	// 付值
	v := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < t.Len(); i++ {
		var d string
		err = convertAssign(&d, v.Index(i))
		if err != nil {
			continue
		}
		n.Val = append(n.Val, d)
	}
	n.Valid = err == nil
	return err
}

func (n Strings) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (ns Strings) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Val)
	}
	return nullString, nil
}

func (ns *Strings) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return ns.Scan(s)
}
