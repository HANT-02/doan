package repositories

type Paging struct {
	Page  uint64 `json:"page" form:"page"`
	Limit uint64 `json:"limit" form:"limit"`
}

type Sorting struct {
	Field string `json:"field" form:"field"`
	Order string `json:"order" form:"order"`
}
type Meta struct {
	ItemsPerPage uint64 `json:"items_per_page"`
	TotalItems   uint64 `json:"total_items"`
	CurrentPage  uint64 `json:"current_page"`
	TotalPages   uint64 `json:"total_pages"`
}
type Pagination[T any] struct {
	Data []*T `json:"data"`
	Meta Meta `json:"meta"`
}

func (meta Meta) ToDto() Meta {
	return Meta{
		ItemsPerPage: uint64(int(meta.ItemsPerPage)),
		TotalItems:   uint64(int64(meta.TotalItems)),
		CurrentPage:  uint64(int(meta.CurrentPage)),
		TotalPages:   uint64(int(meta.TotalPages)),
	}
}

func NewMeta(paging *Paging, totalItems uint64) Meta {
	totalPage := uint64(1)
	itemsPerPage := totalItems
	currentPage := uint64(1)
	if paging != nil && paging.Limit > 0 && totalItems > 0 {
		totalPage = totalItems / paging.Limit
		if totalItems%paging.Limit > 0 {
			totalPage++
		}
		itemsPerPage = paging.Limit
		currentPage = paging.Page
	}
	return Meta{
		ItemsPerPage: itemsPerPage,
		TotalItems:   totalItems,
		CurrentPage:  currentPage,
		TotalPages:   totalPage,
	}
}
