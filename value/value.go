package value

import (
	"fmt"
	"strconv"
)

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

	VT_FLOAT
	VT_INT
	VT_COMPLEX
	VT_HEX

	VT_NUMBER
)

type V struct {
	_bool bool
	_int  int
	_f64  float64
	_nil  bool
}
type Value struct {
	VT VALUE_TYPE
	_V V
}

func PrintValue(value Value) {
	fmt.Printf("val: %g\n", value)
}

func NewBool(value bool, vt VALUE_TYPE) Value {
	return Value{
		_V: V{_bool: value},
		VT: VT_BOOL,
	}
}

func New(rawValue string, vt VALUE_TYPE) Value {
	switch vt {
	case VT_INT:
		b, _ := strconv.Atoi(rawValue)
		// int32
		return Value{
			_V: V{_int: b},
			VT: VT_INT,
		}
	case VT_FLOAT:
		// float64
		b, _ := strconv.ParseFloat(rawValue, 64)
		return Value{
			_V: V{_f64: b},
			VT: VT_FLOAT,
		}
	//case VT_COMPLEX:
	//case VT_HEX:
	case VT_NIL:
		return Value{
			_V: V{_nil: true},
			VT: VT_NIL,
		}
	default:
		// should never reach here!!!!
		return Value{}
	}
}

func Add(a *Value, b *Value) Value {
	t := a._V._f64 + b._V._f64
	return Value{
		_V: V{_f64: t},
		VT: VT_NUMBER,
	}
}

func Sub(a *Value, b *Value) Value {
	t := a._V._f64 - b._V._f64
	return Value{
		_V: V{_f64: t},
		VT: VT_NUMBER,
	}
}

func Divide(a *Value, b *Value) Value {
	t := a._V._f64 / b._V._f64
	return Value{
		_V: V{_f64: t},
		VT: VT_NUMBER,
	}
}

func Multiply(a *Value, b *Value) Value {
	t := a._V._f64 * b._V._f64
	return Value{
		_V: V{_f64: t},
		VT: VT_NUMBER,
	}
}

func Negate(a Value) Value {
	if !IsNumberType(a.VT) {
		// error
	}
	t := -a._V._f64
	return Value{
		_V: V{_f64: t},
		VT: VT_NUMBER,
	}
}

func DetectNumberTypeByConversion(v string) VALUE_TYPE {
	if _, err := strconv.Atoi(v); err == nil {
		return VT_INT
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return VT_FLOAT
	}
	return VT_NIL
}

func IsNumberType(v VALUE_TYPE) bool  { return v == VT_FLOAT || v == VT_INT }
func IsBooleanType(v VALUE_TYPE) bool { return v == VT_BOOL }
