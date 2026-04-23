package duplicatereceipt

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/duplicatereceipt"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	"gorm.io/gorm"
)

type Service struct {
	Repo  ports.DuplicateReceiptRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req domain.ListRequest) (domain.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	result, err := s.Repo.ListDuplicateReceipts(ctx, branchID, req)
	if err != nil {
		return domain.ListResult{}, apperror.New(http.StatusInternalServerError, "Get duplicate receipts failed", err.Error())
	}
	return result, nil
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (domain.Detail, error) {
	item, err := s.Repo.FindDuplicateReceiptByID(ctx, branchID, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Get duplicate receipt failed", "duplicate receipt not found")
	}
	return domain.Detail{ID: item.ID, MemberID: item.MemberID, MemberName: item.MemberName, Description: item.Description, DuplicateReceiptDate: item.DuplicateReceiptDate.Format("2006-01-02"), TotalDuplicateReceipt: item.TotalDuplicateReceipt, ProfitEstimate: item.ProfitEstimate, Payment: string(item.Payment)}, nil
}

func (s Service) Create(ctx context.Context, branchID, userID, defaultMember string, req domain.CreateRequest) (domain.Detail, error) {
	date, err := time.Parse("2006-01-02", req.DuplicateReceipt.DuplicateReceiptDate)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create duplicate receipt failed", "duplicate_receipt_date must be in YYYY-MM-DD format")
	}
	payment := common.PaymentStatus(req.DuplicateReceipt.Payment)
	if payment == "" {
		payment = common.PaymentCash
	}
	memberID := req.DuplicateReceipt.MemberID
	if memberID == "" {
		memberID = defaultMember
	}
	if memberID != "" {
		if _, err := s.Repo.FindMemberByID(ctx, memberID); err != nil {
			memberID = defaultMember
		}
	}
	headerID := s.IDs.New("DUR")
	now := s.Clock.Now()
	header := domain.DuplicateReceipt{ID: headerID, MemberID: memberID, Description: req.DuplicateReceipt.Description, DuplicateReceiptDate: date, Payment: payment, BranchID: branchID, UserID: userID, CreatedAt: now, UpdatedAt: now}
	if err := s.Repo.WithinTransactionDuplicateReceipt(ctx, func(repo ports.DuplicateReceiptTxRepository) error {
		items := make([]domain.Item, 0, len(req.Items))
		for _, in := range req.Items {
			prod, err := repo.FindProduct(ctx, in.ProductID)
			if err != nil {
				return apperror.New(http.StatusNotFound, "Create duplicate receipt failed", "product not found")
			}
			if prod.Stock < in.Qty {
				return apperror.New(http.StatusBadRequest, "Create duplicate receipt failed", fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d", prod.Name, prod.Stock, in.Qty))
			}
			price := prod.SalesPrice
			subTotal := price * in.Qty
			header.TotalDuplicateReceipt += subTotal
			header.ProfitEstimate += (price - prod.PurchasePrice) * in.Qty
			items = append(items, domain.Item{ID: s.IDs.New("DRI"), DuplicateReceiptID: headerID, ProductID: in.ProductID, Price: price, Qty: in.Qty, SubTotal: subTotal})
			prod.Stock -= in.Qty
			if err := repo.UpdateProduct(ctx, prod); err != nil {
				return apperror.New(http.StatusInternalServerError, "Create duplicate receipt failed", err.Error())
			}
		}
		if err := repo.CreateDuplicateReceipt(ctx, header); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt failed", err.Error())
		}
		if err := repo.CreateDuplicateReceiptItems(ctx, items); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt failed", err.Error())
		}
		if err := repo.CreateTransactionReport(ctx, header.ID, "sale", userID, branchID, header.TotalDuplicateReceipt, string(payment), now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt failed", err.Error())
		}
		if err := repo.UpsertDailyProfit(ctx, date, userID, branchID, header.TotalDuplicateReceipt, header.ProfitEstimate, now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt failed", err.Error())
		}
		if memberID != "" && memberID != defaultMember {
			member, err := repo.FindMember(ctx, memberID)
			if err != nil {
				return apperror.New(http.StatusNotFound, "Create duplicate receipt failed", "member not found")
			}
			category, err := repo.FindMemberCategory(ctx, member.MemberCategoryID)
			if err != nil {
				return apperror.New(http.StatusNotFound, "Create duplicate receipt failed", "member category not found")
			}
			if category.PointsConversionRate > 0 {
				points := header.TotalDuplicateReceipt / category.PointsConversionRate
				if err := repo.UpdateMemberPoints(ctx, member.ID, member.Points+points); err != nil {
					return apperror.New(http.StatusInternalServerError, "Create duplicate receipt failed", err.Error())
				}
			}
		}
		return nil
	}); err != nil {
		return domain.Detail{}, err
	}
	return s.GetByID(ctx, branchID, headerID)
}

