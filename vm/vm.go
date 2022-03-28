package vm

import (
	"fmt"
	"os"

	"github.com/badc0re/hprog/chunk"
	"github.com/badc0re/hprog/codes"
	"github.com/badc0re/hprog/lexer"
	"github.com/badc0re/hprog/stack"
	"github.com/badc0re/hprog/token"
	"github.com/badc0re/hprog/value"
)

type VM struct {
	chunk   *chunk.Chunk
	ip      *interface{}
	counter int
	stack   stack.Stack
}

type INTER_RESULT int

const (
	INTER_ILLEGAL INTER_RESULT = iota

	INTER_OK
	INTER_COMPILE_ERROR
	INTER_RUNTIME_ERROR
)

type ParseFn func()

//type OpFn func()

type ParseRule struct {
	prefix ParseFn
	infix  ParseFn
	prec   PREC
}

/*
type OpRule struct {
	opfn OpFn
}

type OpKey struct {
	a value.VALUE_TYPE
	b value.VALUE_TYPE
}
*/

func (vm *VM) InitVM() {

}

func (vm *VM) FreeVM() {

}

func (vm *VM) Move() interface{} {
	vm.ip = &vm.chunk.Code[vm.counter]
	vm.counter++
	return *vm.ip
}

func (vm *VM) ReadConstant() value.Value {
	vm.Move()
	index := (*vm.ip).(uint)
	return vm.chunk.Constants.Values[index]
}

func (vm *VM) binaryOP(op string) INTER_RESULT {
	b := vm.stack.Pop()
	a := vm.stack.Pop()

	// allow "+" for strings
	if !(value.IsNumberType(a.VT) &&
		value.IsNumberType(b.VT)) {
		// works!
		fmt.Println("bla")
		return INTER_RUNTIME_ERROR
	}

	switch op {
	case "+":
		vm.stack.Push(value.Add(&a, &b))
	case "-":
		vm.stack.Push(value.Sub(&a, &b))
	case "/":
		vm.stack.Push(value.Divide(&a, &b))
	case "*":
		vm.stack.Push(value.Multiply(&a, &b))
	}
	return INTER_OK
}

func (v *VM) run() INTER_RESULT {
	for {
		instruct := v.Move()

		switch instruct {
		case codes.INSTRUC_CONSTANT:
			constant := v.ReadConstant()
			value.PrintValue(constant)
			v.stack.Push(constant)
			break
		case codes.INSTRUC_NIL:
		case codes.INSTRUC_TRUE:
			// TODO: change
			v.stack.Push(value.NewBool(true, value.VT_BOOL))
		case codes.INSTRUC_FALSE:
			// TODO: change
			v.stack.Push(value.NewBool(false, value.VT_BOOL))
		case codes.INSTRUC_ADDITION:
			v.binaryOP("+")
		case codes.INSTRUC_SUBSTRACT:
			v.binaryOP("-")
		case codes.INSTRUC_MULTIPLY:
			v.binaryOP("*")
		case codes.INSTRUC_DIVIDE:
			v.binaryOP("/")
		case codes.INSTRUC_NEGATE:
			_v := v.stack.Pop()
			if !(value.IsNumberType(_v.VT) ||
				value.IsBooleanType(_v.VT)) {
				// error
				return INTER_RUNTIME_ERROR
			}
			v.stack.Push(value.Negate(_v))
		case codes.INSTRUC_RETURN:
			fmt.Println("RETURN")
			fmt.Printf("STACK POP:, %#v", v.stack.Pop())
			return INTER_OK
		}
	}
}

func (p *Parser) getRule(tknType token.TokenType) ParseRule {
	return p.tknMap[tknType]
}

