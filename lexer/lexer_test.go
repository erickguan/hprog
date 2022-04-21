package lexer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/badc0re/hprog/token"
)

func TestLexerNumbers(t *testing.T) {
	var testCases = map[string]token.TokenType{
		"11a":   token.ERR,
		"11.a0": token.ERR,
		"11a0":  token.ERR,
		"11":    token.NUMBER,
		"11.":   token.NUMBER,
		".11":   token.NUMBER,
		"1.0":   token.NUMBER,
		"a11":   token.IDENTIFIER,
	}
	evaluateExpression1(t, testCases)
}

func TestLexerIdentifiers(t *testing.T) {
	var testCases = map[string]token.TokenType{
		"a11":   token.IDENTIFIER,
		"a11a":  token.IDENTIFIER,
		"AAA":   token.IDENTIFIER,
		"a11 ":  token.IDENTIFIER,
		"a11a ": token.IDENTIFIER,
		"AAA ":  token.IDENTIFIER,
		"11aa":  token.ERR,
		//"AA!":     token.ERR,
		//"AA1.2":   token.ERR,
		//"(AA1.2)": token.ERR,
		//"AA!=":    token.ERR,
	}
	evaluateExpression1(t, testCases)
}

func TestLexerExpression1(t *testing.T) {
	caseMap := map[string][]token.TokenType{
		"(True)": []token.TokenType{token.OP, token.BOOL_TRUE, token.CP},
		"(a)":    []token.TokenType{token.OP, token.IDENTIFIER, token.CP},
	}
	evaluateExpression(t, caseMap)
}

func TestLexerExpression(t *testing.T) {
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
		"False == True":                  []token.TokenType{token.BOOL_FALSE, token.EQUAL_EQUAL, token.BOOL_TRUE},
		"(False == True)":                []token.TokenType{token.OP, token.BOOL_FALSE, token.EQUAL_EQUAL, token.BOOL_TRUE, token.CP},
		"decl b = 10; # (if equal True)": []token.TokenType{token.DECLARE, token.IDENTIFIER, token.EQUAL, token.NUMBER},
		"decl a == 123":                  []token.TokenType{token.DECLARE, token.IDENTIFIER, token.EQUAL_EQUAL, token.NUMBER},
	}
	evaluateExpression(t, caseMap)
}

func TestLexerString(t *testing.T) {
	var caseMap = map[string]token.TokenType{
		/* not supported
		"\"test\"": token.STRING,
		// "\"test":   token.ERR,
		*/
		"\"test\"": token.STRING,
	}
	evaluateExpression1(t, caseMap)
}

func evaluateExpression(t *testing.T, caseMap map[string][]token.TokenType) {
	for inputExp, expectExp := range caseMap {
		lex := Init(inputExp)
		var ttArray []token.TokenType

		for {
			tkn, done := lex.Consume()
			if done == true || tkn.Type == token.EOF {
				break
			}
			ttArray = append(ttArray, tkn.Type)
		}

		fmt.Println(ttArray)

		if !reflect.DeepEqual(ttArray, expectExp) {
			t.Errorf("input: %s, ouput: %+v, expected: %+v", inputExp, ttArray, expectExp)
		}
	}
}

func evaluateExpression1(t *testing.T, caseMap map[string]token.TokenType) {
	for input, expected := range caseMap {
		lex := Init(input)
		for {
			tkn, done := lex.Consume()
			//  Skip these tokens.
			if done == true || tkn.Type == token.EOF {
				break
			}
			if tkn.Type != expected {
				t.Errorf("input %s, input type: %s, output type: %s", input, token.ReversedTokenMap[expected], token.ReversedTokenMap[tkn.Type])
			}
		}
	}
}
