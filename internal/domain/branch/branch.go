package branch

type Branch struct {
	ID               string
	Name             string
	DefaultMemberID  string
	Quota            int
	SubscriptionType string
	RealAsset        string
}
