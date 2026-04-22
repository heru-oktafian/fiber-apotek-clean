package firststock

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/firststock"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	"gorm.io/gorm"
)

type Service struct {
	Repo  ports.FirstStockRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req firststock.ListRequest) (firststock.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Month == "" {
		req.Month = s.Clock.Now().Format("2006-01")
	}
	return s.Repo.ListFirstStocks(ctx, branchID, req)
}

func (s Service) Create(ctx context.Context, branchID, userID string, req firststock.CreateRequest) (firststock.FirstStock, error) {
	now := s.Clock.Now()
	date := now
	if strings.TrimSpace(req.FirstStockDate) != "" {
		parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.FirstStockDate))
		if err != nil {
			return firststock.FirstStock{}, apperror.New(http.StatusBadRequest, "Create first stock failed", "invalid date format, use YYYY-MM-DD")
		}
		date = parsedDate
	}
	item := firststock.FirstStock{ID: s.IDs.New("FST"), Description: strings.TrimSpace(req.Description), FirstStockDate: date, BranchID: branchID, UserID: userID, TotalFirstStock: 0, Payment: common.PaymentStatus("nocost"), CreatedAt: now, UpdatedAt: now}
	if err := s.Repo.CreateFirstStock(ctx, item); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Create first stock failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "first_stock", item.UserID, item.BranchID, item.TotalFirstStock, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Create first stock failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id string, req firststock.UpdateRequest) (firststock.FirstStock, error) {
	item, err := s.Repo.FindFirstStockByID(ctx, branchID, id)
	if err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusNotFound, "Update first stock failed", "first stock not found")
	}
	if strings.TrimSpace(req.FirstStockDate) != "" {
		parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.FirstStockDate))
		if err != nil {
			return firststock.FirstStock{}, apperror.New(http.StatusBadRequest, "Update first stock failed", "invalid date format, use YYYY-MM-DD")
		}
		item.FirstStockDate = parsedDate
	}
	if strings.TrimSpace(req.Description) != "" {
		item.Description = strings.TrimSpace(req.Description)
	}
	if strings.TrimSpace(req.Payment) != "" {
		item.Payment = common.PaymentStatus(strings.TrimSpace(req.Payment))
	}
	items, err := s.Repo.FindFirstStockItems(ctx, item.ID)
	if err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Update first stock failed", err.Error())
	}
	item.TotalFirstStock = 0
	for _, line := range items {
		item.TotalFirstStock += line.SubTotal
	}
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateFirstStock(ctx, item); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Update first stock failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "first_stock", item.UserID, item.BranchID, item.TotalFirstStock, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return firststock.FirstStock{}, apperror.New(http.StatusInternalServerError, "Update first stock failed", err.Error())
	}
	return item, nil
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete first stock failed", "first stock not found")
	}
	items, err := s.Repo.FindFirstStockItems(ctx, id)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
	}
	for _, item := range items {
		prod, err := s.Repo.FindProductByID(ctx, item.ProductID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
		}
		prod.Stock -= item.Qty
		if prod.Stock < 0 {
			prod.Stock = 0
		}
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
		}
	}
	for _, item := range items {
		if err := s.Repo.DeleteFirstStockItem(ctx, item.ID); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
		}
	}
	if err := s.Repo.DeleteTransactionReport(ctx, header.ID, "first_stock"); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
	}
	if err := s.Repo.DeleteFirstStock(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock failed", err.Error())
	}
	return nil
}

func (s Service) GetDetail(ctx context.Context, branchID, id string) (firststock.Detail, error) {
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, id)
	if err != nil {
		return firststock.Detail{}, apperror.New(http.StatusNotFound, "Get first stock detail failed", "first stock not found")
	}
	items, err := s.Repo.FindFirstStockItems(ctx, id)
	if err != nil {
		return firststock.Detail{}, apperror.New(http.StatusInternalServerError, "Get first stock detail failed", err.Error())
	}
	return firststock.Detail{ID: header.ID, Description: header.Description, FirstStockDate: header.FirstStockDate.Format("02-01-2006"), TotalFirstStock: header.TotalFirstStock, Payment: string(header.Payment), Items: items}, nil
}

func (s Service) ListItems(ctx context.Context, firstStockID string) ([]firststock.Item, error) {
	items, err := s.Repo.FindFirstStockItems(ctx, firstStockID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "List first stock items failed", err.Error())
	}
	return items, nil
}

