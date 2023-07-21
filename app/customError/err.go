package customError

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ce CustomError) Error() string {
	return ce.Message
}

var Others = CustomError{
	Code:    0,
	Message: "some error found",
}
var UserNotFound = CustomError{
	Code:    1,
	Message: "user not found",
}
var UserAlreadyExists = CustomError{
	Code:    2,
	Message: "user already exists",
}
var AddressAlreadyConfirmed = CustomError{
	Code:    3,
	Message: "email address is already confirmed",
}

var OperationNotAllowed = CustomError{
	Code:    100,
	Message: "operation not allowed",
}
var MethodNotAllowed = CustomError{
	Code:    101,
	Message: "method not allowed",
}

var InvalidToken = CustomError{
	Code:    102,
	Message: "token is invalid",
}
