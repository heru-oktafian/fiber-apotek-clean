package membercategory

type MemberCategory struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	PointsConversionRate int    `json:"points_conversion_rate"`
	BranchID             string `json:"branch_id,omitempty"`
}

type CreateRequest struct {
	Name                 string `json:"name"`
	PointsConversionRate int    `json:"points_conversion_rate"`
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
	Items []MemberCategory
	Meta  ListMeta
}

type ComboItem struct {
	MemberCategoryID   uint   `json:"member_category_id"`
	MemberCategoryName string `json:"member_category_name"`
}
