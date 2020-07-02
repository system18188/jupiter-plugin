package bd

import (
	"github.com/system18188/jupiter-plugin/store/dbr"
)

type Builder interface {
	Build() (dbr.Builder, bool)
}

func In(cond ...Builder) dbr.Builder {
	build, _ :=  And(cond...).Build()
	return build
}

type Func func() (dbr.Builder, bool)

func (b Func) Build() (dbr.Builder, bool) {
	return b()
}

// And creates AND from a list of conditions.
func And(cond ...Builder) Builder {
	return Func(func() (b dbr.Builder, p bool) {
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
	return Func(func() (b dbr.Builder, p bool) {
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
	return Func(func() (dbr.Builder, bool) {
		return dbr.Eq(column, value), ok
	})
}

// Neq is `!=`.
// When value is nil, it will be translated to `IS NOT NULL`.
// When value is a slice, it will be translated to `NOT IN`.
// Otherwise it will be translated to `!=`.
func Neq(column string, value interface{}, ok bool) Builder {
	return Func(func() (dbr.Builder, bool) {
		return dbr.Neq(column, value), ok
	})
}

// Gt is `>`.
func Gt(column string, value interface{}, ok bool) Builder {
	return Func(func() (dbr.Builder, bool) {
		return dbr.Gt(column, value), ok
	})
}

// Gte is '>='.
func Gte(column string, value interface{}, ok bool) Builder {
	return Func(func() (dbr.Builder, bool) {
		return dbr.Gte(column, value), ok
	})
}

// Lt is '<'.
func Lt(column string, value interface{}, ok bool) Builder {
	return Func(func() (dbr.Builder, bool) {
		return dbr.Lt(column, value), ok
	})
}

// Lte is `<=`.
func Lte(column string, value interface{}, ok bool) Builder {
	return Func(func() (dbr.Builder, bool) {
		return dbr.Lte(column, value), ok
	})
}

// Like is `LIKE`, with an optional `ESCAPE` clause
func Like(column, value string, ok bool, escape ...string) Builder {
	return Func(func() (dbr.Builder, bool) {
		return dbr.Like(column, value, escape...), ok
	})
}

// NotLike is `NOT LIKE`, with an optional `ESCAPE` clause
func NotLike(column, value string, ok bool, escape ...string) Builder {
	return Func(func() (dbr.Builder, bool) {
		return dbr.Like(column, value, escape...), ok
	})
}
