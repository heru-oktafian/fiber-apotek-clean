package ports

import (
	"context"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/anotherincome"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/branch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/buyreturn"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/duplicatereceipt"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/expense"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/firststock"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/member"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/membercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/opname"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/product"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/productcategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/purchase"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/sale"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/salereturn"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/supplier"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/suppliercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/unit"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/user"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/userbranch"
)

type Clock interface{ Now() time.Time }

type IDGenerator interface{ New(prefix string) string }

type PasswordComparer interface {
	Compare(hashed string, plain string) error
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
}

type TokenManager interface {
	GenerateLoginToken(user user.User, expiresAt time.Time) (string, error)
	GenerateBranchToken(claims auth.Claims, expiresAt time.Time) (string, error)
	Parse(token string) (auth.Claims, time.Time, error)
}

type TokenBlacklist interface {
	Blacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type UserRepository interface {
	FindActiveByUsername(ctx context.Context, username string) (user.User, error)
	FindByID(ctx context.Context, id string) (user.User, error)
	ListUsers(ctx context.Context, req user.ListRequest) (user.ListResult, error)
	FindUserWithBranches(ctx context.Context, id string) (user.DetailWithBranches, error)
	CreateUser(ctx context.Context, item user.User) error
	UpdateUser(ctx context.Context, item user.User) error
	CreateUserBranch(ctx context.Context, userID, branchID string) error
}

type BranchRepository interface {
	FindBranchByID(ctx context.Context, id string) (branch.Branch, error)
	ListBranches(ctx context.Context, req branch.ListRequest) (branch.ListResult, error)
	CreateBranch(ctx context.Context, item branch.Branch) error
	DeleteBranch(ctx context.Context, id string) error
	BranchHasUsers(ctx context.Context, branchID string) (bool, error)
	UserHasBranch(ctx context.Context, userID, branchID string) (bool, error)
	ListUserBranches(ctx context.Context, userID string) ([]auth.UserBranch, error)
	FindProfile(ctx context.Context, userID, branchID string) (auth.Profile, error)
	ListAllUserBranches(ctx context.Context) ([]userbranch.Detail, error)
	FindUserBranchDetail(ctx context.Context, userID, branchID string) ([]userbranch.Detail, error)
}

type ProductRepository interface {
	Create(ctx context.Context, item product.Product) error
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	FindProductDetailByID(ctx context.Context, id, branchID string) (product.Product, error)
	ListProducts(ctx context.Context, branchID string, req product.ListRequest) (product.ListResult, error)
	Update(ctx context.Context, item product.Product) error
	DeleteProduct(ctx context.Context, id, branchID string) error
	GetSaleCombo(ctx context.Context, branchID, search string) ([]product.SaleComboItem, error)
	GetPurchaseCombo(ctx context.Context, branchID, search string) ([]product.PurchaseComboItem, error)
	GetOpnameCombo(ctx context.Context, branchID, search string) ([]product.OpnameComboItem, error)
}

type SupplierRepository interface {
	ListSuppliers(ctx context.Context, branchID string, req supplier.ListRequest) (supplier.ListResult, error)
	FindSupplierByID(ctx context.Context, id, branchID string) (supplier.Supplier, error)
	CreateSupplier(ctx context.Context, item supplier.Supplier) error
	UpdateSupplier(ctx context.Context, item supplier.Supplier) error
	DeleteSupplier(ctx context.Context, id, branchID string) error
	GetSupplierCombo(ctx context.Context, branchID, search string) ([]supplier.ComboItem, error)
}

type ProductCategoryRepository interface {
	ListProductCategories(ctx context.Context, branchID string, req productcategory.ListRequest) (productcategory.ListResult, error)
	FindProductCategoryByID(ctx context.Context, id uint, branchID string) (productcategory.ProductCategory, error)
	CreateProductCategory(ctx context.Context, item productcategory.ProductCategory) (productcategory.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, item productcategory.ProductCategory) error
	DeleteProductCategory(ctx context.Context, id uint, branchID string) error
	GetProductCategoryCombo(ctx context.Context, branchID, search string) ([]productcategory.ComboItem, error)
}

type SupplierCategoryRepository interface {
	ListSupplierCategories(ctx context.Context, branchID string, req suppliercategory.ListRequest) (suppliercategory.ListResult, error)
	FindSupplierCategoryByID(ctx context.Context, id uint, branchID string) (suppliercategory.SupplierCategory, error)
	CreateSupplierCategory(ctx context.Context, item suppliercategory.SupplierCategory) (suppliercategory.SupplierCategory, error)
	UpdateSupplierCategory(ctx context.Context, item suppliercategory.SupplierCategory) error
	DeleteSupplierCategory(ctx context.Context, id uint, branchID string) error
	GetSupplierCategoryCombo(ctx context.Context, branchID string) ([]suppliercategory.ComboItem, error)
}

type MemberCategoryRepository interface {
	ListMemberCategories(ctx context.Context, branchID string, req membercategory.ListRequest) (membercategory.ListResult, error)
	FindMemberCategoryByID(ctx context.Context, id uint, branchID string) (membercategory.MemberCategory, error)
	CreateMemberCategory(ctx context.Context, item membercategory.MemberCategory) (membercategory.MemberCategory, error)
	UpdateMemberCategory(ctx context.Context, item membercategory.MemberCategory) error
	DeleteMemberCategory(ctx context.Context, id uint, branchID string) error
	GetMemberCategoryCombo(ctx context.Context, branchID, search string) ([]membercategory.ComboItem, error)
}

type UnitRepository interface {
	FindUnitByID(ctx context.Context, id string) (unit.Unit, error)
	FindConversion(ctx context.Context, productID, initID, finalID, branchID string) (unit.Conversion, error)
	ListMasterUnits(ctx context.Context, branchID string, req unit.MasterUnitListRequest) (unit.MasterUnitListResult, error)
	FindMasterUnitByID(ctx context.Context, id, branchID string) (unit.MasterUnit, error)
	CreateMasterUnit(ctx context.Context, item unit.MasterUnit) error
	UpdateMasterUnit(ctx context.Context, item unit.MasterUnit) error
	DeleteMasterUnit(ctx context.Context, id, branchID string) error
	GetMasterUnitCombo(ctx context.Context, branchID, search string) ([]unit.MasterUnitComboItem, error)
}

type MemberRepository interface {
	FindMemberByID(ctx context.Context, id string) (member.Member, error)
	FindCategoryByID(ctx context.Context, id string) (member.MemberCategory, error)
	UpdatePoints(ctx context.Context, memberID string, points int) error
}

type PurchaseRepository interface {
	WithinTransaction(ctx context.Context, fn func(repo PurchaseTxRepository) error) error
	ListPurchases(ctx context.Context, branchID string, req purchase.ListRequest) (purchase.ListResult, error)
	FindPurchaseDetail(ctx context.Context, branchID, id string) (purchase.Detail, error)
	FindPurchaseByID(ctx context.Context, branchID, id string) (purchase.Purchase, error)
	FindPurchaseItemByID(ctx context.Context, id string) (purchase.Item, error)
	FindPurchaseItems(ctx context.Context, purchaseID string) ([]purchase.Item, error)
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	UpdatePurchaseHeader(ctx context.Context, item purchase.Purchase) error
	UpdatePurchaseItem(ctx context.Context, item purchase.Item) error
	CreatePurchaseItem(ctx context.Context, item purchase.Item) error
	DeletePurchaseItem(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, item product.Product) error
	DeletePurchaseHeader(ctx context.Context, branchID, id string) error
}

type PurchaseTxRepository interface {
	FindProduct(ctx context.Context, id string) (product.Product, error)
	FindUnit(ctx context.Context, id string) (unit.Unit, error)
	FindConversion(ctx context.Context, productID, initID, finalID, branchID string) (unit.Conversion, error)
	CreatePurchase(ctx context.Context, item purchase.Purchase) error
	CreatePurchaseItems(ctx context.Context, items []purchase.Item) error
	UpdateProduct(ctx context.Context, item product.Product) error
	CreateTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time) error
}

