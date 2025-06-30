package wazero_grpc_server

const (
	errCodeEmpty uint32 = iota
	errCodeDone
	errCodeUnknown
	errCodeInvalid
	errCodeUnrecognized
	errCodeNotImplemented
	errCodeMalformed
	errCodeUnexpected
	errCodeMarshal
)

var (
	ErrUnknown        = &Error{errCodeUnknown, "Unknown"}
	ErrInvalid        = &Error{errCodeInvalid, "Invalid"}
	ErrUnrecognized   = &Error{errCodeUnrecognized, "Unrecognized"}
	ErrNotImplemented = &Error{errCodeNotImplemented, "Not Implemented"}
	ErrMalformed      = &Error{errCodeMalformed, "Malformed"}
	ErrUnexpected     = &Error{errCodeUnexpected, "Unexpected"}
	ErrMarshal        = &Error{errCodeMarshal, "Marshal"}
)

var errorsByCode = map[uint32]error{
	errCodeUnknown:        ErrUnknown,
	errCodeInvalid:        ErrInvalid,
	errCodeUnrecognized:   ErrUnrecognized,
	errCodeNotImplemented: ErrNotImplemented,
	errCodeMalformed:      ErrMalformed,
	errCodeUnexpected:     ErrUnexpected,
	errCodeMarshal:        ErrMarshal,
}

type Error struct {
	code uint32
	msg  string
}

func (e Error) Error() string {
	return e.msg
}
