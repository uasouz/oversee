package core

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrorCode int

const (
	ErrorCodeAlreadyPersistedLog = iota + 1001
	ErrorCodeInvalidErrorFormat
)

type Error struct {
	Message string
	Code    ErrorCode
}

func (e *Error) Error() string {
	return fmt.Sprintf("[ERROR] [%d] %s", e.Code, strings.ToLower(e.Message))
}

func ErrorWithMessage(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ParseString(s string) (*Error, error) {
	if s[0] != '[' {
		return nil, fmt.Errorf("not a valid error string")
	}

	parts := strings.SplitN(s, " ", 3)

	if len(parts) != 2 {
		return nil, ErrorWithMessage(ErrorCodeInvalidErrorFormat, "Invalid error format")
	}

	code := strings.Trim(parts[1], "[]")

	codeValue, err := strconv.Atoi(code)

	if err != nil {
		return nil, fmt.Errorf("failed to parse error code")
	}

	message := strings.TrimSpace(parts[2])

	return &Error{
		Code:    ErrorCode(codeValue),
		Message: message,
	}, nil

}

var ErrorAlreadyPersistedLog = ErrorWithMessage(ErrorCodeAlreadyPersistedLog, "Already Persisted")