type SaleRepository interface {
	WithinTransactionSale(ctx context.Context, fn func(repo SaleTxRepository) error) error
	ListSales(ctx context.Context, branchID string, req sale.ListRequest) (sale.ListResult, error)
	FindSaleDetail(ctx context.Context, branchID, id string) (sale.Detail, error)
	FindSaleByID(ctx context.Context, branchID, id string) (sale.Sale, error)
	FindSaleItemByID(ctx context.Context, id string) (sale.Item, error)
	FindSaleItems(ctx context.Context, saleID string) ([]sale.Item, error)
	FindMemberByID(ctx context.Context, memberID string) (member.Member, error)
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	CreateSaleItem(ctx context.Context, item sale.Item) error
	UpdateSaleItem(ctx context.Context, item sale.Item) error
	UpdateSaleHeader(ctx context.Context, item sale.Sale) error
	UpdateTransactionReport(ctx context.Context, id string, total int, payment string, updatedAt time.Time) error
	AdjustDailyProfit(ctx context.Context, reportDate time.Time, userID string, branchID string, totalDelta int, profitDelta int, now time.Time) error
	DeleteSaleItem(ctx context.Context, id string) error
	DeleteSaleHeader(ctx context.Context, branchID, id string) error
	DeleteSaleItems(ctx context.Context, saleID string) error
	DeleteTransactionReport(ctx context.Context, id string, txType string) error
}

