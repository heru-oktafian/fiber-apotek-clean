package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/branch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/member"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/opname"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/product"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/purchase"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/sale"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/unit"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/user"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/userbranch"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"gorm.io/gorm"
)

type Repositories struct {
	DB *gorm.DB
}

func (r Repositories) FindActiveByUsername(ctx context.Context, username string) (user.User, error) {
	var m UserModel
	if err := r.DB.WithContext(ctx).Where("username = ? AND user_status = 'active'", username).First(&m).Error; err != nil {
		return user.User{}, err
	}
	return user.User{ID: m.ID, Name: m.Name, Username: m.Username, Password: m.Password, Role: common.UserRole(m.UserRole), Status: m.UserStatus}, nil
}

func (r Repositories) FindByID(ctx context.Context, id string) (user.User, error) {
	var m UserModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return user.User{}, err
	}
	return user.User{ID: m.ID, Name: m.Name, Username: m.Username, Password: m.Password, Role: common.UserRole(m.UserRole), Status: m.UserStatus}, nil
}

func (r Repositories) ListUsers(ctx context.Context, req user.ListRequest) (user.ListResult, error) {
	query := r.DB.WithContext(ctx).Model(&UserModel{})
	if req.Search != "" {
		like := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("username ILIKE ? OR name ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return user.ListResult{}, err
	}

	var models []UserModel
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("name ASC").Find(&models).Error; err != nil {
		return user.ListResult{}, err
	}

	items := make([]user.User, 0, len(models))
	for _, model := range models {
		items = append(items, user.User{ID: model.ID, Name: model.Name, Username: model.Username, Role: common.UserRole(model.UserRole), Status: model.UserStatus})
	}

	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}

	return user.ListResult{
		Items: items,
		Meta: user.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			Search:    req.Search,
			TotalData: int(total),
			LastPage:  lastPage,
		},
	}, nil
}

func (r Repositories) FindUserWithBranches(ctx context.Context, id string) (user.DetailWithBranches, error) {
	usr, err := r.FindByID(ctx, id)
	if err != nil {
		return user.DetailWithBranches{}, err
	}
	usr.Password = ""

	var userBranches []UserBranchModel
	if err := r.DB.WithContext(ctx).Where("user_id = ?", id).Find(&userBranches).Error; err != nil {
		return user.DetailWithBranches{}, err
	}

	branchIDs := make([]string, 0, len(userBranches))
	for _, item := range userBranches {
		branchIDs = append(branchIDs, item.BranchID)
	}

	branchDetails := make([]user.BranchDetail, 0)
	if len(branchIDs) > 0 {
		var branches []BranchModel
		if err := r.DB.WithContext(ctx).Where("id IN ?", branchIDs).Find(&branches).Error; err != nil {
			return user.DetailWithBranches{}, err
		}
		for _, item := range branches {
			branchDetails = append(branchDetails, user.BranchDetail{
				BranchID:   item.ID,
				BranchName: item.BranchName,
				Address:    item.Address,
				Phone:      item.Phone,
			})
		}
	}

	return user.DetailWithBranches{User: usr, DetailBranches: branchDetails}, nil
}

func (r Repositories) FindBranchByID(ctx context.Context, id string) (branch.Branch, error) {
	var m BranchModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return branch.Branch{}, err
	}
	return toDomainBranch(m), nil
}

func (r Repositories) ListBranches(ctx context.Context, req branch.ListRequest) (branch.ListResult, error) {
	query := r.DB.WithContext(ctx).Model(&BranchModel{})
	if req.Search != "" {
		like := "%" + strings.ToLower(strings.TrimSpace(req.Search)) + "%"
		query = query.Where("LOWER(branch_name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(email) LIKE ? OR LOWER(sia_name) LIKE ?", like, like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return branch.ListResult{}, err
	}

	var models []BranchModel
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("branch_name ASC").Offset(offset).Limit(req.Limit).Find(&models).Error; err != nil {
		return branch.ListResult{}, err
	}

	items := make([]branch.Branch, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainBranch(model))
	}

	lastPage := 0
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
	}
	if lastPage == 0 {
		lastPage = 1
	}

	return branch.ListResult{
		Items: items,
		Meta: branch.ListMeta{
			Page:      req.Page,
			Limit:     req.Limit,
			TotalData: int(total),
			LastPage:  lastPage,
		},
	}, nil
}

