package model

// BaseAPIResponse represents the standard structure for API responses
type BaseAPIResponse struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ErrorResponse represents the structure for error responses, extending BaseAPIResponse
type ErrorResponse struct {
	BaseAPIResponse
	Error *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail contains specific error information
type ErrorDetail struct {
	Hint   int    `json:"hint,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// SuccessResponse represents the structure for successful responses with generic data
type SuccessResponse[T any] struct {
	BaseAPIResponse
	Data T `json:"data,omitempty"`
}

// HealthCheckResponse represents the structure for health check responses
type HealthCheckResponse struct {
	BaseAPIResponse
	Data HealthCheckData `json:"data"`
}

// HealthCheckData contains specific health check information
type HealthCheckData struct {
	Uptime    int64 `json:"uptime"`
	Timestamp int64 `json:"timestamp"`
}
