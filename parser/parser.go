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

type Compiler struct {
	Locals     []Local
	localCount int
	scopeDepth int
}

type Local struct {
	Name  token.Token
	Depth int
}

type Parser struct {
	current     *token.Token
	previous    *token.Token
	lex         *lexer.Lexer
	Perror      bool
	ppanic      bool
	tknMap      map[token.TokenType]ParseRule
	currentComp *Compiler

	// todo
	chk *chunk.Chunk
}

var tknMap = map[token.TokenType]ParseRule{
	token.OP:            {Grouping, fcall, PREC_CALL},
	token.CP:            {nil, nil, PREC_NONE},
	token.LB:            {nil, nil, PREC_NONE},
	token.RB:            {nil, nil, PREC_NONE},
	token.COMMA:         {nil, nil, PREC_NONE},
	token.DOT:           {nil, nil, PREC_NONE},
	token.MINUS:         {Unary, Binary, PREC_TERM},
	token.PLUS:          {nil, Binary, PREC_TERM},
	token.SEMICOLON:     {nil, nil, PREC_NONE},
	token.SLASH:         {nil, Binary, PREC_FACTOR},
	token.STAR:          {nil, Binary, PREC_FACTOR},
	token.EXCL:          {Unary, nil, PREC_TERM},
	token.EXCL_EQUAL:    {nil, Binary, PREC_EQUALLITY},
	token.EQUAL:         {nil, nil, PREC_NONE},
	token.EQUAL_EQUAL:   {nil, Binary, PREC_COMPARE},
	token.GREATER:       {nil, Binary, PREC_COMPARE},
	token.GREATER_EQUAL: {nil, Binary, PREC_COMPARE},
	token.LESS:          {nil, Binary, PREC_COMPARE},
	token.LESS_EQUAL:    {nil, Binary, PREC_COMPARE},
	token.STRING:        {String, nil, PREC_NONE},
	token.NUMBER:        {Number, nil, PREC_NONE},
	token.AND:           {nil, nil, PREC_NONE},
	token.ELSE:          {nil, nil, PREC_NONE},
	token.BOOL_FALSE:    {Literal, nil, PREC_NONE},
	token.BOOL_TRUE:     {Literal, nil, PREC_NONE},
	token.FOR:           {nil, nil, PREC_NONE},
	token.FUNCTION:      {nil, nil, PREC_NONE},
	token.IF:            {nil, nil, PREC_NONE},
	// maybe not
	token.OR:         {nil, nil, PREC_NONE},
	token.NIL:        {Literal, nil, PREC_NONE},
	token.PRINT:      {nil, nil, PREC_NONE},
	token.RETURN:     {nil, nil, PREC_NONE},
	token.IDENTIFIER: {Variable, nil, PREC_NONE},
	token.WHILE:      {nil, nil, PREC_NONE},
	token.ERR:        {nil, nil, PREC_NONE},
	token.EOF:        {nil, nil, PREC_NONE},
}

type ParseFn func(*Parser, bool)

