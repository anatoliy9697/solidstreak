package errors

type Error struct {
	HTTPCode int    `json:"status,string"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
}

func (e Error) Error() string {
	msg := e.Title
	if e.Detail != "" {
		msg += ": " + e.Detail
	}
	return msg
}

func New(httpCode int, title string, detail string) Error {
	return Error{
		HTTPCode: httpCode,
		Title:    title,
		Detail:   detail,
	}
}

func ErrNotFound(detail string) Error {
	return New(404, "not found", detail)
}

func ErrInternal(detail string) Error {
	return New(500, "internal server error", detail)
}

func ErrBadRequest(detail string) Error {
	return New(400, "bad request", detail)
}

func ErrUnauthorized(detail string) Error {
	return New(401, "unauthorized", detail)
}

func ErrForbidden(detail string) Error {
	return New(403, "forbidden", detail)
}
