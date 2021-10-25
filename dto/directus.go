package dto

import (
	"fmt"
	"net/http"
	"strings"
)

type DirectusM2MPatchRequest struct {
	Create []interface{} `json:"create,omitempty"`
	Update []interface{} `json:"update,omitempty"`
	Delete []int64       `json:"delete,omitempty"`
}

type DirectusErrorResponse struct {
	Status int             `json:"-"`
	Errors []DirectusError `json:"errors"`
}

type DirectusError struct {
	Message    string                 `json:"message"`
	Extensions DirectusErrorExtension `json:"extensions"`
}

type DirectusErrorExtension struct {
	Code string `json:"code"`
}

func (e *DirectusErrorResponse) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return ""
	}

	errStr := make([]string, len(e.Errors))
	for i := range e.Errors {
		errStr[i] = fmt.Sprintf("%v: %v", e.Errors[1].Extensions.Code, e.Errors[0].Message)
	}

	return fmt.Sprintf("%v", strings.Join(errStr[:], "\n"))
}

func DirectusErrorResponseFromHttpStstus(httpStatus int) *DirectusErrorResponse {
	err := DirectusError{
		Message: http.StatusText(httpStatus),
		Extensions: DirectusErrorExtension{
			Code: http.StatusText(httpStatus),
		},
	}
	return &DirectusErrorResponse{
		Status: httpStatus,
		Errors: []DirectusError{err},
	}
}

type DirectusL10N interface {
	UpdateTranslation()
}

func UpdateTranslationIfNeed(source, target *string) {
	if len(*target) > 0 {
		*source = *target
	}
}
