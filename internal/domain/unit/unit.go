package unit

type Unit struct {
	ID   string
	Name string
}

type Conversion struct {
	ProductID string
	InitID    string
	FinalID   string
	Value     int
	BranchID  string
}
