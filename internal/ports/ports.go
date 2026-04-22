package ports

import (
	"context"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/branch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/member"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/membercategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/opname"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/product"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/productcategory"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/purchase"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/sale"
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
	Update(ctx context.Context, item product.Product) error
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
