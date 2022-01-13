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

type Value float64

func PrintValue(value Value) {
	fmt.Printf("val: %g\n", value)
}
