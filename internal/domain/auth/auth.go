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
