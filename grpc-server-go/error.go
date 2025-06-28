package grpc

const (
	errCodeEmpty uint32 = iota
	errCodeUnknown
	errCodeUnrecognized
	errCodeNotImplemented
	errCodeMalformed
	errCodeUnexpected
	errCodeMarshal
)

var (
	ErrUnrecognized   = Error{errCodeUnrecognized, "Unrecognized"}
	ErrNotImplemented = Error{errCodeNotImplemented, "Not Implemented"}
	ErrMalformed      = Error{errCodeMalformed, "Malformed"}
	ErrUnexpected     = Error{errCodeUnexpected, "Unexpected"}
	ErrMarshal        = Error{errCodeMarshal, "Marshal"}
)

type Error struct {
	code uint32
	msg  string
}

func (e Error) Error() string {
	return e.msg
}
