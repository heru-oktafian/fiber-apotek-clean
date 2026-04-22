package branch

type Branch struct {
	ID               string `json:"id"`
	BranchName       string `json:"branch_name"`
	Address          string `json:"address"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	SIAID            string `json:"sia_id"`
	SIAName          string `json:"sia_name"`
	PSAID            string `json:"psa_id"`
	PSAName          string `json:"psa_name"`
	SIPA             string `json:"sipa"`
	SIPAName         string `json:"sipa_name"`
	APINGID          string `json:"aping_id"`
	APINGName        string `json:"aping_name"`
	BankName         string `json:"bank_name"`
	AccountName      string `json:"account_name"`
	AccountNumber    string `json:"account_number"`
	TaxPercentage    int    `json:"tax_percentage"`
	JournalMethod    string `json:"journal_method"`
	BranchStatus     string `json:"branch_status"`
	LicenseDate      string `json:"license_date"`
	DefaultMemberID  string `json:"default_member"`
	SubscriptionType string `json:"subscription_type"`
	Quota            int    `json:"quota"`
	RealAsset        string `json:"real_asset"`
}

type ListRequest struct {
	Search string
	Page   int
	Limit  int
}

type ListResult struct {
	Items []Branch
	Meta  ListMeta
}

type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	LastPage  int `json:"last_page"`
}

type CreateRequest struct {
	BranchName       string `json:"branch_name"`
	Address          string `json:"address"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	SIAID            string `json:"sia_id"`
	SIAName          string `json:"sia_name"`
	PSAID            string `json:"psa_id"`
	PSAName          string `json:"psa_name"`
	SIPA             string `json:"sipa"`
	SIPAName         string `json:"sipa_name"`
	APINGID          string `json:"aping_id"`
	APINGName        string `json:"aping_name"`
	BankName         string `json:"bank_name"`
	AccountName      string `json:"account_name"`
	AccountNumber    string `json:"account_number"`
	TaxPercentage    int    `json:"tax_percentage"`
	JournalMethod    string `json:"journal_method"`
	BranchStatus     string `json:"branch_status"`
	LicenseDate      string `json:"license_date"`
	DefaultMemberID  string `json:"default_member"`
	SubscriptionType string `json:"subscription_type"`
	Quota            int    `json:"quota"`
	RealAsset        string `json:"real_asset"`
}
