package parser

import (
	"fmt"
	"os"

	"github.com/badc0re/hprog/chunk"
	"github.com/badc0re/hprog/codes"
	"github.com/badc0re/hprog/lexer"
	"github.com/badc0re/hprog/token"
	"github.com/badc0re/hprog/value"
)

type PREC int

const (
	PREC_ILLEGAL PREC = iota
	PREC_NONE
	PREC_ASSIGN    // =
	PREC_OR        // or
	PREC_AND       // and
	PREC_EQUALLITY // ==, !=
	PREC_COMPARE   // <, >, <=, >=
	PREC_TERM      // +, -
	PREC_FACTOR    // *, /
	PREC_UNARY     // !, -
	PREC_CALL      // ., ()
	PREC_PRIMARY
)

type Parser struct {
	current  *token.Token
	previous *token.Token
	lex      *lexer.Lexer
	Perror   bool
	ppanic   bool
	tknMap   map[token.TokenType]ParseRule

	// todo
	chk *chunk.Chunk
}

type ParseFn func()

type ParseRule struct {
	prefix ParseFn
	infix  ParseFn
	prec   PREC
}

func (p *Parser) getRule(tknType token.TokenType) ParseRule {
	return p.tknMap[tknType]
}

func (p *Parser) Consume(tknType token.TokenType, message string) {
	if p.current.Type == tknType {
		p.Advance()
		return
	}

	p.reportError(p.current, message)
}

func (p *Parser) ParsePrec(prec PREC) {
	p.Advance()
	prefRule := p.getRule(p.previous.Type).prefix
	if prefRule == nil {
		return
	}

	prefRule()

	for {
		prec1 := p.getRule(p.current.Type).prec
		if prec >= prec1 {
			break
		}
		p.Advance()
		infix := p.getRule(p.previous.Type).infix
		infix()
	}
}

func (p *Parser) Unary() {
	tknType := p.previous.Type

	p.ParsePrec(PREC_UNARY)

	switch tknType {
	case token.MINUS:
		p.emit(codes.INSTRUC_NEGATE)
	case token.EXCL:
		p.emit(codes.INSTRUC_NOT)
	default:
		return
	}
}

func (p *Parser) Binary() {
	tknType := p.previous.Type
	rule := p.getRule(tknType)
	p.ParsePrec(rule.prec + 1)

	switch tknType {
	case token.PLUS:
		p.emit(codes.INSTRUC_ADDITION)
	case token.MINUS:
		p.emit(codes.INSTRUC_SUBSTRACT)
	case token.STAR:
		p.emit(codes.INSTRUC_MULTIPLY)
	case token.SLASH:
		p.emit(codes.INSTRUC_DIVIDE)
	case token.EQUAL_EQUAL:
		p.emit(codes.INSTRUC_EQUAL)
	case token.EXCL_EQUAL:
		p.emit2(codes.INSTRUC_EQUAL, codes.INSTRUC_NOT)
	case token.GREATER:
		p.emit(codes.INSTRUC_GREATER)
	case token.GREATER_EQUAL:
		p.emit2(codes.INSTRUC_LESS, codes.INSTRUC_NOT)
	case token.LESS:
		p.emit(codes.INSTRUC_LESS)
	case token.LESS_EQUAL:
		p.emit2(codes.INSTRUC_GREATER, codes.INSTRUC_NOT)
	default:
		return
	}
}

func (p *Parser) emit(code interface{}) {
	p.chk.WriteChunk(code, p.previous.Line)
}

func (p *Parser) emit2(code1 interface{}, code2 interface{}) {
	p.chk.WriteChunk(code1, p.previous.Line)
	p.chk.WriteChunk(code2, p.previous.Line)
}

func (p *Parser) EndCompile() {
	p.emitReturn()
}

func (p *Parser) emitReturn() {
	p.emit(codes.INSTRUC_RETURN)
}

func (p *Parser) emitVariable(v value.Value) {
	p.emit2(codes.INSTRUC_CONSTANT, p.makeConstant(v))
}

func (p *Parser) makeConstant(v value.Value) uint {
	return p.chk.AddVariable(v)
}

func (p *Parser) Number() {
	dt := value.DetectNumberTypeByConversion(p.previous.Value)
	p.emitVariable(value.New(p.previous.Value, dt))
}

func (p *Parser) String() {
}

