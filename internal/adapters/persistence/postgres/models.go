package postgres

import "time"

type UserModel struct {
	ID         string `gorm:"column:id;primaryKey"`
	Name       string `gorm:"column:name"`
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	UserRole   string `gorm:"column:user_role"`
	UserStatus string `gorm:"column:user_status"`
}

func (UserModel) TableName() string { return "users" }

type BranchModel struct {
	ID               string    `gorm:"column:id;primaryKey"`
	BranchName       string    `gorm:"column:branch_name"`
	Address          string    `gorm:"column:address"`
	Phone            string    `gorm:"column:phone"`
	Email            string    `gorm:"column:email"`
	SIAID            string    `gorm:"column:sia_id"`
	SIAName          string    `gorm:"column:sia_name"`
	PSAID            string    `gorm:"column:psa_id"`
	PSAName          string    `gorm:"column:psa_name"`
	SIPA             string    `gorm:"column:sipa"`
	SIPAName         string    `gorm:"column:sipa_name"`
	APINGID          string    `gorm:"column:aping_id"`
	APINGName        string    `gorm:"column:aping_name"`
	BankName         string    `gorm:"column:bank_name"`
	AccountName      string    `gorm:"column:account_name"`
	AccountNumber    string    `gorm:"column:account_number"`
	TaxPercentage    int       `gorm:"column:tax_percentage"`
	JournalMethod    string    `gorm:"column:journal_method"`
	BranchStatus     string    `gorm:"column:branch_status"`
	LicenseDate      time.Time `gorm:"column:license_date"`
	DefaultMember    string    `gorm:"column:default_member"`
	Quota            int       `gorm:"column:quota"`
	SubscriptionType string    `gorm:"column:subscription_type"`
	RealAsset        string    `gorm:"column:real_asset"`
}

func (BranchModel) TableName() string { return "branches" }

type UserBranchModel struct {
	UserID   string `gorm:"column:user_id"`
	BranchID string `gorm:"column:branch_id"`
}

func (UserBranchModel) TableName() string { return "user_branches" }

type ProductModel struct {
	ID                string    `gorm:"column:id;primaryKey"`
	SKU               string    `gorm:"column:sku"`
	Name              string    `gorm:"column:name"`
	Alias             string    `gorm:"column:alias"`
	Description       string    `gorm:"column:description"`
	Ingredient        string    `gorm:"column:ingredient"`
	Dosage            string    `gorm:"column:dosage"`
	SideAffection     string    `gorm:"column:side_affection"`
	BranchID          string    `gorm:"column:branch_id"`
	UnitID            string    `gorm:"column:unit_id"`
	Stock             int       `gorm:"column:stock"`
	PurchasePrice     int       `gorm:"column:purchase_price"`
	SalesPrice        int       `gorm:"column:sales_price"`
	AlternatePrice    int       `gorm:"column:alternate_price"`
	ProductCategoryID uint      `gorm:"column:product_category_id"`
	ExpiredDate       time.Time `gorm:"column:expired_date"`
}

func (ProductModel) TableName() string { return "products" }

type SupplierModel struct {
	ID                 string `gorm:"column:id;primaryKey"`
	Name               string `gorm:"column:name"`
	Phone              string `gorm:"column:phone"`
	Address            string `gorm:"column:address"`
	PIC                string `gorm:"column:pic"`
	SupplierCategoryID uint   `gorm:"column:supplier_category_id"`
	BranchID           string `gorm:"column:branch_id"`
}

func (SupplierModel) TableName() string { return "suppliers" }

type SupplierCategoryModel struct {
	ID       uint   `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	BranchID string `gorm:"column:branch_id"`
}

func (SupplierCategoryModel) TableName() string { return "supplier_categories" }

type ProductCategoryModel struct {
	ID       uint   `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	BranchID string `gorm:"column:branch_id"`
}

func (ProductCategoryModel) TableName() string { return "product_categories" }

type UnitModel struct {
	ID       string `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	BranchID string `gorm:"column:branch_id"`
}

func (UnitModel) TableName() string { return "units" }

type UnitConversionModel struct {
	ProductID string `gorm:"column:product_id"`
	InitID    string `gorm:"column:init_id"`
	FinalID   string `gorm:"column:final_id"`
	ValueConv int    `gorm:"column:value_conv"`
	BranchID  string `gorm:"column:branch_id"`
}

func (UnitConversionModel) TableName() string { return "unit_conversions" }

type PurchaseModel struct {
	ID            string    `gorm:"column:id;primaryKey"`
	SupplierID    string    `gorm:"column:supplier_id"`
	PurchaseDate  time.Time `gorm:"column:purchase_date"`
	BranchID      string    `gorm:"column:branch_id"`
	UserID        string    `gorm:"column:user_id"`
	Payment       string    `gorm:"column:payment"`
	TotalPurchase int       `gorm:"column:total_purchase"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (PurchaseModel) TableName() string { return "purchases" }

type PurchaseItemModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	PurchaseID  string    `gorm:"column:purchase_id"`
	ProductID   string    `gorm:"column:product_id"`
	UnitID      string    `gorm:"column:unit_id"`
	Price       int       `gorm:"column:price"`
	Qty         int       `gorm:"column:qty"`
	SubTotal    int       `gorm:"column:sub_total"`
	ExpiredDate time.Time `gorm:"column:expired_date"`
}

