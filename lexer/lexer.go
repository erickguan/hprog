package lexer

import (
	"bytes"
	"fmt"
	"hprog/errors"
	"hprog/token"
	"io"
	"os"
	"strings"
	"unicode"
)

func IsDigit(ch rune) bool { return unicode.IsDigit(ch) }

func IsLetter(ch rune) bool { return unicode.IsLetter(ch) }

func IsAlphaNumeric(ch rune) bool { return (IsLetter(ch) || IsDigit(ch)) }

func reportError(ts *TokenScanner, what string) {
	fmt.Fprintf(os.Stderr, "[line:%d, pos:%d] Error, %s\n",
		ts.Line, ts.Position, what)
}

type TokenScanner struct {
	Reader   io.RuneScanner
	Position int
	Line     int
	buf      bytes.Buffer
}

func (ts *TokenScanner) unread() {
	ts.Reader.UnreadRune()
	ts.Position--
}

func (ts *TokenScanner) read() (rune, error) {
	ch, _, err := ts.Reader.ReadRune()
	ts.Position++
	return ch, err
}

func (ts *TokenScanner) peek() (rune, error) {
	ch, _, err := ts.Reader.ReadRune()
	ts.Reader.UnreadRune()
	return ch, err
}

type Lexer struct {
	Scanner *TokenScanner
	tokens  chan token.Token
	STOP    bool
}

type stateFunc func(*Lexer) stateFunc

func (lex *Lexer) trimWhitespace() {
	for {
		ch, _ := lex.Scanner.peek()
		if ch == ' ' {
			lex.Scanner.read()
		} else {
			break
		}
	}
}

func (lex *Lexer) Consume() chan token.Token {
	return lex.tokens
}

func (lex *Lexer) skipComment() {
	for {
		ch, err := lex.Scanner.read()
		if err == io.EOF || ch == '\n' {
			break
		} else {
			lex.Scanner.read()
		}
	}
}

func (lex *Lexer) emit(tokenType token.TokenType) {
	tkn := token.Token{
		Type:     tokenType,
		Position: lex.Scanner.Position,
		Line:     lex.Scanner.Line,
		Value:    lex.Scanner.buf.String(),
	}
	if lex.Scanner.buf.Len() != 0 {
		//fmt.Print(lex.Scanner.buf.String(), " ")
	}
	// b is temporary
	token.Print(tkn)
	lex.tokens <- tkn
	/*
		lex.tokens <- token.Token{
			Type:     tokenType,
			Position: lex.Scanner.Position,
			Line:     lex.Scanner.Line,
		}
	*/
	lex.Scanner.buf.Reset()
	lex.Scanner.Position = 1
}

func (lex *Lexer) scanDigit() error {
	lex.Scanner.unread()
	for {
		ch, _ := lex.Scanner.peek()
		if IsDigit(ch) {
			ch, _ := lex.Scanner.read()
			lex.Scanner.buf.WriteRune(ch)
		} else if IsLetter(ch) {
			return errors.NewSyntaxError("Syntax error, mixing digits with characters.")
		} else {
			break
		}
	}
	return nil
}

func (lex *Lexer) scanIdentifier() (token.TokenType, error) {
	lex.Scanner.unread()
	for {
		ch, _ := lex.Scanner.peek()
		/*
			if IsLetter(ch) {
				ch, _ := lex.Scanner.read()
				lex.Scanner.buf.WriteRune(ch)
			} else if IsDigit(ch) {
				lex.Scanner.buf.Reset()
				return token.ERR, errors.NewSyntaxError("Syntax error, identifier mixed with numbers.")
			} else {
				break
			}
		*/
		// to allow int8
		if IsLetter(ch) || IsDigit(ch) {
			ch, _ := lex.Scanner.read()
			lex.Scanner.buf.WriteRune(ch)
		} else {
			break
		}
	}
	return token.IDENTIFIER, nil
}

func (lex *Lexer) identifierToReseved(ttype token.TokenType) token.TokenType {
	resevedToken := token.TokenMap[lex.Scanner.buf.String()]
	if resevedToken != 0 {
		return resevedToken
	}
	return ttype
}

func (lex *Lexer) scanConditions(rcurrent token.TokenType, rfuture token.TokenType) token.TokenType {
	ch, _ := lex.Scanner.peek()
	if ch == '=' {
		return rfuture
	}
	return rcurrent
}

func (lex *Lexer) extractString() bool {
	// good luck
	return true
}

func fullScan(lex *Lexer) stateFunc {
	// NOTE: i don't like this
	// can it be recrusive?
loop:
	for {
		ch, err := lex.Scanner.read()
		if err == io.EOF {
			break loop
		}

		switch ch1 := ch; {
		case IsDigit(ch1):
			// TODO: all the cases for the digits
			// probably a dynamic value creation based
			// on what type it is.
			err := lex.scanDigit()
			if err != nil {
				lex.emit(token.ERR)
				reportError(lex.Scanner, err.Error())
				return nil
			}
			lex.emit(token.NUMBER)
		case IsLetter(ch):
			ttype, err := lex.scanIdentifier()
			// if is reseved (error on assign)
			if err != nil {
				lex.emit(token.ERR)
				reportError(lex.Scanner, err.Error())
				return nil
			}
			ttype = lex.identifierToReseved(ttype)
			lex.emit(ttype)
		default:
			switch ch {
			case ' ':
				lex.trimWhitespace()
			case '\n':
				// TODO: only temporary
				lex.Scanner.Line += 1
			case '#':
				lex.emit(token.COMMENT)
				lex.skipComment()
			case '!':
				// TODO: is it a condition first
				rtoken := lex.scanConditions(token.EXCL, token.EXCL_EQUAL)
				lex.emit(rtoken)
			case '=':
				// TODO: is it a condition first
				rtoken := lex.scanConditions(token.ASSIGN, token.EQUAL_EQUAL)
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
					reportError(lex.Scanner, "Wrong string formatting.")
					lex.emit(token.ERR)
				}
			default:
				token.TokenMap
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

func Init(expression string) Lexer {
	lex := Lexer{
		Scanner: &TokenScanner{
			Reader:   strings.NewReader(expression),
			Position: 1,
			Line:     1,
		},
		tokens: make(chan token.Token),
	}

	go lex.run()
	return lex
}
