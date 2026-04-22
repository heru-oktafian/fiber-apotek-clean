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

func (r Repositories) CreateUser(ctx context.Context, item user.User) error {
	return r.DB.WithContext(ctx).Create(&UserModel{
		ID:         item.ID,
		Name:       item.Name,
		Username:   item.Username,
		Password:   item.Password,
		UserRole:   string(item.Role),
		UserStatus: item.Status,
	}).Error
}

func (r Repositories) UpdateUser(ctx context.Context, item user.User) error {
	updates := map[string]any{
		"name":        item.Name,
		"username":    item.Username,
		"user_role":   string(item.Role),
		"user_status": item.Status,
	}
	if item.Password != "" {
		updates["password"] = item.Password
	}
	return r.DB.WithContext(ctx).Model(&UserModel{}).Where("id = ?", item.ID).Updates(updates).Error
}

func (r Repositories) CreateUserBranch(ctx context.Context, userID, branchID string) error {
	return r.DB.WithContext(ctx).Create(&UserBranchModel{UserID: userID, BranchID: branchID}).Error
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

func (r Repositories) CreateBranch(ctx context.Context, item branch.Branch) error {
	var licenseDate time.Time
	if strings.TrimSpace(item.LicenseDate) != "" {
		parsed, err := time.Parse("2006-01-02", item.LicenseDate)
		if err != nil {
			parsed, err = time.Parse(time.RFC3339, item.LicenseDate)
			if err != nil {
				return err
			}
		}
		licenseDate = parsed
	}

	return r.DB.WithContext(ctx).Create(&BranchModel{
		ID:               item.ID,
		BranchName:       item.BranchName,
		Address:          item.Address,
		Phone:            item.Phone,
		Email:            item.Email,
		SIAID:            item.SIAID,
		SIAName:          item.SIAName,
		PSAID:            item.PSAID,
		PSAName:          item.PSAName,
		SIPA:             item.SIPA,
		SIPAName:         item.SIPAName,
		APINGID:          item.APINGID,
		APINGName:        item.APINGName,
		BankName:         item.BankName,
		AccountName:      item.AccountName,
		AccountNumber:    item.AccountNumber,
		TaxPercentage:    item.TaxPercentage,
		JournalMethod:    item.JournalMethod,
		BranchStatus:     item.BranchStatus,
		LicenseDate:      licenseDate,
		DefaultMember:    item.DefaultMemberID,
		Quota:            item.Quota,
		SubscriptionType: item.SubscriptionType,
		RealAsset:        item.RealAsset,
	}).Error
}

func (r Repositories) DeleteBranch(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&BranchModel{}, "id = ?", id).Error
}

func (r Repositories) BranchHasUsers(ctx context.Context, branchID string) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&UserBranchModel{}).Where("branch_id = ?", branchID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
	return r.DB.WithContext(ctx).Create(&ProductModel{ID: item.ID, SKU: item.SKU, Name: item.Name, Alias: item.Alias, Description: item.Description, Ingredient: item.Ingredient, Dosage: item.Dosage, SideAffection: item.SideAffection, BranchID: item.BranchID, UnitID: item.UnitID, Stock: item.Stock, PurchasePrice: item.PurchasePrice, SalesPrice: item.SalesPrice, AlternatePrice: item.AlternatePrice, ProductCategoryID: item.ProductCategoryID, ExpiredDate: item.ExpiredDate}).Error
}

func (r Repositories) FindProductByID(ctx context.Context, id string) (product.Product, error) {
	var m ProductModel
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return product.Product{}, err
	}
	return toDomainProduct(m), nil
}

