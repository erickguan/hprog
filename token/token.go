package token

import "fmt"

type TokenType int

var EoF = rune(0)

const (
	ILLEGAL TokenType = iota

	// single char tokens
	OP
	CP
	LB
	RB
	PLUS
	SLASH
	STAR
	COMMA
	DOT
	MINUS
	SEMICOLON
	QUOTE

	COLON

	ASSIGN

	GREATER
	GREATER_EQUAL
	EXCL
	EXCL_EQUAL
	LESS
	LESS_EQUAL
	EQUAL_EQUAL

	// Literals
	IDENTIFIER
	NUMBER
	NIL
	STRING

	// Keywords
	IF
	FOR
	ELSE
	NOT
	PLACEHOLDER
	DEFINE
	DECLARE
	FUNCTION

	ARGS

	AND
	OR

	BOOL_FALSE
	BOOL_TRUE

	// types

	COMMENT
	COMMENT_MULTILINE
	ERR
	EOF
	EOP // end of operation
)

var TokenMap = map[string]TokenType{
	// single
	// char
	// tokens
	"(":  OP,
	")":  CP,
	"{":  LB,
	"}":  RB,
	"+":  PLUS,
	"/":  SLASH,
	"*":  STAR,
	",":  COMMA,
	".":  DOT,
	"-":  MINUS,
	";":  SEMICOLON,
	":":  COLON,
	"\"": QUOTE,

	">":  GREATER,
	">=": GREATER_EQUAL,
	"!":  EXCL,
	"!=": EXCL_EQUAL,
	"<":  LESS,
	"<=": LESS_EQUAL,
	"=":  ASSIGN,
	"==": EQUAL_EQUAL,

	// Keywords
	"_":      PLACEHOLDER,
	"if":     IF,
	"else":   ELSE,
	"define": DEFINE,
	"decl":   DECLARE,

	"for": FOR,

	"args": ARGS,

	"and": AND,
	"or":  OR,

	"false": BOOL_FALSE,
	"true":  BOOL_TRUE,

	"#":    COMMENT,
	"func": FUNCTION,

	// reserver
	// for
	// dbg
	"<STRING>":     STRING,
	"<IDENTIFIER>": IDENTIFIER,
	"<NUMBER>":     NUMBER,
	"<NIL>":        NIL,
	// for debugging
	"ERROR": ERR,
	"\\0":   EOF,
}

var ReverseKeys = reverseMap(TokenMap)

func reverseMap(m map[string]TokenType) map[TokenType]string {
	n := make(map[TokenType]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

type Token struct {
	Type     TokenType
	Position int
	Line     int
	Value    string
}

func Print(token Token) {
	tokenTypeReadable, _ := ReverseKeys[token.Type]
	printFormat := "type: %s, position: %d, line:%d, value: %s\n"

	fmt.Printf(printFormat, tokenTypeReadable, token.Position, token.Line, token.Value)
}