func (s Service) Update(ctx context.Context, branchID, id, defaultMember string, req domain.UpdateRequest) (domain.DuplicateReceipt, error) {
	item, err := s.Repo.FindDuplicateReceiptByID(ctx, branchID, id)
	if err != nil {
		return domain.DuplicateReceipt{}, apperror.New(http.StatusNotFound, "Update duplicate receipt failed", "duplicate receipt not found")
	}
	oldTotal := item.TotalDuplicateReceipt
	oldProfit := item.ProfitEstimate
	if req.MemberID != nil {
		memberID := *req.MemberID
		if memberID == "" {
			memberID = defaultMember
		}
		if _, err := s.Repo.FindMemberByID(ctx, memberID); err == nil {
			item.MemberID = memberID
		}
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Payment != "" {
		item.Payment = common.PaymentStatus(req.Payment)
	}
	items, err := s.Repo.FindDuplicateReceiptItems(ctx, id)
	if err != nil {
		return domain.DuplicateReceipt{}, apperror.New(http.StatusInternalServerError, "Update duplicate receipt failed", err.Error())
	}
	item.TotalDuplicateReceipt = 0
	item.ProfitEstimate = 0
	for _, line := range items {
		item.TotalDuplicateReceipt += line.SubTotal
		prod, prodErr := s.Repo.FindProductByID(ctx, line.ProductID)
		if prodErr != nil {
			return domain.DuplicateReceipt{}, apperror.New(http.StatusInternalServerError, "Update duplicate receipt failed", prodErr.Error())
		}
		item.ProfitEstimate += (line.Price - prod.PurchasePrice) * line.Qty
	}
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateDuplicateReceipt(ctx, item); err != nil {
		return domain.DuplicateReceipt{}, apperror.New(http.StatusInternalServerError, "Update duplicate receipt failed", err.Error())
	}
	if err := s.Repo.UpdateTransactionReport(ctx, item.ID, item.TotalDuplicateReceipt, string(item.Payment), item.UpdatedAt); err != nil {
		return domain.DuplicateReceipt{}, apperror.New(http.StatusInternalServerError, "Update duplicate receipt failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, item.DuplicateReceiptDate, item.UserID, item.BranchID, item.TotalDuplicateReceipt-oldTotal, item.ProfitEstimate-oldProfit, item.UpdatedAt); err != nil {
		return domain.DuplicateReceipt{}, apperror.New(http.StatusInternalServerError, "Update duplicate receipt failed", err.Error())
	}
	return s.Repo.FindDuplicateReceiptByID(ctx, branchID, id)
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindDuplicateReceiptByID(ctx, branchID, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete duplicate receipt failed", "duplicate receipt not found")
	}
	items, err := s.Repo.FindDuplicateReceiptItems(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt failed", err.Error())
	}
	now := s.Clock.Now()
	if err := s.Repo.WithinTransactionDuplicateReceipt(ctx, func(repo ports.DuplicateReceiptTxRepository) error {
		for _, line := range items {
			prod, err := repo.FindProduct(ctx, line.ProductID)
			if err != nil {
				return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt failed", err.Error())
			}
			prod.Stock += line.Qty
			if err := repo.UpdateProduct(ctx, prod); err != nil {
				return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt failed", err.Error())
			}
		}
		if err := repo.DeleteDuplicateReceiptItems(ctx, id); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt failed", err.Error())
		}
		if err := repo.DeleteTransactionReport(ctx, id, "sale"); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt failed", err.Error())
		}
		if err := repo.DeleteDuplicateReceipt(ctx, branchID, id); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt failed", err.Error())
		}
		return nil
	}); err != nil {
		return err
	}
	if err := s.Repo.AdjustDailyProfit(ctx, item.DuplicateReceiptDate, item.UserID, item.BranchID, -item.TotalDuplicateReceipt, -item.ProfitEstimate, now); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt failed", err.Error())
	}
	return nil
}

func (s Service) ListItems(ctx context.Context, branchID, duplicateReceiptID string) ([]domain.Item, error) {
	if _, err := s.Repo.FindDuplicateReceiptByID(ctx, branchID, duplicateReceiptID); err != nil {
		return nil, apperror.New(http.StatusNotFound, "Get duplicate receipt items failed", "duplicate receipt not found")
	}
	items, err := s.Repo.FindDuplicateReceiptItems(ctx, duplicateReceiptID)
	if err != nil {
		return nil, apperror.New(http.StatusInternalServerError, "Get duplicate receipt items failed", err.Error())
	}
	return items, nil
}