func (r Repositories) FindProductDetailByID(ctx context.Context, id, branchID string) (product.Product, error) {
	var row struct {
		ID                  string
		SKU                 string
		Name                string
		Alias               string
		Description         string
		Ingredient          string
		Dosage              string
		SideAffection       string
		UnitID              string
		UnitName            string
		Stock               int
		PurchasePrice       int
		ExpiredDate         time.Time
		SalesPrice          int
		AlternatePrice      int
		ProductCategoryID   uint
		ProductCategoryName string
		BranchID            string
	}
	if err := r.DB.WithContext(ctx).
		Table("products pro").
		Select("pro.id, pro.sku, pro.name, pro.alias, pro.description, pro.ingredient, pro.dosage, pro.side_affection, pro.unit_id, un.name AS unit_name, pro.stock, pro.purchase_price, pro.expired_date, pro.sales_price, pro.alternate_price, pro.product_category_id, pc.name AS product_category_name, pro.branch_id").
		Joins("LEFT JOIN product_categories pc ON pc.id = pro.product_category_id").
		Joins("LEFT JOIN units un ON un.id = pro.unit_id").
		Where("pro.id = ? AND pro.branch_id = ?", id, branchID).
		Scan(&row).Error; err != nil {
		return product.Product{}, err
	}
	if row.ID == "" {
		return product.Product{}, gorm.ErrRecordNotFound
	}
	return product.Product{ID: row.ID, SKU: row.SKU, Name: row.Name, Alias: row.Alias, Description: row.Description, Ingredient: row.Ingredient, Dosage: row.Dosage, SideAffection: row.SideAffection, UnitID: row.UnitID, UnitName: row.UnitName, Stock: row.Stock, PurchasePrice: row.PurchasePrice, ExpiredDate: row.ExpiredDate, SalesPrice: row.SalesPrice, AlternatePrice: row.AlternatePrice, ProductCategoryID: row.ProductCategoryID, ProductCategoryName: row.ProductCategoryName, BranchID: row.BranchID}, nil
}

func (r Repositories) ListProducts(ctx context.Context, branchID string, req product.ListRequest) (product.ListResult, error) {
	query := r.DB.WithContext(ctx).
		Table("products pro").
		Select("pro.id, pro.sku, pro.name, pro.alias, pro.description, pro.ingredient, pro.dosage, pro.side_affection, pro.unit_id, un.name AS unit_name, pro.stock, pro.purchase_price, pro.sales_price, pro.alternate_price, pro.expired_date, pro.product_category_id, pc.name AS product_category_name").
		Joins("LEFT JOIN product_categories pc ON pc.id = pro.product_category_id").
		Joins("LEFT JOIN units un ON un.id = pro.unit_id").
		Where("pro.branch_id = ?", branchID)
	if req.Search != "" {
		like := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("pro.name ILIKE ? OR pro.alias ILIKE ? OR pro.description ILIKE ? OR pro.ingredient ILIKE ? OR pro.dosage ILIKE ? OR pro.side_affection ILIKE ?", like, like, like, like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return product.ListResult{}, err
	}
	var rows []product.Product
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("pro.name ASC").Offset(offset).Limit(req.Limit).Scan(&rows).Error; err != nil {
		return product.ListResult{}, err
	}
	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}
	return product.ListResult{Items: rows, Meta: product.ListMeta{Page: req.Page, Limit: req.Limit, Search: req.Search, TotalData: int(total), LastPage: lastPage}}, nil
}

func (r Repositories) Update(ctx context.Context, item product.Product) error {
	return r.UpdateProduct(ctx, item)
}

func (r Repositories) UpdateProduct(ctx context.Context, item product.Product) error {
	return r.DB.WithContext(ctx).Model(&ProductModel{}).Where("id = ?", item.ID).Updates(map[string]any{
		"sku": item.SKU, "name": item.Name, "alias": item.Alias, "description": item.Description, "ingredient": item.Ingredient, "dosage": item.Dosage, "side_affection": item.SideAffection,
		"branch_id": item.BranchID, "unit_id": item.UnitID, "stock": item.Stock, "purchase_price": item.PurchasePrice, "sales_price": item.SalesPrice, "alternate_price": item.AlternatePrice,
		"product_category_id": item.ProductCategoryID, "expired_date": item.ExpiredDate,
	}).Error
}

func (r Repositories) DeleteProduct(ctx context.Context, id, branchID string) error {
	return r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).Delete(&ProductModel{}).Error
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