func Compile(source string, chk *chunk.Chunk) INTER_RESULT {
	lex := lexer.Init(source)
	p := Parser{
		lex: &lex,
		chk: chk,
	}
	/*
			opMap := map[OpKey]OpRule{
				OpKey{a: value.VT_FLOAT, b: value.VT_FLOAT}: {nil},
				OpKey{a: value.VT_INT, b: value.VT_FLOAT}:   {nil},
				OpKey{a: value.VT_FLOAT, b: value.VT_INT}:   {nil},
			}
		fmt.Println(opMap)
	*/
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
		token.EXCL:          {nil, nil, PREC_NONE},
		token.EXCL_EQUAL:    {nil, nil, PREC_NONE},
		token.EQUAL:         {nil, nil, PREC_NONE},
		token.EQUAL_EQUAL:   {nil, nil, PREC_NONE},
		token.GREATER:       {nil, nil, PREC_NONE},
		token.GREATER_EQUAL: {nil, nil, PREC_NONE},
		token.LESS:          {nil, nil, PREC_NONE},
		token.LESS_EQUAL:    {nil, nil, PREC_NONE},
		token.IDENTIFIER:    {nil, nil, PREC_NONE},
		token.STRING:        {nil, nil, PREC_NONE},
		token.NUMBER:        {p.Number, nil, PREC_NONE},
		token.AND:           {nil, nil, PREC_NONE},
		// maybe not
		// token.CLASS:      {nil, nil, PREC_NONE},
		token.ELSE:       {nil, nil, PREC_NONE},
		token.BOOL_FALSE: {p.Literal, nil, PREC_NONE},
		token.FOR:        {nil, nil, PREC_NONE},
		token.FUNCTION:   {nil, nil, PREC_NONE},
		token.IF:         {nil, nil, PREC_NONE},
		token.NIL:        {p.Literal, nil, PREC_NONE},
		token.OR:         {nil, nil, PREC_NONE},
		token.PRINT:      {nil, nil, PREC_NONE},
		token.RETURN:     {nil, nil, PREC_NONE},
		//token.SUPER:     {nil, nil, PREC_NONE},
		//token.THIS:     {nil, nil, PREC_NONE},
		token.BOOL_TRUE: {p.Literal, nil, PREC_NONE},
		token.VAR:       {nil, nil, PREC_NONE},
		token.WHILE:     {nil, nil, PREC_NONE},
		token.ERR:       {nil, nil, PREC_NONE},
		token.EOF:       {nil, nil, PREC_NONE},
	}
	p.tknMap = tknMap

	p.Advance()
	p.Expression()
	p.Consume(token.EOF, "Expected end.")
	p.endCompile()

	if p.perror {
		return INTER_COMPILE_ERROR
	}
	return INTER_OK
}

type Parser struct {
	current  *token.Token
	previous *token.Token
	lex      *lexer.Lexer
	perror   bool
	ppanic   bool
	tknMap   map[token.TokenType]ParseRule

	// todo
	chk *chunk.Chunk
}

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

func (p *Parser) Consume(tknType token.TokenType, message string) {
	if p.current.Type == tknType {
		p.Advance()
		return
	}

	fmt.Println(p.current.Type, tknType)
	// TODO: ERROR
	p.reportError(p.current, "Error")
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
		break
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
		break
	case token.MINUS:
		p.emit(codes.INSTRUC_SUBSTRACT)
		break
	case token.STAR:
		p.emit(codes.INSTRUC_MULTIPLY)
		break
	case token.SLASH:
		p.emit(codes.INSTRUC_DIVIDE)
		break
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

func (p *Parser) endCompile() {
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

func (p *Parser) Advance() {
	p.previous = p.current

	tkn := <-p.lex.Consume()
	p.current = &tkn
	if tkn.Type == token.ERR {
		p.reportError(p.current, "Error")
	}
}

func (p *Parser) reportError(tkn *token.Token, what string) {
	p.ppanic = true
	fmt.Fprintf(os.Stderr, "[line:%d, pos:%d] Error, %s\n",
		tkn.Line, tkn.Position, what)
	p.perror = true
}

func (v *VM) Interpret(source string) INTER_RESULT {
	chk := chunk.Chunk{}

	if Compile(source, &chk) == INTER_COMPILE_ERROR {
		// parser.ppanic = true
		// parser.perror = true
		return INTER_COMPILE_ERROR
	}

	chunk.DissasChunk(&chk, "test")
	if len(chk.Code) != 0 {
		v.chunk = &chk
		v.counter = 0
		v.ip = &v.chunk.Code[v.counter]
		return v.run()
	}
	return INTER_OK
}
