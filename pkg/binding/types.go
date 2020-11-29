package binding


import (
	"reflect"
	"sync"
	"time"
)

var (
	nullString   = []byte("null")
	ScannerType  = reflect.TypeOf((*Scanner)(nil)).Elem()
	StringType   = reflect.TypeOf(String{})
	Int32Type    = reflect.TypeOf(Int32{})
	IntType      = reflect.TypeOf(Int{})
	Int64Type    = reflect.TypeOf(Int64{})
	Float32Type  = reflect.TypeOf(Float32{})
	Float64Type  = reflect.TypeOf(Float64{})
	StringsType  = reflect.TypeOf(Strings{})
	Int32sType   = reflect.TypeOf(Int32s{})
	IntsType     = reflect.TypeOf(Ints{})
	Int64sType   = reflect.TypeOf(Int64s{})
	Float32sType = reflect.TypeOf(Float32s{})
	Float64sType = reflect.TypeOf(Float64s{})
	BoolType     = reflect.TypeOf(Bool{})
	TimeType     = reflect.TypeOf(Time{})
	DateType     = reflect.TypeOf(Date{})
	sysTimeType  = reflect.TypeOf(time.Time{})
	RegisterFormType = NewBindType()
)

// Scanner is an interface used by Scan.
type Scanner interface {
	Scan(value interface{}) (err error)
}

type bindType struct {
	m    map[reflect.Type]FormBindType
	lock *sync.RWMutex
}

// RegisterFormBindType
// 表单绑定类型注册(表单值) 付值,错误
type FormBindType func(inputValue []string) (reflect.Value, error)

func NewBindType() *bindType {
	return &bindType{
		m:    make(map[reflect.Type]FormBindType, 0),
		lock: new(sync.RWMutex),
	}
}

// Bind 绑定类型
func (p *bindType) Bind(t reflect.Type, fn FormBindType) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.m[t] = fn
}

// 取得数据
func (p *bindType) Get() map[reflect.Type]FormBindType {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.m
}

func (p *bindType) Exists(p2 reflect.Type) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	for key := range p.m {
		if key == p2 {
			return true
		}
	}
	return false
}

func setTypesField(value reflect.Value, inputValue []string) error {
	// 给数据库类型付值
	for keyType, vFn := range RegisterFormType.Get() {
		if value.Type() == keyType {
			v, err := vFn(inputValue)
			if err != nil {
				return err
			}
			value.Set(v)
			continue
		}
	}
	return nil
}

func init() {
	RegisterFormType.Bind(BoolType, func(inputValue []string) (reflect.Value, error) {
		v := Bool{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Bool{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(TimeType, func(inputValue []string) (reflect.Value, error) {
		v := Time{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Time{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(DateType, func(inputValue []string) (reflect.Value, error) {
		v := Date{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Date{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(StringType, func(inputValue []string) (reflect.Value, error) {
		v := String{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(String{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(IntType, func(inputValue []string) (reflect.Value, error) {
		v := Int{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Int{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Int32Type, func(inputValue []string) (reflect.Value, error) {
		v := Int32{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Int32{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Int64Type, func(inputValue []string) (reflect.Value, error) {
		v := Int64{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Int64{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Float32Type, func(inputValue []string) (reflect.Value, error) {
		v := Float32{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Float32{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Float64Type, func(inputValue []string) (reflect.Value, error) {
		v := Float64{}
		err := v.Scan(inputValue[0])
		if err != nil {
			return reflect.ValueOf(Float64{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(StringsType, func(inputValue []string) (reflect.Value, error) {
		v := Strings{}
		err := v.Scan(inputValue)
		if err != nil {
			return reflect.ValueOf(Strings{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(IntsType, func(inputValue []string) (reflect.Value, error) {
		v := Ints{}
		err := v.Scan(inputValue)
		if err != nil {
			return reflect.ValueOf(Ints{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Int32sType, func(inputValue []string) (reflect.Value, error) {
		v := Int32s{}
		err := v.Scan(inputValue)
		if err != nil {
			return reflect.ValueOf(Int32s{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Int64sType, func(inputValue []string) (reflect.Value, error) {
		v := Int64s{}
		err := v.Scan(inputValue)
		if err != nil {
			return reflect.ValueOf(Int64s{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Float32sType, func(inputValue []string) (reflect.Value, error) {
		v := Float32s{}
		err := v.Scan(inputValue)
		if err != nil {
			return reflect.ValueOf(Float32s{}), err
		}
		return reflect.ValueOf(v), err
	})

	RegisterFormType.Bind(Float64sType, func(inputValue []string) (reflect.Value, error) {
		v := Float64s{}
		err := v.Scan(inputValue)
		if err != nil {
			return reflect.ValueOf(Float64s{}), err
		}
		return reflect.ValueOf(v), err
	})
}
