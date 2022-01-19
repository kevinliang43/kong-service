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
	SearchQuery *string `json:"searchQuery,omitempty"`
	Limit       *int64  `json:"limit,omitempty"`
	Offset      *int64  `json:"offset,omitempty"`
	SortType    *string `json:"sortType,omitempty"`
	FilterType  *string `json:"filterType,omitempty"`
}

type ServicesSearchResponse struct {
	Services   []*Service `json:"services"`
	NextOffset int64      `json:"nextOffset"`
}

func (ssr ServicesSearchRequest) GetNextOffset(responseSize int64) int64 {
	if ssr.Offset == nil {
		return responseSize
	} else {
		return *ssr.Offset + responseSize
	}
}
