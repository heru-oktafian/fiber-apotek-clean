package auth

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResult struct {
	Token string
}

type BranchSelectionRequest struct {
	BranchID string `json:"branch_id" validate:"required"`
}

type Claims struct {
	Subject          string
	Name             string
	BranchID         string
	UserRole         string
	DefaultMember    string
	Quota            int
	SubscriptionType string
	RealAsset        string
}

type UserBranch struct {
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	SIAName    string `json:"sia_name"`
	SIPAName   string `json:"sipa_name"`
	Phone      string `json:"phone"`
}

type Profile struct {
	UserID        string `json:"user_id"`
	ProfileName   string `json:"profile_name"`
	BranchID      string `json:"branch_id"`
	BranchName    string `json:"branch_name"`
	Address       string `json:"address"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	SIAID         string `json:"sia_id"`
	SIAName       string `json:"sia_name"`
	PSAID         string `json:"psa_id"`
	PSAName       string `json:"psa_name"`
	SIPA          string `json:"sipa"`
	SIPAName      string `json:"sipa_name"`
	APINGID       string `json:"aping_id"`
	APINGName     string `json:"aping_name"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	TaxPercentage int    `json:"tax_percentage"`
	JournalMethod string `json:"journal_method"`
	BranchStatus  string `json:"branch_status"`
	LicenseDate   string `json:"license_date"`
	DefaultMember string `json:"default_member"`
	MemberName    string `json:"member_name"`
}
