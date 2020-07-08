package condition

import (
	"github.com/system18188/jupiter-plugin/store/dbr"
	"reflect"
)

var (
	Where   = reflect.TypeOf(WhereFunc(func() (dbr.Builder, bool) { return nil, false }))
	GroupBy = reflect.TypeOf(GroupFunc(func() (dbr.Builder, bool) { return nil, false }))
	Having  = reflect.TypeOf(HavingFunc(func() (dbr.Builder, bool) { return nil, false }))
	Order   = reflect.TypeOf(OrderFunc(func() (dbr.Builder, bool) { return nil, false }))
	Limit   = reflect.TypeOf(LimitFunc(func() (dbr.Builder, bool) { return nil, false }))
	Offset  = reflect.TypeOf(OffsetFunc(func() (dbr.Builder, bool) { return nil, false }))
)

func NewBuilders() *Builders {
	return &Builders{
		Br: make([]Builder, 0),
	}
}

type (
	Builder interface {
		Build() (dbr.Builder, bool)
	}
	Builders struct {
		Br []Builder
	}
	Func func() (dbr.Builder, bool)
	WhereFunc func() (dbr.Builder, bool)
	GroupFunc func() (dbr.Builder, bool)
	HavingFunc func() (dbr.Builder, bool)
	OrderFunc func() (dbr.Builder, bool)
	LimitFunc func() (dbr.Builder, bool)
	OffsetFunc func() (dbr.Builder, bool)
)

// ScanSelectBuilder 扫描Builder写入查询SQL语句
func (b *Builders) ScanSelectBuilder(stmt dbr.SelectBuilder) {
	for _, v := range b.Br {
		switch reflect.TypeOf(v) {
		case Where:
			if build, ok := v.Build(); ok {
				stmt.Where(build)
			}
		case GroupBy:
			if build, ok := v.Build(); ok {
				stmt.GroupBy(build)
			}
		case Having:
			if build, ok := v.Build(); ok {
				stmt.Having(build)
			}
		}
	}
}

func (b *Builders) Append(v Builder) {
	b.Br = append(b.Br, v)
}

func (b Func) Build() (dbr.Builder, bool) {
	return b()
}

func (b WhereFunc) Build() (dbr.Builder, bool) {
	return b()
}

func (b GroupFunc) Build() (dbr.Builder, bool) {
	return b()
}

func (b HavingFunc) Build() (dbr.Builder, bool) {
	return b()
}

func (b OrderFunc) Build() (dbr.Builder, bool) {
	return b()
}

func (b LimitFunc) Build() (dbr.Builder, bool) {
	return b()
}

func (b OffsetFunc) Build() (dbr.Builder, bool) {
	return b()
}

func In(cond ...Builder) dbr.Builder {
	build, _ := And(cond...).Build()
	return build
}

// And creates AND from a list of conditions.
func And(cond ...Builder) Builder {
	return WhereFunc(func() (b dbr.Builder, p bool) {
		builders := make([]dbr.Builder, 0)
		for _, v := range cond {
			if build, ok := v.Build(); ok {
				builders = append(builders, build)
			}
		}
		if len(builders) == 1 {
			b = builders[0]
			return b, true
		}
		if len(builders) > 1 {
			b = dbr.And(builders...)
			return b, true
		}
		return b, false
	})
}

// Or creates OR from a list of conditions.
func Or(cond ...Builder) Builder {
	return WhereFunc(func() (b dbr.Builder, p bool) {
		builders := make([]dbr.Builder, 0)
		for _, v := range cond {
			if build, ok := v.Build(); ok {
				builders = append(builders, build)
			}
		}
		if len(builders) == 1 {
			b = builders[0]
			return b, true
		}
		if len(builders) > 1 {
			b = dbr.Or(builders...)
			return b, true
		}
		return b, false
	})
}

// Eq is `=`.
// When value is nil, it will be translated to `IS NULL`.
// When value is a slice, it will be translated to `IN`.
// Otherwise it will be translated to `=`.
func Eq(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Eq(column, value), ok
	})
}

// Neq is `!=`.
// When value is nil, it will be translated to `IS NOT NULL`.
// When value is a slice, it will be translated to `NOT IN`.
// Otherwise it will be translated to `!=`.
func Neq(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Neq(column, value), ok
	})
}

// Gt is `>`.
func Gt(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Gt(column, value), ok
	})
}

// Gte is '>='.
func Gte(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Gte(column, value), ok
	})
}

// Lt is '<'.
func Lt(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Lt(column, value), ok
	})
}

// Lte is `<=`.
func Lte(column string, value interface{}, ok bool) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Lte(column, value), ok
	})
}

// Like is `LIKE`, with an optional `ESCAPE` clause
/*func Like(column, value string, ok bool, escape ...string) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Like(column, value, escape...), ok
	})
}

// NotLike is `NOT LIKE`, with an optional `ESCAPE` clause
func NotLike(column, value string, ok bool, escape ...string) Builder {
	return WhereFunc(func() (dbr.Builder, bool) {
		return dbr.Like(column, value, escape...), ok
	})
}*/
