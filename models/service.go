package models

type ServiceRecord struct {
	Id          string  `json:"id"`
	ServiceId   *int64  `json:"serviceId,omitempty"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Version     float64 `json:"version" binding:"required"`
}

type Service struct {
	Id          string  `json:"id"`
	ServiceId   int64   `json:"serviceId" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Version     float64 `json:"version" binding:"required"`
	Versions    int64   `json:"versions" binding:"required"`
}

type ServicesSearchRequest struct {
	NameFilter *string `json:"nameFilter,omitempty"`
}

type ServicesSearchResponse struct {
	Services []*Service `json:"services"`
}
