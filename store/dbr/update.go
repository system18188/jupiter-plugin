package dbr

import (
	"context"
	"database/sql"
	"github.com/system18188/jupiter-plugin/pkg/check"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// UpdateStmt builds `UPDATE ...`.
type UpdateStmt struct {
	runner
	EventReceiver
	Dialect

	raw

	Table        string
	Value        map[string]interface{}
	WhereCond    []Builder
	ReturnColumn []string
	LimitCount   int64
	comments     Comments
}

type UpdateBuilder = UpdateStmt

func (b *UpdateStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	if len(b.Value) == 0 {
		return ErrColumnNotSpecified
	}

	err := b.comments.Build(d, buf)
	if err != nil {
		return err
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

	if len(b.ReturnColumn) > 0 {
		buf.WriteString(" RETURNING ")
		for i, col := range b.ReturnColumn {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(d.QuoteIdent(col))
		}
	}

	if b.LimitCount >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.FormatInt(b.LimitCount, 10))
	}

	return nil
}

// Update creates an UpdateStmt.
func Update(table string) *UpdateStmt {
	return &UpdateStmt{
		Table:      table,
		Value:      make(map[string]interface{}),
		LimitCount: -1,
	}
}

// Update creates an UpdateStmt.
func (sess *Session) Update(table string) *UpdateStmt {
	b := Update(table)
	b.runner = sess
	b.EventReceiver = sess.EventReceiver
	b.Dialect = sess.Dialect
	return b
}

// Update creates an UpdateStmt.
func (tx *Tx) Update(table string) *UpdateStmt {
	b := Update(table)
	b.runner = tx
	b.EventReceiver = tx.EventReceiver
	b.Dialect = tx.Dialect
	return b
}

// UpdateBySql creates an UpdateStmt with raw query.
func UpdateBySql(query string, value ...interface{}) *UpdateStmt {
	return &UpdateStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
		Value:      make(map[string]interface{}),
		LimitCount: -1,
	}
}

// UpdateBySql creates an UpdateStmt with raw query.
func (sess *Session) UpdateBySql(query string, value ...interface{}) *UpdateStmt {
	b := UpdateBySql(query, value...)
	b.runner = sess
	b.EventReceiver = sess.EventReceiver
	b.Dialect = sess.Dialect
	return b
}

// UpdateBySql creates an UpdateStmt with raw query.
func (tx *Tx) UpdateBySql(query string, value ...interface{}) *UpdateStmt {
	b := UpdateBySql(query, value...)
	b.runner = tx
	b.EventReceiver = tx.EventReceiver
	b.Dialect = tx.Dialect
	return b
}

// Where adds a where condition.
// query can be Builder or string. value is used only if query type is string.
func (b *UpdateStmt) Where(query interface{}, value ...interface{}) *UpdateStmt {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Builder:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}

// Returning specifies the returning columns for postgres.
func (b *UpdateStmt) Returning(column ...string) *UpdateStmt {
	b.ReturnColumn = column
	return b
}

// Set updates column with value.
func (b *UpdateStmt) Set(column string, value interface{}) *UpdateStmt {
	b.Value[column] = value
	return b
}

// SetMap specifies a map of (column, value) to update in bulk.
func (b *UpdateStmt) SetMap(m map[string]interface{}) *UpdateStmt {
	for col, val := range m {
		b.Set(col, val)
	}
	return b
}

// 扫描struct按column绑定值 (传入值, 表列名) tag 绑定db名
func (b *UpdateStmt) Scan(formStruct interface{}, columns ...string) *UpdateStmt {
	tyof := reflect.TypeOf(formStruct)
	vaof := reflect.ValueOf(formStruct)
	if tyof.Kind() == reflect.Ptr {
		tyof = tyof.Elem()
		vaof = vaof.Elem()
	}
	if tyof.Kind() != reflect.Struct {
		return b
	}
	columnlen := len(columns)
	for i := 0; i < tyof.NumField(); i++ {
		typeField := tyof.Field(i)  // struct取行类型
		valueField := vaof.Field(i) // struct取行值
		// 取列名
		columnName := typeField.Tag.Get("db")
		formName := typeField.Tag.Get("form")
		columnNameList := strings.Split(columnName, ",")
		formNameList := strings.Split(formName, ",")
		columnName = columnNameList[0]
		formName = formNameList[0]
		// 列名为空转到下一行
		if columnName == "" || formName == "" || columnName == "-" || formName == "-" {
			continue
		}
		// 检查类名是否存在
		if columnlen > 0 && !check.IsSliceContainsString(formName, columns...) {
			continue
		}
		switch typeField.Type.Kind() {
		case reflect.String:
			b.Set(columnName, valueField.String())
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			b.Set(columnName, valueField.Int())
		case reflect.Bool:
			b.Set(columnName, valueField.Bool())
		case reflect.Float32, reflect.Float64:
			b.Set(columnName, valueField.Float())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b.Set(columnName, valueField.Uint())
		default:
			if now, isTime := valueField.Interface().(time.Time); isTime {
				b.Set(columnName, now)
				continue
			}
		}
	}
	return b
}

// IncrBy increases column by value
func (b *UpdateStmt) IncrBy(column string, value interface{}) *UpdateStmt {
	b.Value[column] = Expr("? + ?", I(column), value)
	return b
}

// DecrBy decreases column by value
func (b *UpdateStmt) DecrBy(column string, value interface{}) *UpdateStmt {
	b.Value[column] = Expr("? - ?", I(column), value)
	return b
}

func (b *UpdateStmt) Limit(n uint64) *UpdateStmt {
	b.LimitCount = int64(n)
	return b
}

func (b *UpdateStmt) Comment(comment string) *UpdateStmt {
	b.comments = b.comments.Append(comment)
	return b
}

func (b *UpdateStmt) Exec() (sql.Result, error) {
	return b.ExecContext(context.Background())
}

func (b *UpdateStmt) ExecContext(ctx context.Context) (sql.Result, error) {
	return exec(ctx, b.runner, b.EventReceiver, b, b.Dialect)
}

func (b *UpdateStmt) LoadContext(ctx context.Context, value interface{}) error {
	_, err := query(ctx, b.runner, b.EventReceiver, b, b.Dialect, value)
	return err
}

func (b *UpdateStmt) Load(value interface{}) error {
	return b.LoadContext(context.Background(), value)
}
