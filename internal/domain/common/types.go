package common

type UserRole string

type PaymentStatus string

const (
	RoleAdministrator UserRole = "administrator"
	RoleOperator      UserRole = "operator"
	RoleCashier       UserRole = "cashier"
	RoleFinance       UserRole = "finance"
	RoleSuperadmin    UserRole = "superadmin"
)

const (
	PaymentCash PaymentStatus = "paid_by_cash"
)