func (r Repositories) UserHasBranch(ctx context.Context, userID, branchID string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&UserBranchModel{}).Where("user_id = ? AND branch_id = ?", userID, branchID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r Repositories) ListUserBranches(ctx context.Context, userID string) ([]auth.UserBranch, error) {
	var items []auth.UserBranch
	err := r.DB.WithContext(ctx).
		Table("user_branches").
		Select("user_branches.user_id, users.name AS user_name, user_branches.branch_id, branches.branch_name, branches.sia_name, branches.sipa_name, branches.phone").
		Joins("LEFT JOIN users ON users.id = user_branches.user_id").
		Joins("LEFT JOIN branches ON branches.id = user_branches.branch_id").
		Where("branches.branch_status = 'active' AND user_branches.user_id = ?", userID).
		Scan(&items).Error
	return items, err
}

func (r Repositories) FindProfile(ctx context.Context, userID, branchID string) (auth.Profile, error) {
	var item auth.Profile
	err := r.DB.WithContext(ctx).
		Table("user_branches usrbrc").
		Select("usrbrc.user_id AS user_id, usr.name AS profile_name, usrbrc.branch_id AS branch_id, brc.branch_name AS branch_name, brc.address, brc.phone, brc.email, brc.sia_id, brc.sia_name, brc.psa_id, brc.psa_name, brc.sipa, brc.sipa_name, brc.aping_id, brc.aping_name, brc.bank_name, brc.account_name, brc.account_number, brc.tax_percentage, brc.journal_method, brc.branch_status, brc.license_date, brc.default_member AS default_member, mbr.name AS member_name").
		Joins("LEFT JOIN users usr ON usr.id = usrbrc.user_id").
		Joins("LEFT JOIN branches brc ON brc.id = usrbrc.branch_id").
		Joins("LEFT JOIN members mbr ON mbr.id = brc.default_member").
		Where("usrbrc.branch_id = ? AND usrbrc.user_id = ?", branchID, userID).
		Scan(&item).Error
	return item, err
}

func (r Repositories) ListAllUserBranches(ctx context.Context) ([]userbranch.Detail, error) {
	var items []userbranch.Detail
	err := r.DB.WithContext(ctx).
		Table("user_branches usrb").
		Select("usrb.user_id, usr.name AS user_name, usrb.branch_id, brc.branch_name AS branch_name, brc.sia_name, brc.sipa_name, brc.phone").
		Joins("LEFT JOIN users usr ON usr.id = usrb.user_id").
		Joins("LEFT JOIN branches brc ON brc.id = usrb.branch_id").
		Scan(&items).Error
	return items, err
}

func (r Repositories) FindUserBranchDetail(ctx context.Context, userID, branchID string) ([]userbranch.Detail, error) {
	var items []userbranch.Detail
	err := r.DB.WithContext(ctx).
		Table("user_branches").
		Select("user_branches.user_id, users.name AS user_name, user_branches.branch_id, branches.branch_name, branches.sia_name, branches.sipa_name, branches.phone").
		Joins("LEFT JOIN users ON users.id = user_branches.user_id").
		Joins("LEFT JOIN branches ON branches.id = user_branches.branch_id").
		Where("branches.branch_status = 'active' AND user_branches.branch_id = ? AND user_branches.user_id = ?", branchID, userID).
		Scan(&items).Error
	return items, err
}

func (r Repositories) Create(ctx context.Context, item product.Product) error {
	return r.DB.WithContext(ctx).Create(&ProductModel{ID: item.ID, SKU: item.SKU, Name: item.Name, Description: item.Description, BranchID: item.BranchID, UnitID: item.UnitID, Stock: item.Stock, PurchasePrice: item.PurchasePrice, SalesPrice: item.SalesPrice, AlternatePrice: item.AlternatePrice, ProductCategoryID: item.ProductCategoryID, ExpiredDate: item.ExpiredDate}).Error
}

func (r Repositories) FindProductByID(ctx context.Context, id string) (product.Product, error) {
	var m ProductModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return product.Product{}, err
	}
	return toDomainProduct(m), nil
}

func (r Repositories) Update(ctx context.Context, item product.Product) error {
	return r.UpdateProduct(ctx, item)
}

