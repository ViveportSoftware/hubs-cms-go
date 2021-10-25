package errors

import "net/http"

// AccountsInvalidRequestFormat shows the error response when request payload is invalid
var AccountsInvalidRequestFormat = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400200,
		Status:  "Bad Request",
		Message: "Invalid request format",
	},
}

// AccountsInvalidAccountID shows the error response when request payload is invalid
var AccountsInvalidAccountID = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400201,
		Status:  "Bad Request",
		Message: "Invalid URI: account_id",
	},
}

// AccountsInvalidActiveAvatarID shows the error response when request payload is invalid
var AccountsInvalidActiveAvatarID = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400202,
		Status:  "Bad Request",
		Message: "Invalid payload: active_avatar",
	},
}

// AccountsInvalidDisplayName shows the error response when request payload is invalid
var AccountsInvalidDisplayName = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400203,
		Status:  "Bad Request",
		Message: "Invalid payload: display_name",
	},
}