func (r Repositories) ListSuppliers(ctx context.Context, branchID string, req supplier.ListRequest) (supplier.ListResult, error) {
	query := r.DB.WithContext(ctx).
		Table("suppliers s").
		Select("s.id, s.name, s.phone, s.address, s.pic, s.supplier_category_id, sc.name AS supplier_category").
		Joins("LEFT JOIN supplier_categories sc ON sc.id = s.supplier_category_id").
		Where("s.branch_id = ?", branchID)

	if req.Search != "" {
		like := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("s.name ILIKE ? OR s.address ILIKE ? OR sc.name ILIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return supplier.ListResult{}, err
	}

	var rows []struct {
		ID               string
		Name             string
		Phone            string
		Address          string
		PIC              string
		SupplierCategory string
		SupplierCategoryID uint
	}
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("s.name ASC").Offset(offset).Limit(req.Limit).Scan(&rows).Error; err != nil {
		return supplier.ListResult{}, err
	}

	items := make([]supplier.Supplier, 0, len(rows))
	for _, row := range rows {
		items = append(items, supplier.Supplier{ID: row.ID, Name: row.Name, Phone: row.Phone, Address: row.Address, PIC: row.PIC, SupplierCategoryID: row.SupplierCategoryID, SupplierCategory: row.SupplierCategory})
	}

	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}

	return supplier.ListResult{Items: items, Meta: supplier.ListMeta{Page: req.Page, Limit: req.Limit, Search: req.Search, TotalData: int(total), LastPage: lastPage}}, nil
}

func (r Repositories) FindSupplierByID(ctx context.Context, id, branchID string) (supplier.Supplier, error) {
	var row struct {
		ID                 string
		Name               string
		Phone              string
		Address            string
		PIC                string
		SupplierCategoryID uint
		SupplierCategory   string
		BranchID           string
	}
	if err := r.DB.WithContext(ctx).
		Table("suppliers s").
		Select("s.id, s.name, s.phone, s.address, s.pic, s.supplier_category_id, s.branch_id, sc.name AS supplier_category").
		Joins("LEFT JOIN supplier_categories sc ON sc.id = s.supplier_category_id").
		Where("s.id = ? AND s.branch_id = ?", id, branchID).
		Scan(&row).Error; err != nil {
		return supplier.Supplier{}, err
	}
	if row.ID == "" {
		return supplier.Supplier{}, gorm.ErrRecordNotFound
	}
	return supplier.Supplier{ID: row.ID, Name: row.Name, Phone: row.Phone, Address: row.Address, PIC: row.PIC, SupplierCategoryID: row.SupplierCategoryID, SupplierCategory: row.SupplierCategory, BranchID: row.BranchID}, nil
}

func (r Repositories) CreateSupplier(ctx context.Context, item supplier.Supplier) error {
	return r.DB.WithContext(ctx).Create(&SupplierModel{ID: item.ID, Name: item.Name, Phone: item.Phone, Address: item.Address, PIC: item.PIC, SupplierCategoryID: item.SupplierCategoryID, BranchID: item.BranchID}).Error
}

func (r Repositories) UpdateSupplier(ctx context.Context, item supplier.Supplier) error {
	return r.DB.WithContext(ctx).Model(&SupplierModel{}).Where("id = ? AND branch_id = ?", item.ID, item.BranchID).Updates(map[string]any{
		"name": item.Name,
		"phone": item.Phone,
		"address": item.Address,
		"pic": item.PIC,
		"supplier_category_id": item.SupplierCategoryID,
	}).Error
}

func (r Repositories) DeleteSupplier(ctx context.Context, id, branchID string) error {
	return r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).Delete(&SupplierModel{}).Error
}

func (r Repositories) GetSupplierCombo(ctx context.Context, branchID, search string) ([]supplier.ComboItem, error) {
	search = strings.TrimSpace(strings.ToLower(search))
	var items []supplier.ComboItem
	query := r.DB.WithContext(ctx).Table("suppliers").Select("id AS supplier_id, name AS supplier_name").Where("branch_id = ?", branchID)
	if search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+search+"%")
	}
	return items, query.Order("name ASC").Scan(&items).Error
}

