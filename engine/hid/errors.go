package hid

type BindingError string

func (e BindingError) Error() string {
	return string(e)
}

const (
	// ErrNilCallback binding callback must not be nil.
	ErrNilCallback = BindingError("binding callback is nil")

	// ErrBindingAlreadyExists binding with same trigger already registered.
	ErrBindingAlreadyExists = BindingError("binding already exists")
)
