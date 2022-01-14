package errors

type SyntaxError struct {
	s string
}

func (e *SyntaxError) Error() string {
	return e.s
}

func NewSyntaxError(text string) error {
	return &SyntaxError{text}
}

type CompileError struct {
	s string
}

func (e *CompileError) Error() string {
	return e.s
}

func NewCompileError(text string) error {
	return &CompileError{text}
}
