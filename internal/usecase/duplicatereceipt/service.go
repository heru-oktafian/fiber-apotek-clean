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