func (r Repositories) ListMasterUnits(ctx context.Context, branchID string, req unit.MasterUnitListRequest) (unit.MasterUnitListResult, error) {
	query := r.DB.WithContext(ctx).Table("units un").Select("un.id, un.name, un.branch_id").Where("un.branch_id = ?", branchID)
	if req.Search != "" {
		like := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("un.name ILIKE ?", like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return unit.MasterUnitListResult{}, err
	}
	var rows []UnitModel
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("un.name ASC").Offset(offset).Limit(req.Limit).Scan(&rows).Error; err != nil {
		return unit.MasterUnitListResult{}, err
	}
	items := make([]unit.MasterUnit, 0, len(rows))
	for _, row := range rows {
		items = append(items, unit.MasterUnit{ID: row.ID, Name: row.Name, BranchID: row.BranchID})
	}
	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}
	return unit.MasterUnitListResult{Items: items, Meta: unit.MasterUnitListMeta{Page: req.Page, Limit: req.Limit, Search: req.Search, TotalData: int(total), LastPage: lastPage}}, nil
}

func (r Repositories) FindMasterUnitByID(ctx context.Context, id, branchID string) (unit.MasterUnit, error) {
	var m UnitModel
	if err := r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).First(&m).Error; err != nil {
		return unit.MasterUnit{}, err
	}
	return unit.MasterUnit{ID: m.ID, Name: m.Name, BranchID: m.BranchID}, nil
}

func (r Repositories) CreateMasterUnit(ctx context.Context, item unit.MasterUnit) error {
	return r.DB.WithContext(ctx).Create(&UnitModel{ID: item.ID, Name: item.Name, BranchID: item.BranchID}).Error
}

func (r Repositories) UpdateMasterUnit(ctx context.Context, item unit.MasterUnit) error {
	return r.DB.WithContext(ctx).Model(&UnitModel{}).Where("id = ? AND branch_id = ?", item.ID, item.BranchID).Update("name", item.Name).Error
}

func (r Repositories) DeleteMasterUnit(ctx context.Context, id, branchID string) error {
	return r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).Delete(&UnitModel{}).Error
}

func (r Repositories) GetMasterUnitCombo(ctx context.Context, branchID, search string) ([]unit.MasterUnitComboItem, error) {
	search = strings.TrimSpace(strings.ToLower(search))
	var items []unit.MasterUnitComboItem
	query := r.DB.WithContext(ctx).Table("units").Select("id as unit_id, name as unit_name").Where("branch_id = ?", branchID)
	if search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+search+"%")
	}
	return items, query.Order("name ASC").Scan(&items).Error
}

func (r Repositories) ListProductCategories(ctx context.Context, branchID string, req productcategory.ListRequest) (productcategory.ListResult, error) {
	query := r.DB.WithContext(ctx).Table("product_categories pc").Select("pc.id AS product_category_id, pc.name AS product_category_name").Where("pc.branch_id = ?", branchID)
	if req.Search != "" {
		like := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("pc.name ILIKE ?", like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return productcategory.ListResult{}, err
	}
	var items []productcategory.ComboItem
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("pc.name ASC").Offset(offset).Limit(req.Limit).Scan(&items).Error; err != nil {
		return productcategory.ListResult{}, err
	}
	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}
	return productcategory.ListResult{Items: items, Meta: productcategory.ListMeta{Page: req.Page, Limit: req.Limit, Search: req.Search, TotalData: int(total), LastPage: lastPage}}, nil
}

func (r Repositories) FindProductCategoryByID(ctx context.Context, id uint, branchID string) (productcategory.ProductCategory, error) {
	var m ProductCategoryModel
	if err := r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).First(&m).Error; err != nil {
		return productcategory.ProductCategory{}, err
	}
	return productcategory.ProductCategory{ID: m.ID, Name: m.Name, BranchID: m.BranchID}, nil
}

func (r Repositories) CreateProductCategory(ctx context.Context, item productcategory.ProductCategory) (productcategory.ProductCategory, error) {
	m := ProductCategoryModel{Name: item.Name, BranchID: item.BranchID}
	if err := r.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return productcategory.ProductCategory{}, err
	}
	return productcategory.ProductCategory{ID: m.ID, Name: m.Name, BranchID: m.BranchID}, nil
}

