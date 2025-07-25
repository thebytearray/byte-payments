package dto

type ApiResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}
