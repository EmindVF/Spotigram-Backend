package customerrors

type ErrInvalidInput struct {
	Message string
}

func (e *ErrInvalidInput) Error() string {
	return e.Message
}

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

type ErrInternal struct {
	Message string
}

func (e *ErrInternal) Error() string {
	return e.Message
}
