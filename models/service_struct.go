package models

// Service Data Struct
type Service struct {
	Id          string  `json:"id"`
	ServiceId   *int64  `json:"serviceId,omitempty"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Version     float64 `json:"version" binding:"required"`
}
