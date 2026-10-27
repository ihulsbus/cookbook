package models

type PaginationRequest struct {
	Page  int `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" json:"limit" binding:"omitempty,min=1"`
}

type PaginationMetadata struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

type PaginatedResponse[T any] struct {
	Data       []T                `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

func NormalizePagination(page int, limit int) PaginationRequest {
	const defaultPage = 1
	const defaultLimit = 25
	const maxLimit = 100

	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return PaginationRequest{
		Page:  page,
		Limit: limit,
	}
}

func NewPaginationMetadata(page int, limit int, total int64) PaginationMetadata {
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return PaginationMetadata{
		Page:        page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}