func (p *Parser) Literal() {
	tokenType := p.previous.Type
	switch tokenType {
	case token.BOOL_FALSE:
		p.emit(codes.INSTRUC_FALSE)
	case token.BOOL_TRUE:
		p.emit(codes.INSTRUC_TRUE)
	case token.NIL:
		p.emit(codes.INSTRUC_NIL)
	default:
		return
	}
}

func (p *Parser) Grouping() {
	p.Expression()
	p.Consume(token.CP, "Expected ')' after expression.")
}

func (p *Parser) Expression() {
	p.ParsePrec(PREC_ASSIGN)
}

func (p *Parser) Match(tokenType token.TokenType) bool {
	if !p.Check(tokenType) {
		return false
	}
	p.Advance()
	return true
}

func (p *Parser) Decl() {
	p.Statement()
}

func (p *Parser) Statement() {
	if p.Match(token.PRINT) {
		p.PrintStmt()
	}
}

func (p *Parser) Check(tokenType token.TokenType) bool {
	return p.current.Type == tokenType
}

func (p *Parser) PrintStmt() {
	p.Consume(token.OP, "Expected '(' after expression.")
	p.Grouping()
	p.emit(codes.INSTRUC_PRINT)
}

func (p *Parser) Advance() {
	p.previous = p.current

	tkn, done := p.lex.Consume()
	if done {
		return
	}
	p.current = tkn
}

func (p *Parser) reportError(tkn *token.Token, what string) {
	p.ppanic = true
	fmt.Fprintf(os.Stderr, "[line:%d, pos:%d] Error, %s\n",
		tkn.Line, tkn.Position, what)
	p.Perror = true
}

func Init(lex *lexer.Lexer, chk *chunk.Chunk) *Parser {
	p := Parser{
		lex: lex,
		chk: chk,
	}
	tknMap := map[token.TokenType]ParseRule{
		token.OP:            {p.Grouping, nil, PREC_NONE},
		token.CP:            {nil, nil, PREC_NONE},
		token.LB:            {nil, nil, PREC_NONE},
		token.RB:            {nil, nil, PREC_NONE},
		token.COMMA:         {nil, nil, PREC_NONE},
		token.DOT:           {nil, nil, PREC_NONE},
		token.MINUS:         {p.Unary, p.Binary, PREC_TERM},
		token.PLUS:          {nil, p.Binary, PREC_TERM},
		token.SEMICOLON:     {nil, nil, PREC_NONE},
		token.SLASH:         {nil, p.Binary, PREC_FACTOR},
		token.STAR:          {nil, p.Binary, PREC_FACTOR},
		token.EXCL:          {p.Unary, nil, PREC_TERM},
		token.EXCL_EQUAL:    {nil, p.Binary, PREC_EQUALLITY},
		token.EQUAL:         {nil, nil, PREC_NONE},
		token.EQUAL_EQUAL:   {nil, p.Binary, PREC_COMPARE},
		token.GREATER:       {nil, p.Binary, PREC_COMPARE},
		token.GREATER_EQUAL: {nil, p.Binary, PREC_COMPARE},
		token.LESS:          {nil, p.Binary, PREC_COMPARE},
		token.LESS_EQUAL:    {nil, p.Binary, PREC_COMPARE},
		token.IDENTIFIER:    {nil, nil, PREC_NONE},
		token.STRING:        {p.String, nil, PREC_NONE},
		token.NUMBER:        {p.Number, nil, PREC_NONE},
		token.AND:           {nil, nil, PREC_NONE},
		token.ELSE:          {nil, nil, PREC_NONE},
		token.BOOL_FALSE:    {p.Literal, nil, PREC_NONE},
		token.BOOL_TRUE:     {p.Literal, nil, PREC_NONE},
		token.FOR:           {nil, nil, PREC_NONE},
		token.FUNCTION:      {nil, nil, PREC_NONE},
		token.IF:            {nil, nil, PREC_NONE},
		token.NIL:           {p.Literal, nil, PREC_NONE},
		token.OR:            {nil, nil, PREC_NONE},
		token.PRINT:         {nil, nil, PREC_NONE},
		token.RETURN:        {nil, nil, PREC_NONE},
		token.VAR:           {nil, nil, PREC_NONE},
		token.WHILE:         {nil, nil, PREC_NONE},
		token.ERR:           {nil, nil, PREC_NONE},
		token.EOF:           {nil, nil, PREC_NONE},
	}
	/*
		Init MUST return a reference, otherwise
		the functions would not point to the correct
		address.

		HEAP HEAP HEAP...hurray?
	*/
	p.tknMap = tknMap
	return &p
}