func (r Repositories) UpdateProductCategory(ctx context.Context, item productcategory.ProductCategory) error {
	return r.DB.WithContext(ctx).Model(&ProductCategoryModel{}).Where("id = ? AND branch_id = ?", item.ID, item.BranchID).Update("name", item.Name).Error
}

func (r Repositories) DeleteProductCategory(ctx context.Context, id uint, branchID string) error {
	return r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).Delete(&ProductCategoryModel{}).Error
}

func (r Repositories) GetProductCategoryCombo(ctx context.Context, branchID, search string) ([]productcategory.ComboItem, error) {
	search = strings.TrimSpace(strings.ToLower(search))
	var items []productcategory.ComboItem
	query := r.DB.WithContext(ctx).Table("product_categories").Select("id as product_category_id, name as product_category_name").Where("branch_id = ?", branchID)
	if search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+search+"%")
	}
	return items, query.Order("name ASC").Scan(&items).Error
}

func (r Repositories) ListSupplierCategories(ctx context.Context, branchID string, req suppliercategory.ListRequest) (suppliercategory.ListResult, error) {
	query := r.DB.WithContext(ctx).Table("supplier_categories sc").Select("sc.id, sc.name, sc.branch_id").Where("sc.branch_id = ?", branchID).Order("sc.name ASC")
	if req.Search != "" {
		like := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("sc.name ILIKE ?", like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return suppliercategory.ListResult{}, err
	}
	var rows []SupplierCategoryModel
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Scan(&rows).Error; err != nil {
		return suppliercategory.ListResult{}, err
	}
	items := make([]suppliercategory.SupplierCategory, 0, len(rows))
	for _, row := range rows {
		items = append(items, suppliercategory.SupplierCategory{ID: row.ID, Name: row.Name, BranchID: row.BranchID})
	}
	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}
	return suppliercategory.ListResult{Items: items, Meta: suppliercategory.ListMeta{Page: req.Page, Limit: req.Limit, Search: req.Search, TotalData: int(total), LastPage: lastPage}}, nil
}

func (r Repositories) FindSupplierCategoryByID(ctx context.Context, id uint, branchID string) (suppliercategory.SupplierCategory, error) {
	var m SupplierCategoryModel
	if err := r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).First(&m).Error; err != nil {
		return suppliercategory.SupplierCategory{}, err
	}
	return suppliercategory.SupplierCategory{ID: m.ID, Name: m.Name, BranchID: m.BranchID}, nil
}

func (r Repositories) CreateSupplierCategory(ctx context.Context, item suppliercategory.SupplierCategory) (suppliercategory.SupplierCategory, error) {
	m := SupplierCategoryModel{Name: item.Name, BranchID: item.BranchID}
	if err := r.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return suppliercategory.SupplierCategory{}, err
	}
	return suppliercategory.SupplierCategory{ID: m.ID, Name: m.Name, BranchID: m.BranchID}, nil
}

func (r Repositories) UpdateSupplierCategory(ctx context.Context, item suppliercategory.SupplierCategory) error {
	return r.DB.WithContext(ctx).Model(&SupplierCategoryModel{}).Where("id = ? AND branch_id = ?", item.ID, item.BranchID).Update("name", item.Name).Error
}

func (r Repositories) DeleteSupplierCategory(ctx context.Context, id uint, branchID string) error {
	return r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).Delete(&SupplierCategoryModel{}).Error
}

func (r Repositories) GetSupplierCategoryCombo(ctx context.Context, branchID string) ([]suppliercategory.ComboItem, error) {
	var items []suppliercategory.ComboItem
	query := r.DB.WithContext(ctx).Table("supplier_categories").Select("id AS supplier_category_id, name AS supplier_category_name").Where("branch_id = ?", branchID).Order("name ASC")
	return items, query.Scan(&items).Error
}

