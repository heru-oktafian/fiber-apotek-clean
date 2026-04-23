package member

type Member struct {
	ID               string `json:"id"`
	Name             string `json:"name,omitempty"`
	Phone            string `json:"phone,omitempty"`
	Address          string `json:"address,omitempty"`
	MemberCategoryID string `json:"member_category_id"`
	MemberCategory   string `json:"member_category,omitempty"`
	Points           int    `json:"points"`
	BranchID         string `json:"branch_id,omitempty"`
}

type MemberCategory struct {
	ID                   string `json:"id"`
	PointsConversionRate int    `json:"points_conversion_rate"`
}

type CreateRequest struct {
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	Address          string `json:"address"`
	MemberCategoryID uint   `json:"member_category_id"`
	Points           int    `json:"points"`
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
	Items []Member
	Meta  ListMeta
}

type ComboItem struct {
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
}
