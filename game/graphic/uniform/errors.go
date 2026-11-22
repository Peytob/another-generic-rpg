package uniform

import "errors"

var (
	ErrVariableNotFound       = errors.New("variable not found")
	ErrBlockNotFound          = errors.New("uniform block not found in repository")
	ErrBlockAlreadyRegistered = errors.New("uniform block already registered")
	ErrBlockNotRegistered     = errors.New("uniform block not found in registered uniform blocks")
	ErrInvalidBlock           = errors.New("uniform block already deleted or id is not buffer")
)
