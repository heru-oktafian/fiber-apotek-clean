package sale

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/sale"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Repo  ports.SaleRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req sale.ListRequest) (sale.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	items, err := s.Repo.ListSales(ctx, branchID, req)
	if err != nil {
		return sale.ListResult{}, apperror.New(http.StatusInternalServerError, "Get sales failed", err.Error())
	}
	return items, nil
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (sale.Detail, error) {
	item, err := s.Repo.FindSaleDetail(ctx, branchID, id)
	if err != nil {
		return sale.Detail{}, apperror.New(http.StatusNotFound, "Get sale failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id, defaultMember string, req sale.UpdateRequest) (sale.Sale, error) {
	item, err := s.Repo.FindSaleByID(ctx, branchID, id)
	if err != nil {
		return sale.Sale{}, apperror.New(http.StatusNotFound, "Update sale failed", "sale not found")
	}
	oldTotal := item.TotalSale
	oldProfit := item.ProfitEstimate
	if req.MemberID != nil {
		memberID := *req.MemberID
		if memberID == "" {
			memberID = defaultMember
		}
		if _, err := s.Repo.FindMemberByID(ctx, memberID); err != nil {
			memberID = defaultMember
		}
		item.MemberID = memberID
	}
	if req.Payment != "" {
		item.Payment = common.PaymentStatus(req.Payment)
	}
	if req.Discount != nil {
		item.Discount = *req.Discount
	}
	items, err := s.Repo.FindSaleItems(ctx, id)
	if err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	var total int
	var profit int
	for _, line := range items {
		total += line.SubTotal
		profit += line.SubTotal - (line.Price * line.Qty)
	}
	item.TotalSale = total - item.Discount
	item.ProfitEstimate = profit
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateSaleHeader(ctx, item); err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	if err := s.Repo.UpdateTransactionReport(ctx, item.ID, item.TotalSale, string(item.Payment), item.UpdatedAt); err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, item.SaleDate, item.UserID, item.BranchID, item.TotalSale-oldTotal, item.ProfitEstimate-oldProfit, item.UpdatedAt); err != nil {
		return sale.Sale{}, apperror.New(http.StatusInternalServerError, "Update sale failed", err.Error())
	}
	return s.Repo.FindSaleByID(ctx, branchID, id)
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindSaleByID(ctx, branchID, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete sale failed", "sale not found")
	}
	items, err := s.Repo.FindSaleItems(ctx, id)
	if err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	for _, line := range items {
		prod, err := s.Repo.FindProductByID(ctx, line.ProductID)
		if err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
		}
		prod.Stock += line.Qty
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
		}
	}
	now := s.Clock.Now()
	if err := s.Repo.DeleteSaleItems(ctx, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	if err := s.Repo.DeleteTransactionReport(ctx, id, "sale"); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	if err := s.Repo.DeleteSaleHeader(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	if err := s.Repo.AdjustDailyProfit(ctx, item.SaleDate, item.UserID, item.BranchID, -item.TotalSale, -item.ProfitEstimate, now); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete sale failed", err.Error())
	}
	return nil
}

func (s Service) CreateTransaction(ctx context.Context, branchID, userID, defaultMember, subscriptionType string, req sale.CreateSaleRequest) (sale.Sale, []sale.Item, error) {
	now := s.Clock.Now()
	payment := common.PaymentStatus(req.Sale.Payment)
	if payment == "" {
		payment = common.PaymentCash
	}
	memberID := req.Sale.MemberID
	if memberID == "" {
		memberID = defaultMember
	}
	saleEntity := sale.Sale{
		ID:        s.IDs.New("SAL"),
		MemberID:  memberID,
		UserID:    userID,
		BranchID:  branchID,
		Payment:   payment,
		Discount:  req.Sale.Discount,
		SaleDate:  now,
		CreatedAt: now,
		UpdatedAt: now,
	}
	items := make([]sale.Item, 0, len(req.SaleItems))
	if err := s.Repo.WithinTransactionSale(ctx, func(repo ports.SaleTxRepository) error {
		var totalSale int
		var profitEstimate int
		for _, input := range req.SaleItems {
			prod, err := repo.FindProduct(ctx, input.ProductID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Product with ID %s not found", input.ProductID), err)
			}
			if prod.Stock < input.Qty {
				return apperror.New(http.StatusBadRequest, fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d", prod.Name, prod.Stock, input.Qty), nil)
			}
			serverPrice := prod.SalesPrice
			subtotal := serverPrice * input.Qty
			items = append(items, sale.Item{
				ID:        s.IDs.New("SIT"),
				SaleID:    saleEntity.ID,
				ProductID: input.ProductID,
				Price:     serverPrice,
				Qty:       input.Qty,
				SubTotal:  subtotal,
			})
			prod.Stock -= input.Qty
			if err := repo.UpdateProduct(ctx, prod); err != nil {
				return apperror.New(http.StatusInternalServerError, fmt.Sprintf("Failed to update stock for product %s", prod.Name), err)
			}
			totalSale += subtotal
			profitEstimate += (serverPrice - prod.PurchasePrice) * input.Qty
		}
		saleEntity.TotalSale = totalSale - saleEntity.Discount
		saleEntity.ProfitEstimate = profitEstimate
		if err := repo.CreateSale(ctx, saleEntity); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create sale", err)
		}
		if err := repo.CreateSaleItems(ctx, items); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create sale items", err)
		}
		if err := repo.CreateTransactionReport(ctx, saleEntity.ID, "sale", userID, branchID, saleEntity.TotalSale, string(saleEntity.Payment), now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create transaction report", err)
		}
		if err := repo.UpsertDailyProfit(ctx, saleEntity.SaleDate, userID, branchID, saleEntity.TotalSale, saleEntity.ProfitEstimate, now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to update daily profit report", err)
		}
		if subscriptionType == "quota" {
			br, err := repo.FindBranch(ctx, branchID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Branch with ID %s not found", branchID), err)
			}
			if br.Quota <= 0 {
				return apperror.New(http.StatusBadRequest, fmt.Sprintf("No quota available for branch %s", br.BranchName), nil)
			}
			if err := repo.UpdateBranchQuota(ctx, branchID, br.Quota-1); err != nil {
				return apperror.New(http.StatusInternalServerError, fmt.Sprintf("Failed to update quota for branch %s", br.BranchName), err)
			}
		}
		if memberID != "" && memberID != defaultMember {
			m, err := repo.FindMember(ctx, memberID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Member with ID %s not found", memberID), err)
			}
			mc, err := repo.FindMemberCategory(ctx, m.MemberCategoryID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Member category with ID %s not found for member %s", m.MemberCategoryID, m.ID), err)
			}
			if mc.PointsConversionRate > 0 {
				points := saleEntity.TotalSale / mc.PointsConversionRate
				if err := repo.UpdateMemberPoints(ctx, m.ID, m.Points+points); err != nil {
					return apperror.New(http.StatusInternalServerError, "Failed to update member points", err)
				}
			}
		}
		return nil
	}); err != nil {
		return sale.Sale{}, nil, err
	}
	return saleEntity, items, nil
}
