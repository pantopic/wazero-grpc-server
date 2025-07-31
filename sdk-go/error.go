package grpc_server

const (
	errCodeEmpty uint32 = iota
	errCodeUnknown
	errCodeInvalid
	errCodeUnrecognized
	errCodeNotImplemented
	errCodeMalformed
	errCodeUnexpected
	errCodeMarshal
)

var (
	ErrUnknown        = Error{errCodeUnknown, "Unknown"}
	ErrInvalid        = Error{errCodeInvalid, "Invalid"}
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
