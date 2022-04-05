package lexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/badc0re/hprog/token"
)

func IsDigit(ch rune) bool { return unicode.IsDigit(ch) }

func IsLetter(ch rune) bool { return unicode.IsLetter(ch) }

func IsAlphaNumeric(ch rune) bool { return (IsLetter(ch) || IsDigit(ch)) }

func (lex *Lexer) reportError(reason string) {
	fmt.Fprintf(os.Stderr, "[line:%d, pos:%d], %s\n",
		lex.line, lex.position, reason)
}

func (lex *Lexer) unread() {
	lex.position--
}

func (lex *Lexer) read() rune {
	// ERROR?
	if lex.position == len(lex.input) {
		return token.EoF
	}
	ch, _ := utf8.DecodeRuneInString(lex.input[lex.position:])
	lex.position++
	fmt.Println("read():", string(ch))
	return ch
}

func (lex *Lexer) peek() rune {
	ch, _ := utf8.DecodeRuneInString(lex.input[lex.position:])
	fmt.Println("peek():", string(ch))
	return ch
}

type Lexer struct {
	input    string
	position int
	line     int
	start    int
	end      int
	tokens   chan token.Token
}

type stateFunc func(*Lexer) stateFunc

func (lex *Lexer) trimWhitespace() {
	for {
		ch := lex.peek()
		if ch == ' ' {
			lex.read()
		} else {
			break
		}
	}
}

func (lex *Lexer) Consume() (*token.Token, bool) {
	if tkn, ok := <-lex.tokens; ok {
		return &tkn, false
	} else {
		return nil, true
	}
}

func (lex *Lexer) skipComment() {
	for {
		ch := lex.read()
		if ch == '\n' || ch == token.EoF {
			break
		} else {
			lex.read()
		}
	}
}

func (lex *Lexer) emit(tokenType token.TokenType) {
	lex.tokens <- token.Token{
		Type:     tokenType,
		Position: lex.position,
		Line:     lex.line,
		Value:    lex.input[lex.start:lex.end],
	}
	lex.start = lex.position
	lex.end = lex.position
}

func (lex *Lexer) accept(v string) bool {
	if strings.ContainsRune(v, lex.peek()) {
		lex.read()
		return true
	}
	return false
}

func (lex *Lexer) acceptRun(v string) {
	for strings.ContainsRune(v, lex.peek()) {
		lex.read()
	}
}

func (lex *Lexer) scanDigit() bool {
	lex.unread()
	lex.start = lex.position

	digits := "0123456789"
	lex.acceptRun(digits)

	dot := "."
	if lex.accept(dot) {
		lex.acceptRun(digits)
	}

	/*
		if IsAlphaNumeric(lex.read()) {
			return false
		}
	*/

	lex.end = lex.position
	return true
}

func (lex *Lexer) scanIdentifier() bool {
	if !IsLetter(lex.peek()) {
		return false
	}
	for IsLetter(lex.peek()) || IsDigit(lex.peek()) {
		lex.read()
	}
	/*
		if !IsAlphaNumeric(lex.read()) {
			return false
		}
	*/
	return true
}

/*
func (lex *Lexer) identifierToReseved(ttype token.TokenType) token.TokenType {
	resevedToken := token.TokenMap[lex.Scanner.buf.String()]
	if resevedToken != 0 {
		return resevedToken
	}
	return ttype
}
*/

func (lex *Lexer) scanConditions(rcurrent token.TokenType, rfuture token.TokenType) token.TokenType {
	ch := lex.peek()
	if ch == '=' {
		lex.read()
		return rfuture
	}
	return rcurrent
}

func (lex *Lexer) extractString() bool {
	// good luck
	return true
}

func fullScan(lex *Lexer) stateFunc {
loop:
	for {
		ch := lex.read()
		if ch == token.EoF {
			lex.emit(token.EOF)
			break loop
		}

		switch ch1 := ch; {
		case IsDigit(ch1):
			// TODO: all the cases for the digilex
			// probably a dynamic value creation based
			// on what type it is.
			done := lex.scanDigit()
			fmt.Println(done)
			if !done {
				lex.emit(token.ERR)
				lex.reportError("SyntaxError")
				return nil
			}
			lex.emit(token.NUMBER)
		case IsLetter(ch):
			lex.unread()
			fmt.Println("Ident")
			done := lex.scanIdentifier()
			// if is reseved (error on assign)
			if !done {
				fmt.Println("ERR")
				lex.emit(token.ERR)
				return nil
			}
			//ttype = lex.identifierToReseved(ttype)
			lex.emit(token.IDENTIFIER)
		default:
			switch ch {
			case ' ':
				lex.trimWhitespace()
			case '\n':
				// TODO: only temporary
				lex.line += 1
			case '#':
				lex.emit(token.COMMENT)
				lex.skipComment()
			case '+':
				lex.emit(token.PLUS)
			case '-':
				lex.emit(token.MINUS)
			case '/':
				lex.emit(token.SLASH)
			case '*':
				lex.emit(token.STAR)
			case '(':
				lex.emit(token.OP)
			case ')':
				lex.emit(token.CP)
			case '{':
				lex.emit(token.LB)
			case '}':
				lex.emit(token.RB)
			case ',':
				lex.emit(token.COMMA)
			case '.':
				done := lex.scanDigit()
				if !done {
					lex.reportError("SyntaxError")
					lex.emit(token.ERR)
					return nil
				}
				lex.emit(token.NUMBER)
			case ';':
				lex.emit(token.SEMICOLON)
			case ':':
				lex.emit(token.COLON)
			case '!':
				// TODO: is it a condition first
				rtoken := lex.scanConditions(token.EXCL, token.EXCL_EQUAL)
				lex.emit(rtoken)
			case '=':
				// TODO: is it a condition first
				rtoken := lex.scanConditions(token.EQUAL, token.EQUAL_EQUAL)
				lex.emit(rtoken)
			case '<':
				// TODO: is it a condition first
				rtoken := lex.scanConditions(token.LESS, token.LESS_EQUAL)
				lex.emit(rtoken)
			case '>':
				// TODO: is it a condition first
				rtoken := lex.scanConditions(token.GREATER, token.GREATER_EQUAL)
				lex.emit(rtoken)
			case '"':
				if lex.extractString() {
					lex.emit(token.STRING)
				} else {
					lex.reportError("Wrong string formatting.")
					lex.emit(token.ERR)
				}
			default:
				break
				lex.reportError("Token not recognized!")
			}
		}
	}
	return nil
}

func (lex *Lexer) run() {
	for state := fullScan; state != nil; {
		state = state(lex)
	}
	close(lex.tokens)
}

func Init(expression string) *Lexer {
	lex := Lexer{
		input:    expression,
		position: 0,
		line:     1,
		tokens:   make(chan token.Token),
	}

	go lex.run()
	return &lex
}
