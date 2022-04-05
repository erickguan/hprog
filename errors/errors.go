package errors

type SyntaxError struct {
	s string
}

type CompileError struct {
	s string
}

func (e *SyntaxError) Error() string {
	return e.s
}

func (e *CompileError) Error() string {
	return e.s
}

func NewCompileError(text string) error {
	return &CompileError{text}
}

func NewSyntaxError(text string) error {
	return &SyntaxError{text}
}