func (r Repositories) UpdateProduct(ctx context.Context, item product.Product) error {
	return r.DB.WithContext(ctx).Model(&ProductModel{}).Where("id = ?", item.ID).Updates(map[string]any{
		"sku": item.SKU, "name": item.Name, "description": item.Description, "branch_id": item.BranchID, "unit_id": item.UnitID,
		"stock": item.Stock, "purchase_price": item.PurchasePrice, "sales_price": item.SalesPrice, "alternate_price": item.AlternatePrice,
		"product_category_id": item.ProductCategoryID, "expired_date": item.ExpiredDate,
	}).Error
}

func (r Repositories) GetSaleCombo(ctx context.Context, branchID, search string) ([]product.SaleComboItem, error) {
	search = strings.TrimSpace(strings.ToLower(search))
	var items []product.SaleComboItem
	query := r.DB.WithContext(ctx).Table("products").Select("products.id as product_id, products.name as product_name, sales_price AS price, products.stock, products.unit_id, units.name AS unit_name").Joins("LEFT JOIN units ON units.id = products.unit_id").Where("products.branch_id = ?", branchID)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ? OR products.id ILIKE ?", like, like, like)
	}
	return items, query.Order("products.name ASC").Scan(&items).Error
}

func (r Repositories) GetPurchaseCombo(ctx context.Context, branchID, search string) ([]product.PurchaseComboItem, error) {
	search = strings.TrimSpace(strings.ToLower(search))
	var items []product.PurchaseComboItem
	query := r.DB.WithContext(ctx).Table("products").Select("products.id as product_id, products.name as product_name, purchase_price AS price, products.unit_id, units.name AS unit_name").Joins("LEFT JOIN units ON units.id = products.unit_id").Where("products.branch_id = ?", branchID)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ? OR products.id ILIKE ?", like, like, like)
	}
	return items, query.Order("products.name ASC").Scan(&items).Error
}

func (r Repositories) GetOpnameCombo(ctx context.Context, branchID, search string) ([]product.OpnameComboItem, error) {
	search = strings.TrimSpace(strings.ToLower(search))
	var items []product.OpnameComboItem
	query := r.DB.WithContext(ctx).Table("products pro").Select("pro.id AS pro_id, pro.name AS pro_name, pro.unit_id, pro.stock, unt.name AS unit_name, pro.purchase_price AS price").Joins("LEFT JOIN units unt ON unt.id = pro.unit_id").Where("pro.branch_id = ?", branchID)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("LOWER(pro.name) LIKE ? OR LOWER(pro.id) LIKE ?", like, like)
	}
	return items, query.Order("pro.name ASC").Scan(&items).Error
}

func (r Repositories) FindUnitByID(ctx context.Context, id string) (unit.Unit, error) {
	var m UnitModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return unit.Unit{}, err
	}
	return unit.Unit{ID: m.ID, Name: m.Name}, nil
}

func (r Repositories) FindConversion(ctx context.Context, productID, initID, finalID, branchID string) (unit.Conversion, error) {
	var m UnitConversionModel
	if err := r.DB.WithContext(ctx).Where("product_id = ? AND init_id = ? AND final_id = ? AND branch_id = ?", productID, initID, finalID, branchID).First(&m).Error; err != nil {
		return unit.Conversion{}, err
	}
	return unit.Conversion{ProductID: m.ProductID, InitID: m.InitID, FinalID: m.FinalID, Value: m.ValueConv, BranchID: m.BranchID}, nil
}

func (r Repositories) FindMemberByID(ctx context.Context, id string) (member.Member, error) {
	var m MemberModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return member.Member{}, err
	}
	return member.Member{ID: m.ID, MemberCategoryID: m.MemberCategoryID, Points: m.Points}, nil
}

func (r Repositories) FindCategoryByID(ctx context.Context, id string) (member.MemberCategory, error) {
	var m MemberCategoryModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return member.MemberCategory{}, err
	}
	return member.MemberCategory{ID: m.ID, PointsConversionRate: m.PointsConversionRate}, nil
}

func (r Repositories) UpdatePoints(ctx context.Context, memberID string, points int) error {
	return r.DB.WithContext(ctx).Model(&MemberModel{}).Where("id = ?", memberID).Update("points", points).Error
}

func (r Repositories) CreateOpname(ctx context.Context, item opname.Opname) error {
	return r.DB.WithContext(ctx).Create(&OpnameModel{ID: item.ID, Description: item.Description, BranchID: item.BranchID, UserID: item.UserID, OpnameDate: item.OpnameDate, TotalOpname: item.TotalOpname, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}).Error
}

