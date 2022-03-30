package value

import (
	"fmt"
	"strconv"
)

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

func NewInt(value int, vt VALUE_TYPE) Value {
	return Value{
		_V: V{_int: value},
		VT: VT_BOOL,
	}
}

func NewFloat(value int, vt VALUE_TYPE) Value {
	return Value{
		_V: V{_int: value},
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
	switch a.VT {
	case VT_FLOAT:
		t := a._V._f64 + b._V._f64
		return Value{
			_V: V{_f64: t},
			VT: VT_FLOAT,
		}
	case VT_INT:
		t := a._V._int + b._V._int
		return Value{
			_V: V{_int: t},
			VT: VT_INT,
		}
	}
	// TODO: return error!
	return Value{}
}

func Sub(a *Value, b *Value) Value {
	switch a.VT {
	case VT_FLOAT:
		t := a._V._f64 - b._V._f64
		return Value{
			_V: V{_f64: t},
			VT: VT_FLOAT,
		}
	case VT_INT:
		t := a._V._int - b._V._int
		return Value{
			_V: V{_int: t},
			VT: VT_INT,
		}
	}
	// TODO: return error!
	return Value{}
}

func Divide(a *Value, b *Value) Value {
	switch a.VT {
	case VT_FLOAT:
		t := a._V._f64 / b._V._f64
		return Value{
			_V: V{_f64: t},
			VT: VT_FLOAT,
		}
	case VT_INT:
		t := a._V._int / b._V._int
		return Value{
			_V: V{_int: t},
			VT: VT_INT,
		}
	}
	// TODO: return error!
	return Value{}
}

func Multiply(a *Value, b *Value) Value {
	switch a.VT {
	case VT_FLOAT:
		t := a._V._f64 * b._V._f64
		return Value{
			_V: V{_f64: t},
			VT: VT_FLOAT,
		}
	case VT_INT:
		t := a._V._int * b._V._int
		return Value{
			_V: V{_int: t},
			VT: VT_INT,
		}
	}
	// TODO: return error!
	return Value{}
}

func Negate(a Value) Value {
	if !IsNumberType(a.VT) {
		// error
	}
	switch a.VT {
	case VT_FLOAT:
		t := -a._V._f64
		return Value{
			_V: V{_f64: t},
			VT: VT_NUMBER,
		}
	case VT_INT:
		t := -a._V._int
		return Value{
			_V: V{_int: t},
			VT: VT_INT,
		}
	}
	// TODO: return error!
	return Value{}
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

func ConvertToExpectedType1(a Value, v VALUE_TYPE) Value {
	_a := a
	if _a.VT != v {
		switch v {
		case VT_INT:
			_a = Value{
				_V: V{_int: int(a._V._f64)},
				VT: v,
			}
		case VT_FLOAT:
			_a = Value{
				_V: V{_f64: float64(a._V._int)},
				VT: v,
			}
		}
	}
	return _a
}

func ConvertToExpectedType2(a Value, b Value, v VALUE_TYPE) (Value, Value) {
	a = ConvertToExpectedType1(a, v)
	b = ConvertToExpectedType1(b, v)
	return a, b
}

func IsNumberType(v VALUE_TYPE) bool             { return v == VT_FLOAT || v == VT_INT }
func IsSameType(a VALUE_TYPE, b VALUE_TYPE) bool { return a == b }
func IsBooleanType(v VALUE_TYPE) bool            { return v == VT_BOOL }
