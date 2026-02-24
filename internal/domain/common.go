package domain

type Paginated[T any] struct {
	Items    []T    `json:"items"`
	Total    uint64 `json:"total"`
	Page     uint   `json:"page"`
	PageSize uint   `json:"page_size"`
}

type Filter struct {
	Search   *string `form:"search"`
	Page     uint    `form:"page,default=1"`
	PageSize uint    `form:"page_size,default=10"`
}
