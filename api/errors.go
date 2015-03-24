package api

import (
	"fmt"
)

const (
	ERR_UNKNOWN_ACTION = 100
	ERR_INVALID_ARGS   = 101
	ERR_DATA_NOT_FOUND = 102

	ERR_INVALID_TOKEN      = 201
	ERR_INVALID_CREDENTIAL = 202

	ERR_INTERAL_ERROR = 301
)

var ErrMessageMap = map[int]string{
	ERR_UNKNOWN_ACTION: "unknown action: %s",
	ERR_INVALID_ARGS:   "invalid args: %s",

	ERR_DATA_NOT_FOUND:     "data not found: %s",
	ERR_INVALID_TOKEN:      "invalid token: %s",
	ERR_INVALID_CREDENTIAL: "invalid credential",

	ERR_INTERAL_ERROR: "internal error: %s",
}

type ApiErr struct {
	Code    int
	Message string
}

func (e *ApiErr) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

func NewApiErr(code int, a ...interface{}) *ApiErr {
	desc := fmt.Sprintf(ErrMessageMap[code], a...)
	return &ApiErr{
		Code:    code,
		Message: desc,
	}
}
