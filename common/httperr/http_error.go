package httperr

import "errors"

var (
	// 4xx
	ErrBadRequest   = errors.New("bad request")  // 400
	ErrUnauthorized = errors.New("unauthorized") // 401
	ErrForbidden    = errors.New("forbidden")    // 403
	ErrNotFound     = errors.New("not found")    // 404

	// 5xx
	ErrInternalServerError = errors.New("internal server error") // 500
)

type httpError struct {
	StatusCode int
	Message    string

	error
}

func Unauthorized(err error, m ...string) error {
	msg := "Unauthorized"
	if len(m) > 0 {
		msg = m[0]
	}
	return &httpError{
		StatusCode: 401,
		Message:    msg,
		error:      err,
	}
}
