package lexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/badc0re/hprog/token"
)

type Lexer struct {
	input        string
	position     int
	line         int
	start        int
	tokens       chan token.Token
	requiresSemi bool
}

type stateFunc func(*Lexer) stateFunc

func IsDigit(ch rune) bool        { return unicode.IsDigit(ch) }
func IsLetter(ch rune) bool       { return unicode.IsLetter(ch) }
func IsAlphaNumeric(ch rune) bool { return (IsLetter(ch) || IsDigit(ch)) }

func (lex *Lexer) reportError(reason string) {
	fmt.Fprintf(os.Stderr, "[line:%d, pos:%d], %s\n",
		lex.line, lex.position, reason)
}

func (lex *Lexer) unread() {
	lex.position--
}

func (lex *Lexer) read() rune {
	if lex.position >= len(lex.input) {
		return token.EoF
	}
	ch, _ := utf8.DecodeRuneInString(lex.input[lex.position:])
	lex.position++
	return ch
}

func (lex *Lexer) peek() rune {
	if lex.position >= len(lex.input) {
		return token.EoF
	}
	ch, _ := utf8.DecodeRuneInString(lex.input[lex.position:])
	return ch
}

func (lex *Lexer) setRequiresSemi(required bool) {
	fmt.Println(required)
	lex.requiresSemi = required
}

func (lex *Lexer) trimWhitespace() {
	_trim := " "
	lex.acceptRun(_trim)
	lex.start = lex.position
}

func (lex *Lexer) trimNewline() {
	_trim := "\n"
	lex.acceptRun(_trim)
	lex.start = lex.position
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
		ch := lex.peek()
		if ch == '\n' || ch == token.EoF {
			break
		} else {
			lex.read()
		}
	}
}

func (lex *Lexer) emit(tokenType token.TokenType) {
	tkn := token.Token{
		Type:     tokenType,
		Position: lex.position,
		Line:     lex.line,
		Value:    lex.input[lex.start:lex.position],
	}
	lex.tokens <- tkn
	lex.start = lex.position
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

func (lex *Lexer) scanNumber() bool {
	lex.unread()
	lex.start = lex.position

	/* ACCEPT DIGITS */
	// token.INIT

	digits := "0123456789"
	lex.acceptRun(digits)

	dot := "."
	/* ACCEPT DIGITS.DIGITS */
	if lex.accept(dot) {
		// token.FLOAT
		lex.acceptRun(digits)
	}

	if IsAlphaNumeric(lex.peek()) {
		return false
	}
	return true
}

func (lex *Lexer) scanIdentifier() bool {
	lex.start = lex.position
	/* ACCEPT ^ALPHA */
	if !IsLetter(lex.peek()) {
		return false
	}
	/* ACCEPT ALPHA | DIGIT */
	for IsLetter(lex.peek()) || IsDigit(lex.peek()) {
		lex.read()
	}
	if IsAlphaNumeric(lex.peek()) {
		return false
	}
	return true
}

func (lex *Lexer) identifierToReseved(defaultType token.TokenType) token.TokenType {
	reservedToken := token.TokenMap[lex.input[lex.start:lex.position]]
	if reservedToken != 0 {
		return reservedToken
	}
	return defaultType
}

func (lex *Lexer) scanConditions(rcurrent token.TokenType, rfuture token.TokenType) token.TokenType {
	ch := lex.peek()
	if ch == '=' {
		lex.read()
		return rfuture
	}
	return rcurrent
}

func (lex *Lexer) scanString() bool {
	lex.start = lex.position
	for {
		ch := lex.read()
		if ch == '"' {
			// don't consume '"'
			lex.unread()
			break
		}
		if ch == '\n' || ch == token.EoF {
			return false
		}
	}
	return true
}

func fullScan(lex *Lexer) stateFunc {
	for {
		ch := lex.read()

		switch ch1 := ch; {
		case IsDigit(ch1):
			done := lex.scanNumber()
			if !done {
				lex.emit(token.ERR)
				lex.reportError("SyntaxError, number malformed.")
				return nil
			}
			lex.emit(token.NUMBER)
		case IsLetter(ch):
			lex.unread()
			done := lex.scanIdentifier()
			if !done {
				lex.reportError("SyntaxError, indentifier malformed.")
				lex.emit(token.ERR)
				return nil
			}
			detectedType := lex.identifierToReseved(token.IDENTIFIER)
			lex.setRequiresSemi(true)
			lex.emit(detectedType)
		default:
			switch ch {
			case ' ':
				lex.trimWhitespace()
			case '\n':
				lex.line += 1
				if lex.requiresSemi == true {
					fmt.Println("AAA")
					lex.emit(token.SEMICOLON)
				}
				lex.trimNewline()
				lex.setRequiresSemi(false)
			case '#':
				lex.skipComment()
			case '+':
				lex.emit(token.PLUS)
				lex.setRequiresSemi(true)
			case '-':
				lex.emit(token.MINUS)
				lex.setRequiresSemi(true)
			case '/':
				lex.emit(token.SLASH)
				lex.setRequiresSemi(true)
			case '*':
				lex.emit(token.STAR)
				lex.setRequiresSemi(true)
			case '(':
				lex.emit(token.OP)
				lex.setRequiresSemi(false)
			case ')':
				lex.emit(token.CP)
				lex.setRequiresSemi(true)
			case '{':
				lex.emit(token.LB)
			case '}':
				lex.emit(token.RB)
			case ',':
				lex.emit(token.COMMA)
			case '.':
				done := lex.scanNumber()
				if !done {
					lex.reportError("SyntaxError, number malformed.")
					lex.emit(token.ERR)
					return nil
				}
				lex.emit(token.NUMBER)
			case ';':
				// TODO: is it needed?
				// lex.emit(token.SEMICOLON)
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
			case '\'':
				lex.emit(token.SINGLE_QUOTE)
			case '"':
				if lex.scanString() {
					lex.emit(token.STRING)
					// consume the trailing '"'
					lex.read()
				} else {
					lex.reportError("Wrong string formatting.")
					lex.emit(token.ERR)
					return nil
				}
			case token.EoF:
				lex.emit(token.EOF)
				return nil
			default:
				lex.emit(token.ERR)
				lex.reportError("Token not recognized.")
				return nil
			}
		}
	}
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