type SaleTxRepository interface {
	FindProduct(ctx context.Context, id string) (product.Product, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	CreateSale(ctx context.Context, item sale.Sale) error
	CreateSaleItems(ctx context.Context, items []sale.Item) error
	CreateTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time) error
	UpsertDailyProfit(ctx context.Context, reportDate time.Time, userID string, branchID string, totalSales int, profitEstimate int, now time.Time) error
	FindBranch(ctx context.Context, branchID string) (branch.Branch, error)
	UpdateBranchQuota(ctx context.Context, branchID string, quota int) error
	FindMember(ctx context.Context, memberID string) (member.Member, error)
	FindMemberCategory(ctx context.Context, categoryID string) (member.MemberCategory, error)
	UpdateMemberPoints(ctx context.Context, memberID string, points int) error
}

type BuyReturnRepository interface {
	ListBuyReturns(ctx context.Context, branchID string, req buyreturn.ListRequest) (buyreturn.ListResult, error)
	FindBuyReturnByID(ctx context.Context, branchID, id string) (buyreturn.BuyReturn, error)
	FindBuyReturnItems(ctx context.Context, buyReturnID string) ([]buyreturn.Item, error)
	CreateBuyReturn(ctx context.Context, item buyreturn.BuyReturn) error
	CreateBuyReturnItems(ctx context.Context, items []buyreturn.Item) error
	ListPurchaseReturnSources(ctx context.Context, branchID, search, month string) ([]buyreturn.PurchaseComboItem, error)
	ListPurchaseReturnableItems(ctx context.Context, purchaseID string) ([]buyreturn.ReturnableItem, error)
	FindPurchaseByID(ctx context.Context, branchID, id string) (purchase.Purchase, error)
	FindPurchaseItemByPurchaseAndProduct(ctx context.Context, purchaseID, productID string) (purchase.Item, error)
	SumBuyReturnedQty(ctx context.Context, purchaseID, productID string) (int, error)
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	FindUnit(ctx context.Context, id string) (unit.Unit, error)
	FindConversion(ctx context.Context, productID, initID, finalID, branchID string) (unit.Conversion, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	CreateTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time) error
}

type SaleReturnRepository interface {
	ListSaleReturns(ctx context.Context, branchID string, req salereturn.ListRequest) (salereturn.ListResult, error)
	FindSaleReturnByID(ctx context.Context, branchID, id string) (salereturn.SaleReturn, error)
	FindSaleReturnItems(ctx context.Context, saleReturnID string) ([]salereturn.Item, error)
	CreateSaleReturn(ctx context.Context, item salereturn.SaleReturn) error
	CreateSaleReturnItems(ctx context.Context, items []salereturn.Item) error
	ListSaleReturnSources(ctx context.Context, branchID, search, month string) ([]salereturn.SaleComboItem, error)
	ListSaleReturnableItems(ctx context.Context, saleID string) ([]salereturn.ReturnableItem, error)
	FindSaleByID(ctx context.Context, branchID, id string) (sale.Sale, error)
	FindSaleItemBySaleAndProduct(ctx context.Context, saleID, productID string) (sale.Item, error)
	SumSaleReturnedQty(ctx context.Context, saleID, productID string) (int, error)
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	CreateTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time) error
}

type DuplicateReceiptRepository interface {
	WithinTransactionDuplicateReceipt(ctx context.Context, fn func(repo DuplicateReceiptTxRepository) error) error
	ListDuplicateReceipts(ctx context.Context, branchID string, req duplicatereceipt.ListRequest) (duplicatereceipt.ListResult, error)
	FindDuplicateReceiptByID(ctx context.Context, branchID, id string) (duplicatereceipt.DuplicateReceipt, error)
	UpdateDuplicateReceipt(ctx context.Context, item duplicatereceipt.DuplicateReceipt) error
	FindDuplicateReceiptItems(ctx context.Context, duplicateReceiptID string) ([]duplicatereceipt.Item, error)
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	UpdateTransactionReport(ctx context.Context, id string, total int, payment string, updatedAt time.Time) error
	DeleteTransactionReport(ctx context.Context, id string, txType string) error
	AdjustDailyProfit(ctx context.Context, reportDate time.Time, userID string, branchID string, totalDelta int, profitDelta int, now time.Time) error
	FindMemberByID(ctx context.Context, id string) (member.Member, error)
}