func (s Service) CreateItem(ctx context.Context, branchID string, req domain.CreateItemRequest) (domain.Item, error) {
	header, err := s.Repo.FindDuplicateReceiptByID(ctx, branchID, req.DuplicateReceiptID)
	if err != nil {
		return domain.Item{}, apperror.New(http.StatusNotFound, "Create duplicate receipt item failed", "duplicate receipt not found")
	}
	items, err := s.Repo.FindDuplicateReceiptItems(ctx, req.DuplicateReceiptID)
	if err != nil {
		return domain.Item{}, apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
	}
	oldTotal := header.TotalDuplicateReceipt
	oldProfit := header.ProfitEstimate
	var resultID string
	if err := s.Repo.WithinTransactionDuplicateReceipt(ctx, func(repo ports.DuplicateReceiptTxRepository) error {
		prod, err := repo.FindProduct(ctx, req.ProductID)
		if err != nil {
			return apperror.New(http.StatusNotFound, "Create duplicate receipt item failed", "product not found")
		}
		if prod.Stock < req.Qty {
			return apperror.New(http.StatusBadRequest, "Create duplicate receipt item failed", fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d", prod.Name, prod.Stock, req.Qty))
		}
		for _, existing := range items {
			if existing.ProductID == req.ProductID {
				existing.Qty += req.Qty
				existing.Price = prod.SalesPrice
				existing.SubTotal = existing.Qty * existing.Price
				if err := repo.UpdateDuplicateReceiptItem(ctx, existing); err != nil {
					return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
				}
				prod.Stock -= req.Qty
				if err := repo.UpdateProduct(ctx, prod); err != nil {
					return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
				}
				header.TotalDuplicateReceipt += prod.SalesPrice * req.Qty
				header.ProfitEstimate += (prod.SalesPrice - prod.PurchasePrice) * req.Qty
				header.UpdatedAt = s.Clock.Now()
				if err := s.Repo.UpdateDuplicateReceipt(ctx, header); err != nil {
					return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
				}
				if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalDuplicateReceipt, string(header.Payment), header.UpdatedAt); err != nil {
					return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
				}
				resultID = existing.ID
				return nil
			}
		}
		item := domain.Item{ID: s.IDs.New("DRI"), DuplicateReceiptID: req.DuplicateReceiptID, ProductID: req.ProductID, Price: prod.SalesPrice, Qty: req.Qty, SubTotal: prod.SalesPrice * req.Qty}
		if err := repo.CreateDuplicateReceiptItem(ctx, item); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
		}
		prod.Stock -= req.Qty
		if err := repo.UpdateProduct(ctx, prod); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
		}
		header.TotalDuplicateReceipt += item.SubTotal
		header.ProfitEstimate += (item.Price - prod.PurchasePrice) * item.Qty
		header.UpdatedAt = s.Clock.Now()
		if err := s.Repo.UpdateDuplicateReceipt(ctx, header); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
		}
		if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalDuplicateReceipt, string(header.Payment), header.UpdatedAt); err != nil {
			return apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
		}
		resultID = item.ID
		return nil
	}); err != nil {
		return domain.Item{}, err
	}
	if err := s.Repo.AdjustDailyProfit(ctx, header.DuplicateReceiptDate, header.UserID, header.BranchID, header.TotalDuplicateReceipt-oldTotal, header.ProfitEstimate-oldProfit, header.UpdatedAt); err != nil {
		return domain.Item{}, apperror.New(http.StatusInternalServerError, "Create duplicate receipt item failed", err.Error())
	}
	return s.Repo.FindDuplicateReceiptItemByID(ctx, resultID)
}

