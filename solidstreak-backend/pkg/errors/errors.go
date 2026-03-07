package errors

type Error struct {
	Err      error  `json:"-"`
	HTTPCode int    `json:"status,string"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
}

func (e *Error) Error() string {
	msg := e.Title
	if e.Detail != "" {
		msg += ": " + e.Detail
	}
	return msg
}

func (e *Error) Unwrap() error {
	return e.Err
}

var ErrNotFound = &Error{Err: nil, HTTPCode: 404, Title: "not found", Detail: ""}
var ErrInternal = &Error{Err: nil, HTTPCode: 500, Title: "internal server error", Detail: ""}
var ErrBadRequest = &Error{Err: nil, HTTPCode: 400, Title: "bad request", Detail: ""}
var ErrUnauthorized = &Error{Err: nil, HTTPCode: 401, Title: "unauthorized", Detail: ""}
var ErrForbidden = &Error{Err: nil, HTTPCode: 403, Title: "forbidden", Detail: ""}

func NewNotFoundErr(detail string) *Error {
	return &Error{Err: ErrNotFound, HTTPCode: 404, Title: "not found", Detail: detail}
}

func NewInternalErr(detail string) *Error {
	return &Error{Err: ErrInternal, HTTPCode: 500, Title: "internal server error", Detail: detail}
}

func NewBadRequestErr(detail string) *Error {
	return &Error{Err: ErrBadRequest, HTTPCode: 400, Title: "bad request", Detail: detail}
}

func NewUnauthorizedErr(detail string) *Error {
	return &Error{Err: ErrUnauthorized, HTTPCode: 401, Title: "unauthorized", Detail: detail}
}

func NewForbiddenErr(detail string) *Error {
	return &Error{Err: ErrForbidden, HTTPCode: 403, Title: "forbidden", Detail: detail}
}
