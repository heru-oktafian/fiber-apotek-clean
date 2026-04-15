package member

type Member struct {
	ID               string
	MemberCategoryID string
	Points           int
}

type MemberCategory struct {
	ID                   string
	PointsConversionRate int
}
