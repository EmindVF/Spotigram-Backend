package customerrors

// Error meant to signal unauthorized behaviour.
type ErrUnauthorized struct {
	Message string
}

func (e *ErrUnauthorized) Error() string {
	return e.Message
}

// Error meant to signal invalid input.
type ErrInvalidInput struct {
	Message string
}

func (e *ErrInvalidInput) Error() string {
	return e.Message
}

// Error meant to signal missing content.
type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

// Error meant to signal internal failure.
type ErrInternal struct {
	Message string
}

func (e *ErrInternal) Error() string {
	return e.Message
}