func (r Repositories) CreateOpnameItem(ctx context.Context, item opname.Item) error {
	return r.DB.WithContext(ctx).Create(&OpnameItemModel{ID: item.ID, OpnameID: item.OpnameID, ProductID: item.ProductID, Qty: item.Qty, QtyExist: item.QtyExist, Price: item.Price, SubTotal: item.SubTotal, SubTotalExist: item.SubTotalExist, ExpiredDate: item.ExpiredDate, CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
}

func (r Repositories) FindOpnameByID(ctx context.Context, id string) (opname.Opname, error) {
	var m OpnameModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return opname.Opname{}, err
	}
	return opname.Opname{ID: m.ID, Description: m.Description, BranchID: m.BranchID, UserID: m.UserID, OpnameDate: m.OpnameDate, TotalOpname: m.TotalOpname, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (r Repositories) FindOpnameItems(ctx context.Context, opnameID string) ([]opname.Item, error) {
	var items []opname.Item
	err := r.DB.WithContext(ctx).Table("opname_items oi").
		Select("oi.id, oi.opname_id, oi.product_id, p.name AS product_name, oi.qty, oi.qty_exist, oi.price, oi.sub_total, oi.sub_total_exist, oi.expired_date").
		Joins("LEFT JOIN products p ON p.id = oi.product_id").
		Where("oi.opname_id = ?", opnameID).
		Order("p.name ASC").
		Scan(&items).Error
	return items, err
}

func (r Repositories) UpdateOpnameTotal(ctx context.Context, opnameID string, total int) error {
	return r.DB.WithContext(ctx).Model(&OpnameModel{}).Where("id = ?", opnameID).Update("total_opname", total).Error
}

func (r Repositories) RecalculateOpnameTotal(ctx context.Context, opnameID string) (int, error) {
	var total int
	err := r.DB.WithContext(ctx).Table("opname_items").Select("COALESCE(SUM(sub_total - sub_total_exist), 0)").Where("opname_id = ?", opnameID).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if err := r.UpdateOpnameTotal(ctx, opnameID, total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r Repositories) WithinOpnameTransaction(ctx context.Context, fn func(repo ports.OpnameTxRepository) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return fn(txRepo{tx: tx}) })
}

func (r Repositories) WithinTransaction(ctx context.Context, fn func(repo ports.PurchaseTxRepository) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return fn(txRepo{tx: tx}) })
}

func (r Repositories) WithinTransactionSale(ctx context.Context, fn func(repo ports.SaleTxRepository) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return fn(txRepo{tx: tx}) })
}

type txRepo struct{ tx *gorm.DB }

func (t txRepo) FindProduct(ctx context.Context, id string) (product.Product, error) {
	var m ProductModel
	if err := t.tx.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return product.Product{}, err
	}
	return toDomainProduct(m), nil
}

func (t txRepo) FindProductByID(ctx context.Context, id string) (product.Product, error) {
	return t.FindProduct(ctx, id)
}

