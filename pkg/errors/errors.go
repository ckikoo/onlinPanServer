package errors

import "github.com/pkg/errors"

// Define alias
var (
	New          = errors.New
	Wrap         = errors.Wrap
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack
	WithMessage  = errors.WithMessage
	WithMessagef = errors.WithMessagef
)

var (
	ErrInvalidToken    = "invalid signature"
	ErrNoPerm          = "no permission"
	ErrNotFound        = "not found"
	ErrMethodNotAllow  = "method not allowed"
	ErrConflict        = "request conflict"
	ErrTooManyRequests = "too many requests"
	ErrInternalServer  = "internal server error"
	ErrBadRequest      = "bad request"
	ErrUserDisable     = "user forbidden"
)