func (PurchaseItemModel) TableName() string { return "purchase_items" }

type SaleModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	MemberID       string    `gorm:"column:member_id"`
	UserID         string    `gorm:"column:user_id"`
	BranchID       string    `gorm:"column:branch_id"`
	Payment        string    `gorm:"column:payment"`
	Discount       int       `gorm:"column:discount"`
	TotalSale      int       `gorm:"column:total_sale"`
	ProfitEstimate int       `gorm:"column:profit_estimate"`
	SaleDate       time.Time `gorm:"column:sale_date"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (SaleModel) TableName() string { return "sales" }

type SaleItemModel struct {
	ID        string `gorm:"column:id;primaryKey"`
	SaleID    string `gorm:"column:sale_id"`
	ProductID string `gorm:"column:product_id"`
	Price     int    `gorm:"column:price"`
	Qty       int    `gorm:"column:qty"`
	SubTotal  int    `gorm:"column:sub_total"`
}

func (SaleItemModel) TableName() string { return "sale_items" }

type TransactionReportModel struct {
	ID              string    `gorm:"column:id;primaryKey"`
	TransactionType string    `gorm:"column:transaction_type"`
	UserID          string    `gorm:"column:user_id"`
	BranchID        string    `gorm:"column:branch_id"`
	Total           int       `gorm:"column:total"`
	Payment         string    `gorm:"column:payment"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (TransactionReportModel) TableName() string { return "transaction_reports" }

type AnotherIncomeModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Description string    `gorm:"column:description"`
	IncomeDate  time.Time `gorm:"column:income_date"`
	BranchID    string    `gorm:"column:branch_id"`
	TotalIncome int       `gorm:"column:total_income"`
	Payment     string    `gorm:"column:payment"`
	UserID      string    `gorm:"column:user_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (AnotherIncomeModel) TableName() string { return "another_incomes" }

type ExpenseModel struct {
	ID           string    `gorm:"column:id;primaryKey"`
	Description  string    `gorm:"column:description"`
	ExpenseDate  time.Time `gorm:"column:expense_date"`
	BranchID     string    `gorm:"column:branch_id"`
	TotalExpense int       `gorm:"column:total_expense"`
	Payment      string    `gorm:"column:payment"`
	UserID       string    `gorm:"column:user_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (ExpenseModel) TableName() string { return "expenses" }

type DailyProfitReportModel struct {
	ID             string    `gorm:"column:id;primaryKey"`
	ReportDate     time.Time `gorm:"column:report_date"`
	UserID         string    `gorm:"column:user_id"`
	BranchID       string    `gorm:"column:branch_id"`
	TotalSales     int       `gorm:"column:total_sales"`
	ProfitEstimate int       `gorm:"column:profit_estimate"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (DailyProfitReportModel) TableName() string { return "daily_profit_reports" }

type OpnameModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Description string    `gorm:"column:description"`
	BranchID    string    `gorm:"column:branch_id"`
	UserID      string    `gorm:"column:user_id"`
	OpnameDate  time.Time `gorm:"column:opname_date"`
	TotalOpname int       `gorm:"column:total_opname"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (OpnameModel) TableName() string { return "opnames" }

type OpnameItemModel struct {
	ID            string    `gorm:"column:id;primaryKey"`
	OpnameID      string    `gorm:"column:opname_id"`
	ProductID     string    `gorm:"column:product_id"`
	Qty           int       `gorm:"column:qty"`
	QtyExist      int       `gorm:"column:qty_exist"`
	Price         int       `gorm:"column:price"`
	SubTotal      int       `gorm:"column:sub_total"`
	SubTotalExist int       `gorm:"column:sub_total_exist"`
	ExpiredDate   time.Time `gorm:"column:expired_date"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (OpnameItemModel) TableName() string { return "opname_items" }

type MemberModel struct {
	ID               string `gorm:"column:id;primaryKey"`
	MemberCategoryID string `gorm:"column:member_category_id"`
	Points           int    `gorm:"column:points"`
}

func (MemberModel) TableName() string { return "members" }

type MemberCategoryModel struct {
	ID                   uint   `gorm:"column:id;primaryKey"`
	Name                 string `gorm:"column:name"`
	PointsConversionRate int    `gorm:"column:points_conversion_rate"`
	BranchID             string `gorm:"column:branch_id"`
}

func (MemberCategoryModel) TableName() string { return "member_categories" }