func (t txRepo) FindOpnameByID(ctx context.Context, id string) (opname.Opname, error) {
	var m OpnameModel
	if err := t.tx.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return opname.Opname{}, err
	}
	return opname.Opname{ID: m.ID, Description: m.Description, BranchID: m.BranchID, UserID: m.UserID, OpnameDate: m.OpnameDate, TotalOpname: m.TotalOpname, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (t txRepo) FindOpnameItemByOpnameAndProduct(ctx context.Context, opnameID, productID string) (opname.Item, error) {
	var m OpnameItemModel
	if err := t.tx.WithContext(ctx).Where("opname_id = ? AND product_id = ?", opnameID, productID).First(&m).Error; err != nil {
		return opname.Item{}, err
	}
	return opname.Item{ID: m.ID, OpnameID: m.OpnameID, ProductID: m.ProductID, Qty: m.Qty, QtyExist: m.QtyExist, Price: m.Price, SubTotal: m.SubTotal, SubTotalExist: m.SubTotalExist, ExpiredDate: m.ExpiredDate}, nil
}
func (t txRepo) FindUnit(ctx context.Context, id string) (unit.Unit, error) {
	var m UnitModel
	if err := t.tx.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return unit.Unit{}, err
	}
	return unit.Unit{ID: m.ID, Name: m.Name}, nil
}
func (t txRepo) FindConversion(ctx context.Context, productID, initID, finalID, branchID string) (unit.Conversion, error) {
	var m UnitConversionModel
	if err := t.tx.WithContext(ctx).Where("product_id = ? AND init_id = ? AND final_id = ? AND branch_id = ?", productID, initID, finalID, branchID).First(&m).Error; err != nil {
		return unit.Conversion{}, err
	}
	return unit.Conversion{ProductID: m.ProductID, InitID: m.InitID, FinalID: m.FinalID, Value: m.ValueConv, BranchID: m.BranchID}, nil
}
func (t txRepo) CreatePurchase(ctx context.Context, item purchase.Purchase) error {
	return t.tx.WithContext(ctx).Create(&PurchaseModel{ID: item.ID, SupplierID: item.SupplierID, PurchaseDate: item.PurchaseDate, BranchID: item.BranchID, UserID: item.UserID, Payment: string(item.Payment), TotalPurchase: item.TotalPurchase, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}).Error
}
func (t txRepo) CreatePurchaseItems(ctx context.Context, items []purchase.Item) error {
	models := make([]PurchaseItemModel, 0, len(items))
	for _, item := range items {
		models = append(models, PurchaseItemModel{ID: item.ID, PurchaseID: item.PurchaseID, ProductID: item.ProductID, UnitID: item.UnitID, Price: item.Price, Qty: item.Qty, SubTotal: item.SubTotal, ExpiredDate: item.ExpiredDate})
	}
	return t.tx.WithContext(ctx).Create(&models).Error
}
func (t txRepo) UpdateProduct(ctx context.Context, item product.Product) error {
	return t.tx.WithContext(ctx).Model(&ProductModel{}).Where("id = ?", item.ID).Updates(map[string]any{"stock": item.Stock, "expired_date": item.ExpiredDate, "purchase_price": item.PurchasePrice, "sales_price": item.SalesPrice}).Error
}

func (t txRepo) CreateOpnameItem(ctx context.Context, item opname.Item) error {
	return t.tx.WithContext(ctx).Create(&OpnameItemModel{ID: item.ID, OpnameID: item.OpnameID, ProductID: item.ProductID, Qty: item.Qty, QtyExist: item.QtyExist, Price: item.Price, SubTotal: item.SubTotal, SubTotalExist: item.SubTotalExist, ExpiredDate: item.ExpiredDate, CreatedAt: time.Now(), UpdatedAt: time.Now()}).Error
}

func (t txRepo) UpdateOpnameItem(ctx context.Context, item opname.Item) error {
	return t.tx.WithContext(ctx).Model(&OpnameItemModel{}).Where("id = ?", item.ID).Updates(map[string]any{"qty": item.Qty, "qty_exist": item.QtyExist, "price": item.Price, "sub_total": item.SubTotal, "sub_total_exist": item.SubTotalExist, "expired_date": item.ExpiredDate, "updated_at": time.Now()}).Error
}

func (t txRepo) UpdateOpnameTotal(ctx context.Context, opnameID string, total int) error {
	return t.tx.WithContext(ctx).Model(&OpnameModel{}).Where("id = ?", opnameID).Update("total_opname", total).Error
}

