package dbr

import (
	"bytes"
	"fmt"
	"github.com/system18188/jupiter-plugin/pkg/check"
	"reflect"
	"strings"
)

// ConflictStmt is ` ON CONFLICT ...` part of InsertStmt
type ConflictStmt interface {
	Action(column string, action interface{}) ConflictStmt
}

type conflictStmt struct {
	constraint string
	actions    map[string]interface{}
}

// Action adds action for column which will do if conflict happens
func (b *conflictStmt) Action(column string, action interface{}) ConflictStmt {
	b.actions[column] = action
	return b
}

// InsertStmt builds `INSERT INTO ...`
type InsertStmt interface {
	Builder
	Columns(column ...string) InsertStmt
	Values(value ...interface{}) InsertStmt
	ScanStruct(value interface{}, column ...string) InsertStmt
	Record(structValue interface{}) InsertStmt
	OnConflictMap(constraint string, actions map[string]interface{}) InsertStmt
	OnConflict(constraint string) ConflictStmt
	Returning(column ...string) InsertStmt
}

type insertStmt struct {
	raw

	Table    string
	Column   []string
	Value    [][]interface{}
	ReturnColumn []string
	Conflict *conflictStmt
}

// Proposed is reference to proposed value in on conflict clause
func Proposed(column string) Builder {
	return BuildFunc(func(d Dialect, b Buffer) error {
		_, err := b.WriteString(d.Proposed(column))
		return err
	})
}

// Build builds `INSERT INTO ...` in dialect
func (b *insertStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	if len(b.Column) == 0 {
		return ErrColumnNotSpecified
	}

	buf.WriteString("INSERT INTO ")
	buf.WriteString(d.QuoteIdent(b.Table))

	placeholderBuf := new(bytes.Buffer)
	placeholderBuf.WriteString("(")
	buf.WriteString(" (")
	for i, col := range b.Column {
		if i > 0 {
			buf.WriteString(",")
			placeholderBuf.WriteString(",")
		}
		buf.WriteString(d.QuoteIdent(col))
		placeholderBuf.WriteString(placeholder)
	}
	buf.WriteString(") VALUES ")
	placeholderBuf.WriteString(")")
	placeholderStr := placeholderBuf.String()

	for i, tuple := range b.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(placeholderStr)

		buf.WriteValue(tuple...)
	}
	if b.Conflict != nil && len(b.Conflict.actions) > 0 {
		keyword := d.OnConflict(b.Conflict.constraint)
		if len(keyword) == 0 {
			return fmt.Errorf("Dialect %s does not support OnConflict", d)
		}
		buf.WriteString(" ")
		buf.WriteString(keyword)
		buf.WriteString(" ")
		needComma := false
		for _, column := range b.Column {
			if v, ok := b.Conflict.actions[column]; ok {
				if needComma {
					buf.WriteString(",")
				}
				buf.WriteString(d.QuoteIdent(column))
				buf.WriteString("=")
				buf.WriteString(placeholder)
				buf.WriteValue(v)
				needComma = true
			}
		}
	}

	if len(b.ReturnColumn) > 0 {
		buf.WriteString(" RETURNING ")
		for i, col := range b.ReturnColumn {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(d.QuoteIdent(col))
		}
	}

	return nil
}

// InsertInto creates an InsertStmt
func InsertInto(table string) InsertStmt {
	return createInsertStmt(table)
}

func createInsertStmt(table string) *insertStmt {
	return &insertStmt{
		Table: table,
	}
}

// InsertBySql creates an InsertStmt from raw query
func InsertBySql(query string, value ...interface{}) InsertStmt {
	return createInsertStmtBySQL(query, value)
}

func createInsertStmtBySQL(query string, value []interface{}) *insertStmt {
	return &insertStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
	}
}

// Columns adds columns
func (b *insertStmt) Columns(column ...string) InsertStmt {
	b.Column = append(b.Column, column...)
	return b
}

// Values adds a tuple for columns
func (b *insertStmt) Values(value ...interface{}) InsertStmt {
	b.Value = append(b.Value, value)
	return b
}

// StructValues  Struct里tag 的db不等于空的写入values
func (b *insertStmt) ScanStruct(value interface{}, column ...string) InsertStmt {
	vType := reflect.TypeOf(value)
	vValue := reflect.ValueOf(value)
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
		vValue = vValue.Elem()
	}
	if vType.Kind() != reflect.Struct {
		return b
	}
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	columnLen := len(column)
	for index := 0; index < vType.NumField(); index++ {
		typeField := vType.Field(index)
		valueField := vValue.Field(index)
		columnName := typeField.Tag.Get("db")
		columnSlice := strings.SplitN(columnName, ",", 2)
		columnName = columnSlice[0]
		if columnName == "" || columnName == "id" || columnName == "-" {
			continue
		}
		// 检查列是否存在
		if columnLen > 0 && !check.IsSliceContainsString(columnName, column...) {
			continue
		}
		columns = append(columns, columnName)
		values = append(values, valueField.Interface())
	}
	return b.Columns(columns...).Values(values...)
}

// Record adds a tuple for columns from a struct
func (b *insertStmt) Record(structValue interface{}) InsertStmt {
	v := reflect.Indirect(reflect.ValueOf(structValue))

	if v.Kind() == reflect.Struct {
		var value []interface{}
		m := structMap(v.Type())
		for _, key := range b.Column {
			if index, ok := m[key]; ok {
				value = append(value, v.FieldByIndex(index).Interface())
			} else {
				value = append(value, nil)
			}
		}
		b.Values(value...)
	}
	return b
}

// OnConflictMap allows to add actions for constraint violation, e.g UPSERT
func (b *insertStmt) OnConflictMap(constraint string, actions map[string]interface{}) InsertStmt {
	b.Conflict = &conflictStmt{constraint: constraint, actions: actions}
	return b
}

// OnConflict creates an empty OnConflict section fo insert statement , e.g UPSERT
func (b *insertStmt) OnConflict(constraint string) ConflictStmt {
	b.Conflict = &conflictStmt{constraint: constraint, actions: make(map[string]interface{})}
	return b.Conflict
}

// Returning specifies the returning columns for postgres.
func (b *insertStmt) Returning(column ...string) InsertStmt {
	b.ReturnColumn = column
	return b
}
