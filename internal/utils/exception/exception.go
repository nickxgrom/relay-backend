package exception

type enumType struct {
	NotFound            string
	InternalServerError string
	BadRequest          string
	Forbidden           string
	Unauthorized        string

	OrganizationNotFound     string
	UserNotFound             string
	EmployeeAlreadyExist     string
	InvalidEmail             string
	EmailTaken               string
	TokenNotFound            string
	TokenExpired             string
	IncorrectEmailOrPassword string
}

var Enum = enumType{
	NotFound:                 "not-found",
	InternalServerError:      "internal-server-error",
	BadRequest:               "bad-request",
	Forbidden:                "resource-forbidden",
	Unauthorized:             "unauthorized",
	OrganizationNotFound:     "organization-not-found",
	UserNotFound:             "user-not-found",
	EmployeeAlreadyExist:     "employee-already-exist",
	InvalidEmail:             "email-not-valid",
	EmailTaken:               "email-already-taken",
	TokenNotFound:            "token-not-found",
	TokenExpired:             "token-expired",
	IncorrectEmailOrPassword: "incorrect-email-or-password",
}

type Exception struct {
	StatusCode int
	Message    string
	Details    map[string]interface{}
}

func (err Exception) Error() string {
	return err.Message
}

func NewException(statusCode int, message string) Exception {
	return Exception{StatusCode: statusCode, Message: message}
}

func NewDetailsException(statusCode int, message string, details map[string]interface{}) Exception {
	return Exception{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

//TODO: make not found methods with generics
