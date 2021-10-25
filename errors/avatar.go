package errors

import "net/http"

// AvatarsInvalidRequestFormat shows the error response when request payload is invalid
var AvatarsInvalidRequestFormat = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400300,
		Status:  "Bad Request",
		Message: "Invalid request format",
	},
}

// AvatarsInvalidStart shows the error response when request payload is invalid
var AvatarsInvalidStart = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400301,
		Status:  "Bad Request",
		Message: "Invalid param: start",
	},
}

// AvatarsInvalidLimit shows the error response when request payload is invalid
var AvatarsInvalidLimit = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400302,
		Status:  "Bad Request",
		Message: "Invalid param: limit",
	},
}

// AvatarsInvalidID shows the error response when request payload is invalid
var AvatarsInvalidID = ErrorInfo{
	HttpStatus: http.StatusBadRequest,
	ErrorBody: ErrorBody{
		Code:    400303,
		Status:  "Bad Request",
		Message: "Invalid uri: avatar_id",
	},
}
