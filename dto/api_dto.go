package dto

import (
	"net/http"
)

type STATUS string

var (
	OK     STATUS = "ok"
	ERROR  STATUS = "error"
	FAILED STATUS = "failed"
)

type ApiResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func NewError(message string, err error) ApiResponse {

	return ApiResponse{
		Status:     string(ERROR),
		StatusCode: http.StatusExpectationFailed,
		Error:      err.Error(),
		Message:    message,
		Data:       nil,
	}
}

func NewSuccess(message string, data any) ApiResponse {
	return ApiResponse{
		Status:     string(OK),
		StatusCode: http.StatusOK,
		Error:      "",
		Message:    message,
		Data:       data,
	}
}

// Email verification DTOs
type SendVerificationCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type SendVerificationCodeResponse struct {
	Message string `json:"message"`
}

type VerifyEmailCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type VerifyEmailCodeResponse struct {
	Message       string `json:"message"`
	IsValid       bool   `json:"is_valid"`
	VerificationToken string `json:"verification_token,omitempty"`
}