func (s Service) CreateItem(ctx context.Context, branchID string, req firststock.CreateItemRequest) (firststock.Item, error) {
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusBadRequest, "Create first stock item failed", "invalid expired_date format, use YYYY-MM-DD")
	}
	_, err = s.Repo.FindFirstStockByID(ctx, branchID, req.FirstStockID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Create first stock item failed", "first stock not found")
	}
	prod, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Create first stock item failed", "product not found")
	}
	unitItem, err := s.Repo.FindUnit(ctx, req.UnitID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Create first stock item failed", "unit not found")
	}
	conversion := 1
	if req.UnitID != prod.UnitID {
		conv, err := s.Repo.FindConversion(ctx, req.ProductID, req.UnitID, prod.UnitID, branchID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
		}
		if err == nil {
			conversion = conv.Value
		}
	}
	actualQty := req.Qty * conversion
	item := firststock.Item{ID: s.IDs.New("FSI"), FirstStockID: req.FirstStockID, ProductID: req.ProductID, ProductName: prod.Name, UnitID: req.UnitID, UnitName: unitItem.Name, Price: prod.PurchasePrice, Qty: req.Qty, SubTotal: prod.PurchasePrice * actualQty, ExpiredDate: expiredDate}
	if err := s.Repo.CreateFirstStockItem(ctx, item); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
	}
	prod.Stock += actualQty
	if expiredDate.Before(prod.ExpiredDate) || prod.ExpiredDate.IsZero() {
		prod.ExpiredDate = expiredDate
	}
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
	}
	total, err := s.Repo.RecalculateFirstStockTotal(ctx, req.FirstStockID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Create first stock item failed", err.Error())
	}
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, req.FirstStockID)
	if err == nil {
		_ = s.Repo.UpsertTransactionReport(ctx, header.ID, "first_stock", header.UserID, header.BranchID, total, string(header.Payment), header.CreatedAt, s.Clock.Now())
	}
	return item, nil
}

func (s Service) UpdateItem(ctx context.Context, branchID, id string, req firststock.UpdateItemRequest) (firststock.Item, error) {
	existing, err := s.Repo.FindFirstStockItemByID(ctx, id)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Update first stock item failed", "item not found")
	}
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusBadRequest, "Update first stock item failed", "invalid expired_date format, use YYYY-MM-DD")
	}
	oldProduct, err := s.Repo.FindProductByID(ctx, existing.ProductID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	oldProduct.Stock -= existing.Qty
	if oldProduct.Stock < 0 {
		oldProduct.Stock = 0
	}
	if err := s.Repo.UpdateProduct(ctx, oldProduct); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	prod, err := s.Repo.FindProductByID(ctx, req.ProductID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Update first stock item failed", "product not found")
	}
	unitItem, err := s.Repo.FindUnit(ctx, req.UnitID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusNotFound, "Update first stock item failed", "unit not found")
	}
	conversion := 1
	if req.UnitID != prod.UnitID {
		conv, err := s.Repo.FindConversion(ctx, req.ProductID, req.UnitID, prod.UnitID, branchID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
		}
		if err == nil {
			conversion = conv.Value
		}
	}
	actualQty := req.Qty * conversion
	existing.ProductID = req.ProductID
	existing.ProductName = prod.Name
	existing.UnitID = req.UnitID
	existing.UnitName = unitItem.Name
	existing.Price = prod.PurchasePrice
	existing.Qty = req.Qty
	existing.SubTotal = prod.PurchasePrice * actualQty
	existing.ExpiredDate = expiredDate
	if err := s.Repo.UpdateFirstStockItem(ctx, existing); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	prod.Stock += actualQty
	if expiredDate.Before(prod.ExpiredDate) || prod.ExpiredDate.IsZero() {
		prod.ExpiredDate = expiredDate
	}
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	total, err := s.Repo.RecalculateFirstStockTotal(ctx, existing.FirstStockID)
	if err != nil {
		return firststock.Item{}, apperror.New(http.StatusInternalServerError, "Update first stock item failed", err.Error())
	}
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, existing.FirstStockID)
	if err == nil {
		_ = s.Repo.UpsertTransactionReport(ctx, header.ID, "first_stock", header.UserID, header.BranchID, total, string(header.Payment), header.CreatedAt, s.Clock.Now())
	}
	return existing, nil
}

func (s Service) DeleteItem(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindFirstStockItemByID(ctx, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete first stock item failed", "item not found")
	}
	prod, err := s.Repo.FindProductByID(ctx, item.ProductID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	prod.Stock -= item.Qty
	if prod.Stock < 0 {
		prod.Stock = 0
	}
	if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	if err := s.Repo.DeleteFirstStockItem(ctx, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	total, err := s.Repo.RecalculateFirstStockTotal(ctx, item.FirstStockID)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete first stock item failed", err.Error())
	}
	header, err := s.Repo.FindFirstStockByID(ctx, branchID, item.FirstStockID)
	if err == nil {
		_ = s.Repo.UpsertTransactionReport(ctx, header.ID, "first_stock", header.UserID, header.BranchID, total, string(header.Payment), header.CreatedAt, s.Clock.Now())
	}
	return nil
}
