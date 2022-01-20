package lexer

import (
	"reflect"
	"testing"

	"github.com/badc0re/hprog/token"
)

func TestParsingNumber(t *testing.T) {
	var caseMap = map[string]token.TokenType{
		".11":   token.ERR,
		"11.":   token.ERR,
		"11a":   token.ERR,
		"11.a0": token.ERR,
		"11a0":  token.ERR,
		"-11":   token.NUMBER,
		"a11":   token.IDENTIFIER,
		"1.0":   token.NUMBER,
	}
	evalCase(t, caseMap)
}

func TestParsingString(t *testing.T) {
	/*
		var caseMap = map[string]token.TokenType{
			"\"dame\"": token.STRING,
			"\"dame":   token.ERR,
			"'dame":    token.ERR,
		}
		evalCase(t, caseMap)
	*/
}

func TestParsingExpression(t *testing.T) {
	caseMap := map[string][]token.TokenType{
		"1 + 2":                          []token.TokenType{token.NUMBER, token.PLUS, token.NUMBER},
		"1.2 + 3":                        []token.TokenType{token.NUMBER, token.PLUS, token.NUMBER},
		"((1 + 2) - 3)":                  []token.TokenType{token.OP, token.OP, token.NUMBER, token.PLUS, token.NUMBER, token.CP, token.MINUS, token.NUMBER, token.CP},
		"a = 4":                          []token.TokenType{token.IDENTIFIER, token.EQUAL, token.NUMBER},
		"a = b + c":                      []token.TokenType{token.IDENTIFIER, token.EQUAL, token.IDENTIFIER, token.PLUS, token.IDENTIFIER},
		"decl a = 10":                    []token.TokenType{token.DECLARE, token.IDENTIFIER, token.EQUAL, token.NUMBER},
		"(a == 10)":                      []token.TokenType{token.OP, token.IDENTIFIER, token.EQUAL_EQUAL, token.NUMBER, token.CP},
		"(a >= 10)":                      []token.TokenType{token.OP, token.IDENTIFIER, token.GREATER_EQUAL, token.NUMBER, token.CP},
		"(a <= 10)":                      []token.TokenType{token.OP, token.IDENTIFIER, token.LESS_EQUAL, token.NUMBER, token.CP},
		"if":                             []token.TokenType{token.IF},
		"(false == true)":                []token.TokenType{token.OP, token.BOOL_FALSE, token.EQUAL_EQUAL, token.BOOL_TRUE, token.CP},
		"decl b = 10; # (if equal true)": []token.TokenType{token.DECLARE, token.IDENTIFIER, token.EQUAL, token.NUMBER, token.SEMICOLON, token.COMMENT},
	}
	evalExpr(t, caseMap)
}

func evalExpr(t *testing.T, caseMap map[string][]token.TokenType) {
	for inputExp, expectExp := range caseMap {
		// fmt.Println(inputExp)
		lex := Init(inputExp)
		var ttArray []token.TokenType
		for tkn := range lex.Consume() {
			ttArray = append(ttArray, tkn.Type)
			token.Print(tkn)
		}
		if !reflect.DeepEqual(ttArray, expectExp) {
			t.Errorf("input: %s,tokenType %+v is %+v", inputExp, ttArray, expectExp)
		}
	}
}

func evalCase(t *testing.T, caseMap map[string]token.TokenType) {
	for inputExp, expectExp := range caseMap {
		// fmt.Println(inputExp, expectExp)

		lex := Init(inputExp)
		for tkn := range lex.Consume() {
			// NOTE: i don't like this
			if tkn.Type != expectExp {
				t.Errorf("tokenType %s is %d", inputExp, tkn.Type)
			}
		}
	}
}
