package utils

type HttpError struct {
	Status  bool   `json:"status" default:"false"`
	Code    int    `json:"code" example:"500"`
	Message string `json:"message"`
}

type HttpSuccess struct {
	Status  bool        `json:"status" default:"true"`
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message"`
	Data    interface{} `json:"data" swaggertype:"object"`
}

// Warp the error info in a object
func NewSuccess(data interface{}, message string) *HttpSuccess {
	formattedMessage := GetMessage(message)

	return &HttpSuccess{
		Status:  true,
		Code:    200,
		Message: formattedMessage,
		Data:    data,
	}
}

// Warp the error info in a object
func NewError(code int, message string) *HttpError {
	formattedMessage := GetMessage(message)

	return &HttpError{
		Status:  false,
		Code:    code,
		Message: formattedMessage,
	}
}

// Error makes it compatible with `error` interface.
func (e *HttpError) Error() string {
	return e.Message
}
