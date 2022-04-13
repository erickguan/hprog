package token

import "fmt"

type TokenType int

const EoF = -1
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
	SINGLE_QUOTE
	NEW_LINE

	COLON

	EQUAL

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
	// FLOAT
	// INT
	// COMPLEX

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
	PRINT
	RETURN
	VAR
	WHILE

	ARGS

	AND
	OR

	CLASS

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
	"=":  EQUAL,
	"==": EQUAL_EQUAL,

	// Keywords
	"_":      PLACEHOLDER,
	"if":     IF,
	"else":   ELSE,
	"define": DEFINE,
	"decl":   DECLARE,

	"for":   FOR,
	"while": WHILE,

	"args": ARGS,

	"and": AND,
	"or":  OR,

	"False": BOOL_FALSE,
	"True":  BOOL_TRUE,

	"#":      COMMENT,
	"fn":     FUNCTION,
	"print":  PRINT,
	"return": RETURN,
	"var":    VAR,
	"nil":    NIL,

	// for dbg
	"<STRING>":     STRING,
	"<IDENTIFIER>": IDENTIFIER,
	"<NUMBER>":     NUMBER,
	// for debugging
	"ERROR": ERR,
	"\\0":   EOF,
	"\\n":   NEW_LINE,
}

var ReversedTokenMap = reverseMap(TokenMap)

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

func Print(token *Token) {
	tokenTypeReadable, _ := ReversedTokenMap[token.Type]
	printFormat := "type: %s, position: %d, line:%d, value: %s\n"

	fmt.Printf(printFormat, tokenTypeReadable, token.Position, token.Line, token.Value)
}
