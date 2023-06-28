package models

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

func BuildResponse(status string, message string, data interface{}, err string) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   err,
	}
}