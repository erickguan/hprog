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
func isTheEnd(ch rune) bool       { return (ch == token.EoF || ch == ' ') }

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

type Lexer struct {
	input    string
	position int
	line     int
	start    int
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
	// lex.start = lex.position
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
	fmt.Println("start:", lex.start, "end:", lex.position, "len:", len(lex.input))
	fmt.Println("Value:", lex.input[lex.start:lex.position])
	lex.tokens <- token.Token{
		Type:     tokenType,
		Position: lex.position,
		Line:     lex.line,
		Value:    lex.input[lex.start:lex.position],
	}
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

func (lex *Lexer) scanDigit() bool {
	lex.unread()
	lex.start = lex.position

	/* ACCEPT DIGITS */
	digits := "0123456789"
	lex.acceptRun(digits)

	dot := "."
	/* ACCEPT DIGITS.DIGITS */
	if lex.accept(dot) {
		lex.acceptRun(digits)
	}

	if isTheEnd(lex.peek()) {
	} else if IsAlphaNumeric(lex.peek()) {
		/* ERROR DIGITS.DIGITS|ALPHA */
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
		/* ERROR ALPHA | DIGIT | NON-ALPHA*/
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

func (lex *Lexer) extractString() bool {
	// good luck
	return true
}

func fullScan(lex *Lexer) stateFunc {
	for {
		ch := lex.read()

		switch ch1 := ch; {
		case IsDigit(ch1):
			// TODO: all the cases for the digilex
			// probably a dynamic value creation based
			// on what type it is.
			done := lex.scanDigit()
			if !done {
				lex.emit(token.ERR)
				lex.reportError("SyntaxError, number malformed.")
				return nil
			}
			lex.emit(token.NUMBER)
		case IsLetter(ch):
			lex.unread()
			done := lex.scanIdentifier()
			// if is reseved (error on assign)
			if !done {
				lex.reportError("SyntaxError, indentifier malformed.")
				lex.emit(token.ERR)
				return nil
			}
			detectedType := lex.identifierToReseved(token.IDENTIFIER)
			lex.emit(detectedType)
		default:
			switch ch {
			case ' ':
				lex.trimWhitespace()
			case '\n':
				// TODO: only temporary
				lex.line += 1
			case '#':
				lex.skipComment()
				lex.emit(token.COMMENT)
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
					lex.reportError("SyntaxError, number misformed.")
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
				fmt.Println("AA", token.ReversedTokenMap[rtoken])
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
					return nil
				}
			case token.EoF:
				lex.emit(token.EOF)
			default:
				lex.reportError("Token not recognized!")
				return nil
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
