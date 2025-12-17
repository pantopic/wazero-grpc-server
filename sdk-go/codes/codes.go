package codes

import (
	"strconv"
)

type Code uint32

const (
	OK                 Code = 0
	Canceled           Code = 1
	Unknown            Code = 2
	InvalidArgument    Code = 3
	DeadlineExceeded   Code = 4
	NotFound           Code = 5
	AlreadyExists      Code = 6
	PermissionDenied   Code = 7
	ResourceExhausted  Code = 8
	FailedPrecondition Code = 9
	Aborted            Code = 10
	OutOfRange         Code = 11
	Unimplemented      Code = 12
	Internal           Code = 13
	Unavailable        Code = 14
	DataLoss           Code = 15
	Unauthenticated    Code = 16
)

func (c Code) String() (name string) {
	name, ok := codeNames[c]
	if !ok {
		name = strconv.Itoa(int(c))
	}
	return
}

var codeNames = map[Code]string{
	0:  "OK",
	1:  "Canceled",
	2:  "Unknown",
	3:  "InvalidArgument",
	4:  "DeadlineExceeded",
	5:  "NotFound",
	6:  "AlreadyExists",
	7:  "PermissionDenied",
	8:  "ResourceExhausted",
	9:  "FailedPrecondition",
	10: "Aborted",
	11: "OutOfRange",
	12: "Unimplemented",
	13: "Internal",
	14: "Unavailable",
	15: "DataLoss",
	16: "Unauthenticated",
}
