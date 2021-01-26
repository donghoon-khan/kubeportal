package dataselect

var NoPagination = NewPaginationQuery(-1, -1)

var EmptyPagination = NewPaginationQuery(0, 0)

var DefaultPagination = NewPaginationQuery(10, 0)

type PaginationQuery struct {
	ItemPerPage int
	Page        int
}

func NewPaginationQuery(itemsPerPage, page int) *PaginationQuery {
	return &PaginationQuery{itemsPerPage, page}
}

func (p *PaginationQuery) IsValidPagination() bool {
	return p.ItemPerPage >= 0 && p.Page >= 0
}

func (p *PaginationQuery) IsPageAvailable(itemsCount, startingIndex int) bool {
	return itemsCount > startingIndex && p.ItemPerPage > 0
}

func (p *PaginationQuery) GetPaginationSettings(itemsCount int) (startIndex int, endIndex int) {
	startIndex = p.ItemPerPage * p.Page
	endIndex = startIndex + p.ItemPerPage

	if endIndex > itemsCount {
		endIndex = itemsCount
	}

	return startIndex, endIndex
}
