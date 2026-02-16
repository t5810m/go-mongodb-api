package helpers

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasMore    bool  `json:"has_more"`
}

// GetSkip calculates the skip value for MongoDB
func (p *Pagination) GetSkip() int {
	return (p.Page - 1) * p.Limit
}

// SetTotal calculates TotalPages and HasMore based on total count
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	p.TotalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))
	p.HasMore = int64(p.Page*p.Limit) < total
}

func NewPagination(page, limit int) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return &Pagination{Page: page, Limit: limit}
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
}
