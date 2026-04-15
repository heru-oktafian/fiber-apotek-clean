package purchase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/purchase"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Repo  ports.PurchaseRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) CreateTransaction(ctx context.Context, branchID, userID string, req purchase.CreatePurchaseRequest) (purchase.Purchase, []purchase.Item, error) {
	now := s.Clock.Now()
	purchaseDate := now
	if req.Purchase.PurchaseDate != "" {
		parsed, err := time.Parse("2006-01-02", req.Purchase.PurchaseDate)
		if err != nil {
			return purchase.Purchase{}, nil, apperror.New(http.StatusBadRequest, "Invalid purchase_date format. Please use `YYYY-MM-DD`.", err)
		}
		purchaseDate = parsed
	}
	payment := common.PaymentStatus(req.Purchase.Payment)
	if payment == "" {
		payment = common.PaymentCash
	}
	p := purchase.Purchase{
		ID:           s.IDs.New("PUR"),
		SupplierID:   req.Purchase.SupplierID,
		PurchaseDate: purchaseDate,
		BranchID:     branchID,
		UserID:       userID,
		Payment:      payment,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	var items []purchase.Item
	if err := s.Repo.WithinTransaction(ctx, func(repo ports.PurchaseTxRepository) error {
		var total int
		for _, input := range req.PurchaseItems {
			prod, err := repo.FindProduct(ctx, input.ProductID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Product with ID %s not found", input.ProductID), err)
			}
			_, err = repo.FindUnit(ctx, input.UnitID)
			if err != nil {
				return apperror.New(http.StatusNotFound, fmt.Sprintf("Unit with ID %s not found", input.UnitID), err)
			}
			conversionValue := 1
			if input.UnitID != prod.UnitID {
				conv, err := repo.FindConversion(ctx, input.ProductID, input.UnitID, prod.UnitID, branchID)
				if err == nil && conv.Value > 0 {
					conversionValue = conv.Value
				}
			}
			expiredDate, err := time.Parse("2006-01-02", input.ExpiredDate)
			if err != nil {
				return apperror.New(http.StatusBadRequest, fmt.Sprintf("Invalid expired_date format for product %s. Please use `YYYY-MM-DD`.", input.ProductID), err)
			}
			itemPrice := input.Price * conversionValue
			itemSubtotal := itemPrice * input.Qty
			items = append(items, purchase.Item{
				ID:          s.IDs.New("PIT"),
				PurchaseID:  p.ID,
				ProductID:   input.ProductID,
				UnitID:      input.UnitID,
				Price:       itemPrice,
				Qty:         input.Qty,
				SubTotal:    itemSubtotal,
				ExpiredDate: expiredDate,
			})
			actualQtyToAdd := input.Qty * conversionValue
			prod.Stock += actualQtyToAdd
			if prod.ExpiredDate.IsZero() || expiredDate.Before(prod.ExpiredDate) {
				prod.ExpiredDate = expiredDate
			}
			if err := repo.UpdateProduct(ctx, prod); err != nil {
				return apperror.New(http.StatusInternalServerError, fmt.Sprintf("Failed to update product details (stock/expired_date) for product %s", prod.Name), err)
			}
			total += itemSubtotal
		}
		p.TotalPurchase = total
		if err := repo.CreatePurchase(ctx, p); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create purchase", err)
		}
		if err := repo.CreatePurchaseItems(ctx, items); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create purchase items", err)
		}
		if err := repo.CreateTransactionReport(ctx, s.IDs.New("TRX"), "purchase", userID, branchID, p.TotalPurchase, string(p.Payment), now); err != nil {
			return apperror.New(http.StatusInternalServerError, "Failed to create transaction report", err)
		}
		return nil
	}); err != nil {
		return purchase.Purchase{}, nil, err
	}
	return p, items, nil
}