type DuplicateReceiptTxRepository interface {
	FindProduct(ctx context.Context, id string) (product.Product, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	CreateDuplicateReceipt(ctx context.Context, item duplicatereceipt.DuplicateReceipt) error
	CreateDuplicateReceiptItems(ctx context.Context, items []duplicatereceipt.Item) error
	CreateTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time) error
	UpsertDailyProfit(ctx context.Context, reportDate time.Time, userID string, branchID string, totalSales int, profitEstimate int, now time.Time) error
	FindMember(ctx context.Context, memberID string) (member.Member, error)
	FindMemberCategory(ctx context.Context, categoryID string) (member.MemberCategory, error)
	UpdateMemberPoints(ctx context.Context, memberID string, points int) error
	DeleteDuplicateReceipt(ctx context.Context, branchID, id string) error
	DeleteDuplicateReceiptItems(ctx context.Context, duplicateReceiptID string) error
	DeleteTransactionReport(ctx context.Context, id string, txType string) error
}

type AnotherIncomeRepository interface {
	ListAnotherIncomes(ctx context.Context, branchID string, req anotherincome.ListRequest) (anotherincome.ListResult, error)
	FindAnotherIncomeByID(ctx context.Context, branchID, id string) (anotherincome.AnotherIncome, error)
	CreateAnotherIncome(ctx context.Context, item anotherincome.AnotherIncome) error
	UpdateAnotherIncome(ctx context.Context, item anotherincome.AnotherIncome) error
	DeleteAnotherIncome(ctx context.Context, branchID, id string) error
	UpsertTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time, updatedAt time.Time) error
	DeleteTransactionReport(ctx context.Context, id string, txType string) error
}

type ExpenseRepository interface {
	ListExpenses(ctx context.Context, branchID string, req expense.ListRequest) (expense.ListResult, error)
	FindExpenseByID(ctx context.Context, branchID, id string) (expense.Expense, error)
	CreateExpense(ctx context.Context, item expense.Expense) error
	UpdateExpense(ctx context.Context, item expense.Expense) error
	DeleteExpense(ctx context.Context, branchID, id string) error
	UpsertTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time, updatedAt time.Time) error
	DeleteTransactionReport(ctx context.Context, id string, txType string) error
}

type FirstStockRepository interface {
	ListFirstStocks(ctx context.Context, branchID string, req firststock.ListRequest) (firststock.ListResult, error)
	FindFirstStockByID(ctx context.Context, branchID, id string) (firststock.FirstStock, error)
	FindFirstStockItems(ctx context.Context, firstStockID string) ([]firststock.Item, error)
	FindFirstStockItemByID(ctx context.Context, id string) (firststock.Item, error)
	CreateFirstStock(ctx context.Context, item firststock.FirstStock) error
	UpdateFirstStock(ctx context.Context, item firststock.FirstStock) error
	DeleteFirstStock(ctx context.Context, branchID, id string) error
	CreateFirstStockItem(ctx context.Context, item firststock.Item) error
	UpdateFirstStockItem(ctx context.Context, item firststock.Item) error
	DeleteFirstStockItem(ctx context.Context, id string) error
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	FindUnit(ctx context.Context, id string) (unit.Unit, error)
	FindConversion(ctx context.Context, productID, initID, finalID, branchID string) (unit.Conversion, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	RecalculateFirstStockTotal(ctx context.Context, firstStockID string) (int, error)
	UpsertTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time, updatedAt time.Time) error
	DeleteTransactionReport(ctx context.Context, id string, txType string) error
}

type OpnameRepository interface {
	CreateOpname(ctx context.Context, item opname.Opname) error
	FindOpnameByID(ctx context.Context, id string) (opname.Opname, error)
	FindOpnameItems(ctx context.Context, opnameID string) ([]opname.Item, error)
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	UpdateProduct(ctx context.Context, item product.Product) error
	UpdateOpnameTotal(ctx context.Context, opnameID string, total int) error
	RecalculateOpnameTotal(ctx context.Context, opnameID string) (int, error)
	WithinOpnameTransaction(ctx context.Context, fn func(repo OpnameTxRepository) error) error
}

type OpnameTxRepository interface {
	FindOpnameByID(ctx context.Context, id string) (opname.Opname, error)
	FindProductByID(ctx context.Context, id string) (product.Product, error)
	FindOpnameItemByOpnameAndProduct(ctx context.Context, opnameID, productID string) (opname.Item, error)
	CreateOpnameItem(ctx context.Context, item opname.Item) error
	UpdateOpnameItem(ctx context.Context, item opname.Item) error
	UpdateProduct(ctx context.Context, item product.Product) error
	UpdateOpnameTotal(ctx context.Context, opnameID string, total int) error
	RecalculateOpnameTotal(ctx context.Context, opnameID string) (int, error)
}
