package utils

const (
	NotFound            string = "not-found"
	InternalServerError        = "internal-server-error"
	BadRequest                 = "bad-request"
	Forbidden                  = "resource-forbidden"
	Unauthorized               = "unauthorized"
)

type Exception struct {
	StatusCode int
	Message    string
}

func (err Exception) Error() string {
	return err.Message
}

func NewException(statusCode int, message string) Exception {
	return Exception{StatusCode: statusCode, Message: message}
}
