package httpmsgerrorhandler

import (
	"errors"
	"net/http"
	"todoapp/internal/pkg/errmsg"
	"todoapp/internal/pkg/richerror"
)

func Error(err error) (message string, code int) {
	// Use errors.As to unwrap properly
	var re *richerror.RichError
	if errors.As(err, &re) {
		msg := re.GetMessage()
		code := MapKindToHTTPStatusCode(re.GetKind())
		
		// Don't expose internal error details
		if code >= 500 {
			msg = errmsg.ErrorMsgSomethingWentWrong
		}
		return msg, code
	}
	
	// Unknown error type
	return errmsg.ErrorMsgSomethingWentWrong, http.StatusInternalServerError
}

func MapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}