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
