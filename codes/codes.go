package codes

type INSTRUC int

const (
	INSTRUC_ILLEGAL INSTRUC = iota

	INSTRUC_ADDITION
	INSTRUC_SUBSTRACT
	INSTRUC_MULTIPLY
	INSTRUC_DIVIDE

	INSTRUC_CONSTANT
	INSTRUC_NEGATE
	INSTRUC_NOT

	INSTRUC_FALSE
	INSTRUC_TRUE

	INSTRUC_EQUAL
	INSTRUC_GREATER
	INSTRUC_LESS

	INSTRUC_DECL_GLOBAL
	INSTRUC_SET_DECL_GLOBAL
	INSTRUC_GET_DECL_GLOBAL

	INSTRUC_DECL_LOCAL
	INSTRUC_SET_DECL_LOCAL
	INSTRUC_GET_DECL_LOCAL

	INSTRUC_NIL
	INSTRUC_POP

	INSTRUC_PRINT
	INSTRUC_RETURN
	INSTRUC_ERR
)
