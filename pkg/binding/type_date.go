package binding

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

const (
	timeFormat10 = "2006-01-02"
)

type Date struct {
	Val   time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (n Date) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val.Format(timeFormat10), nil
}

func (t Date) String() string {
	if !t.Valid {
		return ""
	}
	return t.Val.Format(timeFormat10)
}

func (t *Date) Scan(value interface{}) error {
	var err error

	if value == nil {
		t.Val, t.Valid = time.Time{}, false
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		// 如果时区不是上海 转为上海时区
		t.Val = v
		t.Valid = true
		return nil
	case []byte:
		t.Val, err = parseDateTime(string(v), time.Local)
		t.Valid = (err == nil)
		return err
	case string:
		t.Val, err = parseDateTime(v, time.Local)
		t.Valid = (err == nil)
		return err
	}

	t.Valid = false
	return errors.New("scan on type: unknown type")
}

func (n Date) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val.Format(timeFormat10))
	}
	return nullString, nil
}

func (n *Date) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}
