package lexer

import (
	"reflect"
	"testing"

	"github.com/badc0re/hprog/token"
)

func TestParsingNumbers(t *testing.T) {
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
	evluateTestCases(t, testCases)
}

func TestParsingIdentifiers(t *testing.T) {
	var testCases = map[string]token.TokenType{
		"a11":   token.IDENTIFIER,
		"AAA":   token.IDENTIFIER,
		"AA!":   token.ERR,
		"AA1.2": token.ERR,
	}
	evluateTestCases(t, testCases)
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
	/*
			caseMap := map[string][]token.TokenType{
				"1 + 2":   []token.TokenType{token.NUMBER, token.PLUS, token.NUMBER},
				"1.2 + 3": []token.TokenType{token.NUMBER, token.PLUS, token.NUMBER},
					"((1 + 2) - 3)": []token.TokenType{token.OP, token.OP, token.NUMBER, token.PLUS, token.NUMBER, token.CP, token.MINUS, token.NUMBER, token.CP},
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
		evaluateExpression(t, caseMap)
	*/
}

func evaluateExpression(t *testing.T, caseMap map[string][]token.TokenType) {
	for inputExp, expectExp := range caseMap {
		// fmt.Println(inputExp)
		lex := Init(inputExp)
		var ttArray []token.TokenType
		for tkn, done := lex.Consume(); done != false; {
			ttArray = append(ttArray, tkn.Type)
			// token.Print(tkn)
		}
		if !reflect.DeepEqual(ttArray, expectExp) {
			t.Errorf("input: %s, ouput: %+v, expected: %+v", inputExp, ttArray, expectExp)
		}
	}
}

func evluateTestCases(t *testing.T, caseMap map[string]token.TokenType) {
	for input, expected := range caseMap {

		lex := Init(input)
		for {
			// TODO: this should be an array...
			tkn, done := lex.Consume()
			if done == true || tkn.Type == token.EOF {
				break
			}
			if tkn.Value != input {
				t.Errorf("input: %s, output: %s", input, tkn.Value)
			}
			if tkn.Type != expected {
				t.Errorf("input type: %s, output type: %s", token.ReversedTokenMap[expected], token.ReversedTokenMap[tkn.Type])
			}
		}
	}
}
