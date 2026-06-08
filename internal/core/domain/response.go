package domain

// Envelope defines the standardized outer wrapper
type Envelope[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Error   string `json:"error"`
}
