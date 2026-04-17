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