func (r Repositories) ListMemberCategories(ctx context.Context, branchID string, req membercategory.ListRequest) (membercategory.ListResult, error) {
	query := r.DB.WithContext(ctx).Table("member_categories mc").Select("mc.id, mc.name, mc.points_conversion_rate, mc.branch_id").Where("mc.branch_id = ?", branchID)
	if req.Search != "" {
		like := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("mc.name ILIKE ?", like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return membercategory.ListResult{}, err
	}
	var rows []MemberCategoryModel
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("mc.name ASC").Offset(offset).Limit(req.Limit).Scan(&rows).Error; err != nil {
		return membercategory.ListResult{}, err
	}
	items := make([]membercategory.MemberCategory, 0, len(rows))
	for _, row := range rows {
		items = append(items, membercategory.MemberCategory{ID: row.ID, Name: row.Name, PointsConversionRate: row.PointsConversionRate, BranchID: row.BranchID})
	}
	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}
	return membercategory.ListResult{Items: items, Meta: membercategory.ListMeta{Page: req.Page, Limit: req.Limit, Search: req.Search, TotalData: int(total), LastPage: lastPage}}, nil
}

func (r Repositories) FindMemberCategoryByID(ctx context.Context, id uint, branchID string) (membercategory.MemberCategory, error) {
	var m MemberCategoryModel
	if err := r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).First(&m).Error; err != nil {
		return membercategory.MemberCategory{}, err
	}
	return membercategory.MemberCategory{ID: m.ID, Name: m.Name, PointsConversionRate: m.PointsConversionRate, BranchID: m.BranchID}, nil
}

func (r Repositories) CreateMemberCategory(ctx context.Context, item membercategory.MemberCategory) (membercategory.MemberCategory, error) {
	m := MemberCategoryModel{Name: item.Name, PointsConversionRate: item.PointsConversionRate, BranchID: item.BranchID}
	if err := r.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return membercategory.MemberCategory{}, err
	}
	return membercategory.MemberCategory{ID: m.ID, Name: m.Name, PointsConversionRate: m.PointsConversionRate, BranchID: m.BranchID}, nil
}

func (r Repositories) UpdateMemberCategory(ctx context.Context, item membercategory.MemberCategory) error {
	return r.DB.WithContext(ctx).Model(&MemberCategoryModel{}).Where("id = ? AND branch_id = ?", item.ID, item.BranchID).Updates(map[string]any{"name": item.Name, "points_conversion_rate": item.PointsConversionRate}).Error
}

func (r Repositories) DeleteMemberCategory(ctx context.Context, id uint, branchID string) error {
	return r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).Delete(&MemberCategoryModel{}).Error
}

