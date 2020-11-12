package dbr

import (
	"github.com/system18188/jupiter-plugin/pkg/check"
	"reflect"
	"strings"
)

// UpdateStmt builds `UPDATE ...`
type UpdateStmt interface {
	Builder

	Where(query interface{}, value ...interface{}) UpdateStmt
	Set(column string, value interface{}) UpdateStmt
	SetMap(m map[string]interface{}) UpdateStmt
	ScanStruct(valStruct interface{}, column ...string) UpdateStmt
}

type updateStmt struct {
	raw

	Table     string
	Value     map[string]interface{}
	WhereCond []Builder
}

// Build builds `UPDATE ...` in dialect
func (b *updateStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	if len(b.Value) == 0 {
		return ErrColumnNotSpecified
	}

	buf.WriteString("UPDATE ")
	buf.WriteString(d.QuoteIdent(b.Table))
	buf.WriteString(" SET ")

	i := 0
	for col, v := range b.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(d.QuoteIdent(col))
		buf.WriteString(" = ")
		buf.WriteString(placeholder)

		buf.WriteValue(v)
		i++
	}

	if len(b.WhereCond) > 0 {
		buf.WriteString(" WHERE ")
		err := And(b.WhereCond...).Build(d, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update creates an UpdateStmt
func Update(table string) UpdateStmt {
	return createUpdateStmt(table)
}

func createUpdateStmt(table string) *updateStmt {
	return &updateStmt{
		Table: table,
		Value: make(map[string]interface{}),
	}
}

// UpdateBySql creates an UpdateStmt with raw query
func UpdateBySql(query string, value ...interface{}) UpdateStmt {
	return createUpdateStmtBySQL(query, value)
}

func createUpdateStmtBySQL(query string, value []interface{}) *updateStmt {
	return &updateStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
		Value: make(map[string]interface{}),
	}
}

// Where adds a where condition
func (b *updateStmt) Where(query interface{}, value ...interface{}) UpdateStmt {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Builder:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}

// Set specifies a key-value pair
func (b *updateStmt) Set(column string, value interface{}) UpdateStmt {
	b.Value[column] = value
	return b
}

// SetMap specifies a list of key-value pair
func (b *updateStmt) SetMap(m map[string]interface{}) UpdateStmt {
	for col, val := range m {
		b.Set(col, val)
	}
	return b
}

// ScanStruct 扫描struct按column绑定值 (传入值, 表列名) tag 绑定db名
func (b *updateStmt) ScanStruct(valStruct interface{}, column ...string) UpdateStmt {
	tyOf := reflect.TypeOf(valStruct)
	vaOf := reflect.ValueOf(valStruct)
	if tyOf.Kind() == reflect.Ptr {
		tyOf = tyOf.Elem()
		vaOf = vaOf.Elem()
	}
	if tyOf.Kind() != reflect.Struct {
		return b
	}
	columnLen := len(column)
	for i := 0; i < tyOf.NumField(); i++ {
		typeField := tyOf.Field(i)  // struct取行类型
		valueField := vaOf.Field(i) // struct取行值
		// 取列名
		columnName := typeField.Tag.Get("db")
		columnNameList := strings.Split(columnName, ",")
		columnName = columnNameList[0]
		// 列名为空转到下一行
		if columnName == "" || columnName == "-" || columnName == "id" || columnName == "ID" {
			continue
		}
		// 检查列是否存在
		if columnLen > 0 && !check.IsSliceContainsString(columnName, column...) {
			continue
		}
		if valueField.CanSet() {
			b.Set(columnName, valueField.Interface())
		}
	}
	return b
}
