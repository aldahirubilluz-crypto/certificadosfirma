// ============================================
// internal/dto/response.go
// ============================================
package dto

type Response struct {
	Data    any  `json:"data"`
	Total   int  `json:"total,omitempty"`
	Success bool `json:"success"`
	Error   string `json:"error,omitempty"`
}

