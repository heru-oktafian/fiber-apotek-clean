package suppliercategory

type SupplierCategory struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	BranchID string `json:"branch_id,omitempty"`
}

type CreateRequest struct {
	Name string `json:"name"`
}

type ListRequest struct {
	Search string
	Page   int
	Limit  int
}

type ListMeta struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	TotalData int    `json:"total_data"`
	LastPage  int    `json:"last_page"`
}

type ListResult struct {
	Items []SupplierCategory
	Meta  ListMeta
}

type ComboItem struct {
	SupplierCategoryID   uint   `json:"supplier_category_id"`
	SupplierCategoryName string `json:"supplier_category_name"`
}
