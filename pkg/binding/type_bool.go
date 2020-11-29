package binding

import (
	"database/sql/driver"
	"encoding/json"
)

type Bool struct {
	Val   bool
	Valid bool // Valid is true if Bool is not NULL
}

func (n *Bool) Scan(value interface{}) (err error) {
	if value == nil {
		n.Val, n.Valid = false, false
		return err
	}
	err = convertAssign(&n.Val, value)
	n.Valid = err == nil
	return err
}

func (n Bool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Val, nil
}

func (n Bool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Val)
	}
	return nullString, nil
}

func (n *Bool) UnmarshalJSON(b []byte) error {
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return n.Scan(s)
}