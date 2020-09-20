package contract

type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   Error       `json:"error,omitempty"`
	Success bool        `json:"success"`
}

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Data:    data,
		Success: true,
	}
}

func NewFailureResponse(errorCode, description string) APIResponse {
	return APIResponse{
		Error: Error{
			Code:    errorCode,
			Message: description,
		},
		Success: false,
	}
}
