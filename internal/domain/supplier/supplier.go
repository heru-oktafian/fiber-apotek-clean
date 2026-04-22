package supplier

type Supplier struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	PIC                string `json:"pic"`
	SupplierCategoryID uint   `json:"supplier_category_id"`
	SupplierCategory   string `json:"supplier_category,omitempty"`
	BranchID           string `json:"branch_id,omitempty"`
}

type CreateRequest struct {
	Name               string `json:"name"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	PIC                string `json:"pic"`
	SupplierCategoryID uint   `json:"supplier_category_id"`
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
	Items []Supplier
	Meta  ListMeta
}

type ComboItem struct {
	SupplierID   string `json:"supplier_id"`
	SupplierName string `json:"supplier_name"`
}
