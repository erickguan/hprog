package value

// move to generics soon...

import (
	"fmt"
	"strconv"
)

type VALUE_TYPE int

const (
	VT_ILLEGAL VALUE_TYPE = iota

	VT_BOOL
	VT_NIL

	VT_FLOAT
	VT_INT
	VT_COMPLEX
	VT_HEX

	VT_OBJ

	VT_NUMBER
)

var VTmap = map[VALUE_TYPE]string{
	VT_BOOL: "VT_BOOL",
	VT_NIL:  "VT_NIL",

	VT_FLOAT:   "VT_FLOAT",
	VT_INT:     "VT_INT",
	VT_COMPLEX: "VT_COMPLEX",
	VT_HEX:     "VT_HEX",
	VT_OBJ:     "VT_OBJ",

	VT_NUMBER: "VT_NUMBER",
}

type OType int

const (
	O_ILLEGAL OType = iota
	O_STRING
)

type ObjCtr struct {
	_obj  interface{}
	otype OType
	_next *ObjCtr
}

type ObjString string

type V struct {
	_bool   bool
	_int    int
	_f64    float64
	_nil    bool
	_objCtr *ObjCtr
}

// TODO: (maybe) create Value from _objCtr
// by the _objCtr Type (O_TYPE -> VT)

type Value struct {
	VT VALUE_TYPE
	_V V
}

func PrintValue(v Value) {
	vts := ""
	switch v.VT {
	case VT_INT:
		vts = strconv.Itoa(v._V._int)
	case VT_FLOAT:
		vts = strconv.FormatFloat(v._V._f64, 'E', -1, 64)
	case VT_BOOL:
		vts = strconv.FormatBool(v._V._bool)
	case VT_OBJ:
		if IsString(&v) {
			vts = *AsString(&v)
		}
	case VT_NIL:
		vts = "nil"
	}
	fmt.Printf("%s (%s)", vts, VTmap[v.VT])
}

func NewBool(value bool) Value {
	return Value{
		_V: V{_bool: value},
		VT: VT_BOOL,
	}
}

func NewInt(value int) Value {
	return Value{
		_V: V{_int: value},
		VT: VT_INT,
	}
}

func NewFloat(value float64) Value {
	return Value{
		_V: V{_f64: value},
		VT: VT_FLOAT,
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
	default:
		return Value{
			_V: V{_nil: true},
			VT: VT_NIL,
		}
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
	case VT_OBJ:
		t := ConvertToString(a) + ConvertToString(b)
		return NewString(t)
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
	case VT_BOOL:
		t := !a._V._bool
		return Value{
			_V: V{_bool: t},
			VT: VT_BOOL,
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

func Equal(a *Value, b *Value) Value {
	if a.VT != b.VT {
		return NewBool(false)
	}
	switch a.VT {
	case VT_NIL:
		return NewBool(true)
	case VT_BOOL:
		return NewBool(a._V._bool == b._V._bool)
	case VT_INT:
		return NewBool(a._V._int == b._V._int)
	case VT_FLOAT:
		return NewBool(a._V._f64 == b._V._f64)
	case VT_OBJ:
		return NewBool(ConvertToString(a) == ConvertToString(b))
	default:
		return NewBool(false)
	}
}

func Less(a *Value, b *Value) Value {
	if a.VT != b.VT {
		return NewBool(false)
	}
	switch a.VT {
	case VT_NIL:
		return NewBool(false)
	case VT_INT:
		return NewBool(a._V._int < b._V._int)
	case VT_FLOAT:
		return NewBool(a._V._f64 < b._V._f64)
	default:
		return NewBool(false)
	}
}

func Greater(a *Value, b *Value) Value {
	if a.VT != b.VT {
		return NewBool(false)
	}
	switch a.VT {
	case VT_NIL:
		return NewBool(false)
	case VT_INT:
		return NewBool(a._V._int > b._V._int)
	case VT_FLOAT:
		return NewBool(a._V._f64 > b._V._f64)
	default:
		return NewBool(false)
	}
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

func NewString(v string) Value {
	o := ObjCtr{
		_obj:  &v,
		otype: O_STRING,
	}
	return Value{
		_V: V{_objCtr: &o},
		VT: VT_OBJ,
	}
}

/*
func ObjAsValue(o *Obj) Value {
	return Value{_V: V{_obj: o}, VT: VT_OBJ}
}
*/

func ConvertToString(v *Value) string {
	return *AsString(v) //, true
}

func FreeObj(v *Value)          { v._V._objCtr = nil }
func AsString(v *Value) *string { return v._V._objCtr._obj.(*string) }
func IsString(v *Value) bool    { return ObjType(v) == O_STRING }
func ObjType(v *Value) OType    { return AsObj(v).otype }
func AsObj(v *Value) *ObjCtr    { return v._V._objCtr }
func IsObj(v *Value) bool       { return v.VT == VT_OBJ }

func IsNumberType(v VALUE_TYPE) bool             { return v == VT_FLOAT || v == VT_INT }
func IsSameType(a VALUE_TYPE, b VALUE_TYPE) bool { return a == b }
func IsBooleanType(v VALUE_TYPE) bool            { return v == VT_BOOL }