type ParseRule struct {
	prefix ParseFn
	infix  ParseFn
	prec   PREC
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

func (p *Parser) getRule(tknType token.TokenType) ParseRule {
	rule, found := p.tknMap[tknType]
	if found == false {
		p.reportError(p.current, "Expression not supported")
	}
	return rule
}

func (p *Parser) Consume(tknType token.TokenType, message string) {
	if p.current.Type == tknType {
		p.Advance()
		return
	}

	p.reportError(p.current, message)
}

func (p *Parser) parsePrec(prec PREC, assign bool) {
	p.Advance()

	prefix := p.getRule(p.previous.Type).prefix
	if prefix == nil {
		p.reportError(p.previous, "Expression "+token.ReversedTokenMap[p.previous.Type]+" not supported.")
		return
	}

	canAssign := prec <= PREC_ASSIGN
	prefix(p, canAssign && assign)

	for prec <= p.getRule(p.current.Type).prec {
		p.Advance()
		infix := p.getRule(p.previous.Type).infix
		infix(p, canAssign)
	}

	if canAssign && p.Match(token.EQUAL) {
		p.reportError(p.current, "Syntax Error, cannot assign on operator.")
	}
}

func (p *Parser) emitReturn() {
	p.emit(codes.INSTRUC_RETURN)
}

func (p *Parser) emitConst(v value.Value) {
	p.emit2(codes.INSTRUC_CONSTANT, p.makeConstant(v))
}

func (p *Parser) makeConstant(v value.Value) uint {
	return p.chk.AddVariable(v)
}

func (p *Parser) Expression(assign bool) {
	p.parsePrec(PREC_ASSIGN, assign)
}

func fcall(p *Parser, canAssign bool) {
}

func (p *Parser) Match(tokenType token.TokenType) bool {
	// DEBUG
	if !p.Check(tokenType) {
		return false
	}
	p.Advance()
	return true
}

func (p *Parser) Decl() {
	if p.Match(token.DECLARE) {
		p.declVar()
	} else {
		p.Statement()
	}
	// if p.ppanic {

	// }
}

func (p *Parser) declVar() {
	index := p.parseVar("Expected identifier Name.")
	if p.Match(token.EQUAL) {
		p.Expression(true)
	} else {
		p.emit(codes.INSTRUC_NIL)
	}
	p.Consume(token.SEMICOLON, "Malformed variable declaration.")
	p.defineDeclVar(index)
}

func (p *Parser) parseVar(msg string) (index uint) {
	p.Consume(token.IDENTIFIER, msg)
	return p.identifierConst(p.previous.Value)
}

func (p *Parser) identifierConst(name string) (index uint) {
	return p.makeConstant(value.NewString(name))
}

func (p *Parser) defineDeclVar(index uint) {
	p.emit2(codes.INSTRUC_DECL_GLOBAL, index)
}

func Unary(p *Parser, canAssign bool) {
	tknType := p.previous.Type

	p.parsePrec(PREC_UNARY, canAssign)

	switch tknType {
	case token.MINUS:
		p.emit(codes.INSTRUC_NEGATE)
	case token.EXCL:
		p.emit(codes.INSTRUC_NOT)
	default:
		return
	}
}

func Grouping(p *Parser, assign bool) {
	p.Expression(assign)
	p.Consume(token.CP, "Expected ')' after expression.")
}

func Variable(p *Parser, canAssign bool) {
	p.definedVar(p.previous.Value, canAssign)
}

func (p *Parser) definedVar(name string, canAssign bool) {
	index := p.identifierConst(name)

	if canAssign && p.Match(token.EQUAL) {
		/*
			Needs to be a declared variable before
			assigning.
		*/
		p.emit2(codes.INSTRUC_GET_DECL_GLOBAL, index)
		p.Expression(canAssign)
		/*
			DECL_GLOBAL -> initial declaration
			DECL_SET_GLOBAL -> assign on declared variable
		*/
		p.emit2(codes.INSTRUC_SET_DECL_GLOBAL, index)
	} else {
		p.emit2(codes.INSTRUC_GET_DECL_GLOBAL, index)
	}
}

func (p *Parser) Statement() {
	/*
		statement -> exprRessionStmt
					| printStmt
					| block

		block -> { delcare }
	*/
	if p.Match(token.PRINT) {
		p.PrintStmt()
	} else if p.Match(token.LB) {
		p.beginDeclScope()
		p.insideBlock()
		p.endDeclScope()
	} else {
		p.ExpressionStmt()
	}
}

func (p *Parser) beginDeclScope() {
	p.currentComp.scopeDepth++
}

func (p *Parser) insideBlock() {
	for !p.Check(token.RB) && !p.Check(token.EOF) {
		p.Decl()
	}
	p.Consume(token.RB, "No '}' at the end.")
}

func (p *Parser) endDeclScope() {
	p.currentComp.scopeDepth--
}

func (p *Parser) ExpressionStmt() {
	// variable has it
	p.Expression(true)
	p.Consume(token.SEMICOLON, "Malformed expression statement.")
	p.emit(codes.INSTRUC_POP)
}

func (p *Parser) PrintStmt() {
	p.Consume(token.OP, "Expected '(' after expression.")

	if !p.Match(token.CP) {
		Grouping(p, false)
	} else {
		// CP if only "print()"
		p.emit(codes.INSTRUC_NIL)
	}
	p.Consume(token.SEMICOLON, "Malformed print statement.")
	p.emit(codes.INSTRUC_PRINT)
}

func (p *Parser) Check(tokenType token.TokenType) bool {
	return p.current.Type == tokenType
}

func (p *Parser) Advance() {
	p.previous = p.current
	tkn, done := p.lex.Consume()

	if done {
		fmt.Println("done")
		return
	}
	if tkn.Type == token.ERR {
		p.reportError(p.current, "Error lexer.")
	}

	p.current = tkn
}

func Binary(p *Parser, canAssign bool) {
	tknType := p.previous.Type
	rule := p.getRule(tknType)
	p.parsePrec(rule.prec+1, canAssign)

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

func Number(p *Parser, canAssign bool) {
	dt := value.DetectNumberTypeByConversion(p.previous.Value)
	p.emitConst(value.New(p.previous.Value, dt))
}

func String(p *Parser, canAssign bool) {
	p.emitConst(value.NewString(p.previous.Value))
}

func Literal(p *Parser, canAssign bool) {
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

func (p *Parser) reportError(tkn *token.Token, what string) {
	if p.ppanic {
		return
	}

	fmt.Fprintf(os.Stderr, "[line:%d, pos:%d] Error %s, %s\n",
		tkn.Line, tkn.Position, token.ReversedTokenMap[tkn.Type], what)

	p.Perror = true
	p.ppanic = true

	/*
		for !p.Match(token.EOF) {
			p.Advance()
		}
	*/
}

func Init(lex *lexer.Lexer, chk *chunk.Chunk, comp *Compiler) *Parser {
	p := Parser{
		lex: lex,
		chk: chk,
	}
	p.tknMap = tknMap
	p.currentComp = comp
	/*
		Init MUST return a reference, otherwise
		the functions would not point to the correct
		address.

		HEAP HEAP HEAP...hurray?
	*/
	return &p
}
