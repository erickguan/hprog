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

var MAX_STACK_SIZE = 256

type INTER_RESULT int

const (
	INTER_ILLEGAL INTER_RESULT = iota

	INTER_OK
	INTER_COMPILE_ERROR
	INTER_RUNTIME_ERROR
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

type VM struct {
	chunk        *chunk.Chunk
	ip           *interface{}
	counter      int
	vstack       stack.Stack
	valueTypeMap map[OpKey]value.VALUE_TYPE
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

type ParseFn func()

type ParseRule struct {
	prefix ParseFn
	infix  ParseFn
	prec   PREC
}

type OpKey struct {
	a value.VALUE_TYPE
	b value.VALUE_TYPE
}

func (vm *VM) InitVM() {
	vm.vstack = stack.Stack{
		Sarray: make([]value.Value, MAX_STACK_SIZE),
		Top:    -1,
	}
	valueTypeMap := map[OpKey]value.VALUE_TYPE{
		OpKey{a: value.VT_FLOAT, b: value.VT_FLOAT}: value.VT_FLOAT,
		OpKey{a: value.VT_INT, b: value.VT_FLOAT}:   value.VT_FLOAT,
		OpKey{a: value.VT_FLOAT, b: value.VT_INT}:   value.VT_FLOAT,
		OpKey{a: value.VT_INT, b: value.VT_INT}:     value.VT_INT,
	}
	vm.valueTypeMap = valueTypeMap
}

func (vm *VM) ResetStack() {
	vm.vstack = stack.Stack{
		Sarray: make([]value.Value, MAX_STACK_SIZE),
	}
}

func (vm *VM) FreeVM() {
	vm.vstack = stack.Stack{}
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
	b := vm.vstack.Pop()
	a := vm.vstack.Pop()

	if !value.IsSameType(a.VT, b.VT) {
		vt := vm.valueTypeMap[OpKey{a: a.VT, b: b.VT}]
		a, b = value.ConvertToExpectedType2(a, b, vt)
	}
	// allow "+" for strings
	if !(value.IsNumberType(a.VT) &&
		value.IsNumberType(b.VT)) {
		// works!
		return INTER_RUNTIME_ERROR
	}

	switch op {
	case "+":
		vm.vstack.Push(value.Add(&a, &b))
	case "-":
		vm.vstack.Push(value.Sub(&a, &b))
	case "/":
		vm.vstack.Push(value.Divide(&a, &b))
	case "*":
		vm.vstack.Push(value.Multiply(&a, &b))
	case ">":
		vm.vstack.Push(value.Greater(&a, &b))
	case "<":
		vm.vstack.Push(value.Less(&a, &b))
	}
	return INTER_OK
}

func (v *VM) StackTrace() {
	fmt.Println("== Stack Trace ==")
	fmt.Println("[")
	for i := 0; i < v.vstack.Top+1; i++ {
		value.PrintValue(i, v.vstack.Sarray[i])
	}
	fmt.Println("]")
	fmt.Printf("== End Stack Trace ==\n\n")
}

func (vm *VM) run() INTER_RESULT {
	for {
		instruct := vm.Move()
		switch instruct {
		case codes.INSTRUC_CONSTANT:
			constant := vm.ReadConstant()
			vm.vstack.Push(constant)
			break
		case codes.INSTRUC_NIL:
			vm.vstack.Push(value.New("", value.VT_NIL))
		case codes.INSTRUC_TRUE:
			vm.vstack.Push(value.NewBool(true))
		case codes.INSTRUC_FALSE:
			vm.vstack.Push(value.NewBool(false))
		case codes.INSTRUC_NOT:
			_v, err := vm.vstack.Peek(0)
			if !value.IsBooleanType(_v.VT) || err != nil {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.vstack.Push(value.Negate(vm.vstack.Pop()))
		case codes.INSTRUC_NEGATE:
			a, err := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) || err != nil {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.vstack.Push(value.Negate(vm.vstack.Pop()))
		case codes.INSTRUC_EQUAL:
			b := vm.vstack.Pop()
			a := vm.vstack.Pop()
			if !value.IsSameType(a.VT, b.VT) {
				vt := vm.valueTypeMap[OpKey{a: a.VT, b: b.VT}]
				a, b = value.ConvertToExpectedType2(a, b, vt)
			}
			vm.vstack.Push(value.Equal(&a, &b))
		case codes.INSTRUC_ADDITION:
			vm.binaryOP("+")
		case codes.INSTRUC_SUBSTRACT:
			vm.binaryOP("-")
		case codes.INSTRUC_MULTIPLY:
			vm.binaryOP("*")
		case codes.INSTRUC_DIVIDE:
			vm.binaryOP("/")
		case codes.INSTRUC_GREATER:
			a, _ := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.binaryOP(">")
		case codes.INSTRUC_LESS:
			a, _ := vm.vstack.Peek(0)
			if !value.IsNumberType(a.VT) {
				// error
				return INTER_RUNTIME_ERROR
			}
			vm.binaryOP("<")
		case codes.INSTRUC_RETURN:
			fmt.Printf("RETURN; STACK POP:, %#v", vm.vstack.Pop())
			return INTER_OK
		}
		vm.StackTrace()
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

func (vm *VM) Interpret(source string) INTER_RESULT {
	chk := chunk.Chunk{}

	if Compile(source, &chk) == INTER_COMPILE_ERROR {
		// parser.ppanic = true
		// parser.perror = true
		return INTER_COMPILE_ERROR
	}

	/* DEBUG */
	chunk.DissasChunk(&chk, "test")

	if len(chk.Code) != 0 {
		/* INIT START */
		vm.chunk = &chk
		vm.counter = 0
		vm.ip = &vm.chunk.Code[vm.counter]
		/* INIT END */
		return vm.run()
	}
	return INTER_OK
}
