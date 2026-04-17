package userbranch

type Detail struct {
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	SIAName    string `json:"sia_name"`
	SIPAName   string `json:"sipa_name"`
	Phone      string `json:"phone"`
}
