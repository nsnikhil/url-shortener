package liberr

import "net/http"

var ValidationError = func(description string) ResponseError {
	return NewResponseError(validationError, http.StatusBadRequest, description)
}

var InternalError = func(description string) ResponseError {
	return NewResponseError(internalError, http.StatusInternalServerError, description)
}
