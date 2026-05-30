package http

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Erros map[string]interface{} `json:"errors"`
}

type ErrorsResponse struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	StatusCode int       `json:"status_code"`
	Data       ErrorData `json:"data"`
}

type ErrorData struct {
	Error string `json:"error"`
}

// pkg/http/error.go
func NewResponseError(message string, statusCode int, err error) ErrorsResponse {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return ErrorsResponse{
		Success:    false,
		Message:    message,
		StatusCode: statusCode,
		Data: ErrorData{
			Error: errMsg,
		},
	}
}

func NewError(err error) Error {
	e := Error{}
	e.Erros = make(map[string]interface{})
	switch v := err.(type) {
	default:
		e.Erros["body"] = v.Error()
	}
	return e
}

func NewValidatorError(err error) Error {
	e := Error{}
	e.Erros = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		e.Erros[v.Field()] = fmt.Sprintf("%v", v.Tag())
	}
	return e
}

func AccessForbidden() Error {
	e := Error{}
	e.Erros = make(map[string]interface{})
	e.Erros["body"] = "access forbidden"
	return e
}

func NotFound() Error {
	e := Error{}
	e.Erros = make(map[string]interface{})
	e.Erros["body"] = "resource not found"
	return e
}
