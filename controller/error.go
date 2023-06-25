package controller

import "fmt"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (a APIError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", a.Code, a.Message)
}
