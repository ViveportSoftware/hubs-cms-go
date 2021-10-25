package errors

import "net/http"

// ErrorInfo is the response when error happened
type ErrorInfo struct {
	HttpStatus int       `json:"http_status"`
	ErrorBody  ErrorBody `json:"error"`
}

// ErrorBody is the body part of Error
type ErrorBody struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (e ErrorInfo) IsNil() bool {
	return ErrorInfo{} == e
}

// UnauthorizedError shows the error response when X-Account-Token not in header
var UnauthorizedError = ErrorInfo{
	HttpStatus: http.StatusUnauthorized,
	ErrorBody: ErrorBody{
		Code:    401,
		Status:  "Unauthorized",
		Message: "Unauthorized",
	},
}

// ForbiddenError shows the error response when X-Account-Token in header is not valid
var ForbiddenError = ErrorInfo{
	HttpStatus: http.StatusForbidden,
	ErrorBody: ErrorBody{
		Code:    403,
		Status:  "Forbidden",
		Message: "Forbidden",
	},
}

// InternalError shows the error response when error happened inside
var InternalError = ErrorInfo{
	HttpStatus: http.StatusInternalServerError,
	ErrorBody: ErrorBody{
		Code:    500,
		Status:  "Internal error",
		Message: "Internal error",
	},
}

func BadRequestError(code int, msg string) ErrorInfo {
	if len(msg) == 0 {
		msg = http.StatusText(http.StatusBadRequest)
	}

	return ErrorInfo{
		HttpStatus: http.StatusBadRequest,
		ErrorBody: ErrorBody{
			Code:    code,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: msg,
		},
	}

}
