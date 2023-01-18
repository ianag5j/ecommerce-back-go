package dto

import "fmt"

const (
	INTERNAL_ERROR           = 0000
	STORE_NAME_USED          = 1001
	USER_ALREADY_HAS_STORE   = 1002
	USER_NOT_HAS_CREDENTIALS = 2001
	UALA_ERROR               = 2002
)

type RequestError struct {
	StatusCode int    `json:"-"`
	ErrorCode  int    `json:"error_code"`
	Err        string `json:"error"`
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("%v", r.Err)
}
