package pagination

type PaginationMeta struct {
	TotalItems      int  `json:"total_items"`
	CurrentPage     int  `json:"current_page"`
	ItemsPerPage    int  `json:"items_per_page"`
	TotalPages      int  `json:"total_pages"`
	HasPreviousPage bool `json:"has_previous_page"`
	HasNextPage     bool `json:"has_next_page"`
	IsFirstPage     bool `json:"is_first_page"`
	IsLastPage      bool `json:"is_last_page"`
}

type Paginated[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func newPaginationMeta(totalItems, currentPage, itemsPerPage int) PaginationMeta {
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	hasPreviousPage := currentPage > 1
	hasNextPage := currentPage < totalPages
	isFirstPage := currentPage == 1
	isLastPage := currentPage == totalPages

	return PaginationMeta{
		TotalItems:      totalItems,
		CurrentPage:     currentPage,
		ItemsPerPage:    itemsPerPage,
		TotalPages:      totalPages,
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		IsFirstPage:     isFirstPage,
		IsLastPage:      isLastPage,
	}
}

func New[T any](data []T, totalItems, currentPage, itemsPerPage int) Paginated[T] {
	meta := newPaginationMeta(totalItems, currentPage, itemsPerPage)
	return Paginated[T]{Data: data, Meta: meta}
}