func (r Repositories) GetMemberCategoryCombo(ctx context.Context, branchID, search string) ([]membercategory.ComboItem, error) {
	search = strings.TrimSpace(strings.ToLower(search))
	var items []membercategory.ComboItem
	query := r.DB.WithContext(ctx).Table("member_categories").Select("id AS member_category_id, name AS member_category_name").Where("branch_id = ?", branchID)
	if search != "" {
		query = query.Where("LOWER(member_categories.name) ILIKE ?", "%"+search+"%")
	}
	return items, query.Order("name ASC").Scan(&items).Error
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
	return member.MemberCategory{ID: fmt.Sprintf("%d", m.ID), PointsConversionRate: m.PointsConversionRate}, nil
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

func (r Repositories) ListPurchases(ctx context.Context, branchID string, req purchase.ListRequest) (purchase.ListResult, error) {
	query := r.DB.WithContext(ctx).
		Table("purchases pur").
		Select("pur.id, pur.supplier_id, sup.name AS supplier_name, pur.purchase_date, pur.total_purchase, pur.payment").
		Joins("LEFT JOIN suppliers sup ON sup.id = pur.supplier_id").
		Where("pur.branch_id = ? AND pur.total_purchase > 0", branchID)
	if req.Search != "" {
		like := "%" + strings.TrimSpace(strings.ToLower(req.Search)) + "%"
		query = query.Where("LOWER(sup.name) LIKE ?", like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return purchase.ListResult{}, err
	}
	var items []purchase.Purchase
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("pur.created_at DESC").Offset(offset).Limit(req.Limit).Scan(&items).Error; err != nil {
		return purchase.ListResult{}, err
	}
	lastPage := 1
	if req.Limit > 0 {
		lastPage = int((total + int64(req.Limit) - 1) / int64(req.Limit))
		if lastPage == 0 {
			lastPage = 1
		}
	}
	return purchase.ListResult{Items: items, Meta: purchase.ListMeta{Page: req.Page, Limit: req.Limit, Search: req.Search, TotalData: int(total), LastPage: lastPage}}, nil
}

func (r Repositories) FindPurchaseByID(ctx context.Context, branchID, id string) (purchase.Purchase, error) {
	var m PurchaseModel
	if err := r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).First(&m).Error; err != nil {
		return purchase.Purchase{}, err
	}
	return purchase.Purchase{ID: m.ID, SupplierID: m.SupplierID, PurchaseDate: m.PurchaseDate, BranchID: m.BranchID, UserID: m.UserID, Payment: common.PaymentStatus(m.Payment), TotalPurchase: m.TotalPurchase, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (r Repositories) FindPurchaseDetail(ctx context.Context, branchID, id string) (purchase.Detail, error) {
	var header struct {
		ID            string
		SupplierID    string
		SupplierName  string
		PurchaseDate  time.Time
		TotalPurchase int
		Payment       string
	}
	if err := r.DB.WithContext(ctx).
		Table("purchases pur").
		Select("pur.id, pur.supplier_id, sup.name AS supplier_name, pur.purchase_date, pur.total_purchase, pur.payment").
		Joins("LEFT JOIN suppliers sup ON sup.id = pur.supplier_id").
		Where("pur.id = ? AND pur.branch_id = ?", id, branchID).
		Scan(&header).Error; err != nil {
		return purchase.Detail{}, err
	}
	if header.ID == "" {
		return purchase.Detail{}, gorm.ErrRecordNotFound
	}
	var items []struct {
		ID          string
		ProductID   string
		ProductName string
		UnitID      string
		UnitName    string
		Price       int
		Qty         int
		SubTotal    int
		ExpiredDate time.Time
	}
	if err := r.DB.WithContext(ctx).
		Table("purchase_items pit").
		Select("pit.id, pit.product_id, pro.name AS product_name, pit.unit_id AS unit_id, un.name AS unit_name, pit.price, pit.qty, pit.sub_total, pit.expired_date").
		Joins("LEFT JOIN products pro ON pro.id = pit.product_id").
		Joins("LEFT JOIN units un ON un.id = pit.unit_id").
		Where("pit.purchase_id = ?", id).
		Order("pro.name ASC").
		Scan(&items).Error; err != nil {
		return purchase.Detail{}, err
	}
	formattedItems := make([]purchase.FormattedItem, 0, len(items))
	for _, item := range items {
		formattedItems = append(formattedItems, purchase.FormattedItem{ID: item.ID, ProductID: item.ProductID, ProductName: item.ProductName, UnitID: item.UnitID, UnitName: item.UnitName, Price: item.Price, Qty: item.Qty, SubTotal: item.SubTotal, ExpiredDate: item.ExpiredDate.Format("02 January 2006")})
	}
	return purchase.Detail{ID: header.ID, SupplierID: header.SupplierID, SupplierName: header.SupplierName, PurchaseDate: header.PurchaseDate.Format("02 January 2006"), TotalPurchase: header.TotalPurchase, Payment: header.Payment, Items: formattedItems}, nil
}

func (r Repositories) UpdatePurchaseHeader(ctx context.Context, item purchase.Purchase) error {
	return r.DB.WithContext(ctx).Model(&PurchaseModel{}).Where("id = ? AND branch_id = ?", item.ID, item.BranchID).Updates(map[string]any{"supplier_id": item.SupplierID, "purchase_date": item.PurchaseDate, "payment": string(item.Payment), "updated_at": item.UpdatedAt}).Error
}

func (r Repositories) DeletePurchaseHeader(ctx context.Context, branchID, id string) error {
	return r.DB.WithContext(ctx).Where("id = ? AND branch_id = ?", id, branchID).Delete(&PurchaseModel{}).Error
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
	return member.MemberCategory{ID: fmt.Sprintf("%d", m.ID), PointsConversionRate: m.PointsConversionRate}, nil
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
