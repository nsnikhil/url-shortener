package liberr

type ResponseError interface {
	ErrorCode() string
	StatusCode() int
	error
}

type responseError struct {
	errCode     string
	statusCode  int
	description string
}

func (re responseError) ErrorCode() string {
	return re.errCode
}

func (re responseError) StatusCode() int {
	return re.statusCode
}

func (re responseError) Error() string {
	return re.description
}

func NewResponseError(errCode string, statusCode int, description string) ResponseError {
	return &responseError{
		errCode:     errCode,
		statusCode:  statusCode,
		description: description,
	}
}