func (t txRepo) RecalculateOpnameTotal(ctx context.Context, opnameID string) (int, error) {
	var total int
	err := t.tx.WithContext(ctx).Table("opname_items").Select("COALESCE(SUM(sub_total - sub_total_exist), 0)").Where("opname_id = ?", opnameID).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if err := t.UpdateOpnameTotal(ctx, opnameID, total); err != nil {
		return 0, err
	}
	return total, nil
}
func (t txRepo) CreateTransactionReport(ctx context.Context, id string, txType string, userID string, branchID string, total int, payment string, createdAt time.Time) error {
	return t.tx.WithContext(ctx).Create(&TransactionReportModel{ID: id, TransactionType: txType, UserID: userID, BranchID: branchID, Total: total, Payment: payment, CreatedAt: createdAt, UpdatedAt: createdAt}).Error
}
func (t txRepo) CreateSale(ctx context.Context, item sale.Sale) error {
	return t.tx.WithContext(ctx).Create(&SaleModel{ID: item.ID, MemberID: item.MemberID, UserID: item.UserID, BranchID: item.BranchID, Payment: string(item.Payment), Discount: item.Discount, TotalSale: item.TotalSale, ProfitEstimate: item.ProfitEstimate, SaleDate: item.SaleDate, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}).Error
}
func (t txRepo) CreateSaleItems(ctx context.Context, items []sale.Item) error {
	models := make([]SaleItemModel, 0, len(items))
	for _, item := range items {
		models = append(models, SaleItemModel{ID: item.ID, SaleID: item.SaleID, ProductID: item.ProductID, Price: item.Price, Qty: item.Qty, SubTotal: item.SubTotal})
	}
	return t.tx.WithContext(ctx).Create(&models).Error
}
func (t txRepo) UpsertDailyProfit(ctx context.Context, reportDate time.Time, userID string, branchID string, totalSales int, profitEstimate int, now time.Time) error {
	var report DailyProfitReportModel
	err := t.tx.WithContext(ctx).Where("report_date = ? AND branch_id = ? AND user_id = ?", reportDate.Format("2006-01-02"), branchID, userID).First(&report).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return t.tx.WithContext(ctx).Create(&DailyProfitReportModel{ID: fmt.Sprintf("DPR%s", now.Format("060102150405")), ReportDate: reportDate, UserID: userID, BranchID: branchID, TotalSales: totalSales, ProfitEstimate: profitEstimate, CreatedAt: now, UpdatedAt: now}).Error
	}
	if err != nil {
		return err
	}
	report.TotalSales += totalSales
	report.ProfitEstimate += profitEstimate
	report.UpdatedAt = now
	return t.tx.WithContext(ctx).Save(&report).Error
}
func (t txRepo) FindBranch(ctx context.Context, branchID string) (branch.Branch, error) {
	var m BranchModel
	if err := t.tx.WithContext(ctx).Where("id = ?", branchID).First(&m).Error; err != nil {
		return branch.Branch{}, err
	}
	return toDomainBranch(m), nil
}
func (t txRepo) UpdateBranchQuota(ctx context.Context, branchID string, quota int) error {
	return t.tx.WithContext(ctx).Model(&BranchModel{}).Where("id = ?", branchID).Update("quota", quota).Error
}
func (t txRepo) FindMember(ctx context.Context, memberID string) (member.Member, error) {
	var m MemberModel
	if err := t.tx.WithContext(ctx).Where("id = ?", memberID).First(&m).Error; err != nil {
		return member.Member{}, err
	}
	return member.Member{ID: m.ID, MemberCategoryID: m.MemberCategoryID, Points: m.Points}, nil
}
func (t txRepo) FindMemberCategory(ctx context.Context, categoryID string) (member.MemberCategory, error) {
	var m MemberCategoryModel
	if err := t.tx.WithContext(ctx).Where("id = ?", categoryID).First(&m).Error; err != nil {
		return member.MemberCategory{}, err
	}
	return member.MemberCategory{ID: m.ID, PointsConversionRate: m.PointsConversionRate}, nil
}
func (t txRepo) UpdateMemberPoints(ctx context.Context, memberID string, points int) error {
	return t.tx.WithContext(ctx).Model(&MemberModel{}).Where("id = ?", memberID).Update("points", points).Error
}

func toDomainBranch(m BranchModel) branch.Branch {
	return branch.Branch{
		ID:               m.ID,
		BranchName:       m.BranchName,
		Address:          m.Address,
		Phone:            m.Phone,
		Email:            m.Email,
		SIAID:            m.SIAID,
		SIAName:          m.SIAName,
		PSAID:            m.PSAID,
		PSAName:          m.PSAName,
		SIPA:             m.SIPA,
		SIPAName:         m.SIPAName,
		APINGID:          m.APINGID,
		APINGName:        m.APINGName,
		BankName:         m.BankName,
		AccountName:      m.AccountName,
		AccountNumber:    m.AccountNumber,
		TaxPercentage:    m.TaxPercentage,
		JournalMethod:    m.JournalMethod,
		BranchStatus:     m.BranchStatus,
		LicenseDate:      m.LicenseDate.Format(time.RFC3339),
		DefaultMemberID:  m.DefaultMember,
		SubscriptionType: m.SubscriptionType,
		Quota:            m.Quota,
		RealAsset:        m.RealAsset,
	}
}

func toDomainProduct(m ProductModel) product.Product {
	return product.Product{ID: m.ID, SKU: m.SKU, Name: m.Name, Description: m.Description, BranchID: m.BranchID, UnitID: m.UnitID, Stock: m.Stock, PurchasePrice: m.PurchasePrice, SalesPrice: m.SalesPrice, AlternatePrice: m.AlternatePrice, ProductCategoryID: m.ProductCategoryID, ExpiredDate: m.ExpiredDate}
}