func (s Service) UpdateItem(ctx context.Context, branchID, id string, req domain.UpdateItemRequest) (domain.Item, error) {
	item, err := s.Repo.FindDuplicateReceiptItemByID(ctx, id)
	if err != nil {
		return domain.Item{}, apperror.New(http.StatusNotFound, "Update duplicate receipt item failed", "item not found")
	}
	header, err := s.Repo.FindDuplicateReceiptByID(ctx, branchID, item.DuplicateReceiptID)
	if err != nil {
		return domain.Item{}, apperror.New(http.StatusNotFound, "Update duplicate receipt item failed", "duplicate receipt not found")
	}
	oldTotal := header.TotalDuplicateReceipt
	oldProfit := header.ProfitEstimate
	if err := s.Repo.WithinTransactionDuplicateReceipt(ctx, func(repo ports.DuplicateReceiptTxRepository) error {
		oldProduct, err := repo.FindProduct(ctx, item.ProductID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
		}
		oldProduct.Stock += item.Qty
		if err := repo.UpdateProduct(ctx, oldProduct); err != nil {
			return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
		}
		newProduct, err := repo.FindProduct(ctx, req.ProductID)
		if err != nil {
			return apperror.New(http.StatusNotFound, "Update duplicate receipt item failed", "product not found")
		}
		if newProduct.Stock < req.Qty {
			return apperror.New(http.StatusBadRequest, "Update duplicate receipt item failed", fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d", newProduct.Name, newProduct.Stock, req.Qty))
		}
		newProduct.Stock -= req.Qty
		if err := repo.UpdateProduct(ctx, newProduct); err != nil {
			return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
		}
		item.ProductID = req.ProductID
		item.Price = newProduct.SalesPrice
		item.Qty = req.Qty
		item.SubTotal = newProduct.SalesPrice * req.Qty
		if err := repo.UpdateDuplicateReceiptItem(ctx, item); err != nil {
			return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
		}
		items, err := s.Repo.FindDuplicateReceiptItems(ctx, item.DuplicateReceiptID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
		}
		header.TotalDuplicateReceipt = 0
		header.ProfitEstimate = 0
		for _, line := range items {
			header.TotalDuplicateReceipt += line.SubTotal
			prod, prodErr := repo.FindProduct(ctx, line.ProductID)
			if prodErr != nil {
				return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", prodErr.Error())
			}
			header.ProfitEstimate += (line.Price - prod.PurchasePrice) * line.Qty
		}
		header.UpdatedAt = s.Clock.Now()
		if err := s.Repo.UpdateDuplicateReceipt(ctx, header); err != nil {
			return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
		}
		if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalDuplicateReceipt, string(header.Payment), header.UpdatedAt); err != nil {
			return apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
		}
		return nil
	}); err != nil {
		return domain.Item{}, err
	}
	if err := s.Repo.AdjustDailyProfit(ctx, header.DuplicateReceiptDate, header.UserID, header.BranchID, header.TotalDuplicateReceipt-oldTotal, header.ProfitEstimate-oldProfit, header.UpdatedAt); err != nil {
		return domain.Item{}, apperror.New(http.StatusInternalServerError, "Update duplicate receipt item failed", err.Error())
	}
	return s.Repo.FindDuplicateReceiptItemByID(ctx, id)
}

func (s Service) DeleteItem(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindDuplicateReceiptItemByID(ctx, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete duplicate receipt item failed", "item not found")
	}
	header, err := s.Repo.FindDuplicateReceiptByID(ctx, branchID, item.DuplicateReceiptID)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete duplicate receipt item failed", "duplicate receipt not found")
	}
	oldTotal := header.TotalDuplicateReceipt
	oldProfit := header.ProfitEstimate
	if err := s.Repo.WithinTransactionDuplicateReceipt(ctx, func(repo ports.DuplicateReceiptTxRepository) error {
		prod, err := repo.FindProduct(ctx, item.ProductID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", err.Error())
		}
		prod.Stock += item.Qty
		if err := repo.UpdateProduct(ctx, prod); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", err.Error())
		}
		if err := repo.DeleteDuplicateReceiptItem(ctx, id); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", err.Error())
		}
		items, err := s.Repo.FindDuplicateReceiptItems(ctx, item.DuplicateReceiptID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", err.Error())
		}
		header.TotalDuplicateReceipt = 0
		header.ProfitEstimate = 0
		for _, line := range items {
			header.TotalDuplicateReceipt += line.SubTotal
			lineProduct, lineErr := repo.FindProduct(ctx, line.ProductID)
			if lineErr != nil {
				return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", lineErr.Error())
			}
			header.ProfitEstimate += (line.Price - lineProduct.PurchasePrice) * line.Qty
		}
		header.UpdatedAt = s.Clock.Now()
		if err := s.Repo.UpdateDuplicateReceipt(ctx, header); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", err.Error())
		}
		if err := s.Repo.UpdateTransactionReport(ctx, header.ID, header.TotalDuplicateReceipt, string(header.Payment), header.UpdatedAt); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", err.Error())
		}
		return nil
	}); err != nil {
		return err
	}
	if err := s.Repo.AdjustDailyProfit(ctx, header.DuplicateReceiptDate, header.UserID, header.BranchID, header.TotalDuplicateReceipt-oldTotal, header.ProfitEstimate-oldProfit, header.UpdatedAt); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete duplicate receipt item failed", err.Error())
	}
	return nil
}
