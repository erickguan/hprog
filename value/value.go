package value

import "fmt"

/*
type Variable struct {
	vtype string
	value interface{}
}
*/

/*
type VT_CODE int

const (
	VT_ILLEGAL VT_CODE = iota

	VT_INT
	VT_FLOAT
)
*/

type VALUE_TYPE int

const (
	VT_ILLEGAL = iota

	VT_BOOL
	VT_NIL
	VT_NUMBER
)

type V struct {
	b   bool
	i   int
	f64 float64
}
type Value struct {
	VT VALUE_TYPE
	_V V
}

func PrintValue(value Value) {
	fmt.Printf("val: %g\n", value)
}

// NOTE: Don't like this at all!
func Add(a *Value, b *Value) Value {
	t := a._V.f64 + b._V.f64
	return Value{
		_V: V{f64: t},
		VT: VT_NUMBER,
	}
}

// meh, NO!
func Create(t interface{}, vt VALUE_TYPE) Value {
	switch vt {
	case VT_NUMBER:
		return Value{
			_V: V{f64: t.(float64)},
			VT: VT_NUMBER,
		}
	default:
		// should never reach here!!!!
		return Value{}
	}
}

func Sub(a *Value, b *Value) Value {

	t := a._V.f64 - b._V.f64
	return Value{
		_V: V{f64: t},
		VT: VT_NUMBER,
	}
}

func Divide(a *Value, b *Value) Value {
	t := a._V.f64 / b._V.f64
	return Value{
		_V: V{f64: t},
		VT: VT_NUMBER,
	}
}

func Multiply(a *Value, b *Value) Value {
	t := a._V.f64 * b._V.f64
	return Value{
		_V: V{f64: t},
		VT: VT_NUMBER,
	}
}

func Negate(a Value) Value {
	if a.VT != VT_NUMBER {
		// error
	}
	t := -a._V.f64
	return Value{
		_V: V{f64: t},
		VT: VT_NUMBER,
	}
}
