package unit

type MasterUnit struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BranchID string `json:"branch_id,omitempty"`
}

type MasterUnitCreateRequest struct {
	Name string `json:"name"`
}

type MasterUnitListRequest struct {
	Search string
	Page   int
	Limit  int
}

type MasterUnitListMeta struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	TotalData int    `json:"total_data"`
	LastPage  int    `json:"last_page"`
}

type MasterUnitListResult struct {
	Items []MasterUnit
	Meta  MasterUnitListMeta
}

type MasterUnitComboItem struct {
	UnitID   string `json:"unit_id"`
	UnitName string `json:"unit_name"`
}
