package unit

type ConversionMaster struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name,omitempty"`
	InitID      string `json:"init_id"`
	InitName    string `json:"init_name,omitempty"`
	FinalID     string `json:"final_id"`
	FinalName   string `json:"final_name,omitempty"`
	ValueConv   int    `json:"value_conv"`
	BranchID    string `json:"branch_id,omitempty"`
}

type ConversionCreateRequest struct {
	ProductID string `json:"product_id"`
	InitID    string `json:"init_id"`
	FinalID   string `json:"final_id"`
	ValueConv int    `json:"value_conv"`
}

type ConversionListRequest struct {
	Search string
	Page   int
	Limit  int
}

type ConversionListMeta struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	TotalData int    `json:"total_data"`
	LastPage  int    `json:"last_page"`
}

type ConversionListResult struct {
	Items []ConversionMaster
	Meta  ConversionListMeta
}
