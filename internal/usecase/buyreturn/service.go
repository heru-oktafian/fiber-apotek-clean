package buyreturn

import (
	"context"
	"net/http"
	"strings"
	"time"

	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/buyreturn"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Repo  ports.BuyReturnRepository
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
	return s.Repo.ListBuyReturns(ctx, branchID, req)
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (domain.Detail, error) {
	header, err := s.Repo.FindBuyReturnByID(ctx, branchID, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Get buy return failed", "buy return not found")
	}
	items, err := s.Repo.FindBuyReturnItems(ctx, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Get buy return failed", err.Error())
	}
	formatted := make([]domain.FormattedItem, 0, len(items))
	for _, item := range items {
		formatted = append(formatted, domain.FormattedItem{
			ID:          item.ID,
			BuyReturnID: item.BuyReturnID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			UnitID:      item.UnitID,
			UnitName:    item.UnitName,
			Price:       item.Price,
			Qty:         item.Qty,
			SubTotal:    item.SubTotal,
			ExpiredDate: item.ExpiredDate.Format("2006-01-02"),
		})
	}
	return domain.Detail{ID: header.ID, PurchaseID: header.PurchaseID, ReturnDate: header.ReturnDate.Format("2006-01-02"), TotalReturn: header.TotalReturn, Payment: header.Payment, Items: formatted}, nil
}

func (s Service) ListPurchaseSources(ctx context.Context, branchID, search, month string) ([]domain.PurchaseComboItem, error) {
	return s.Repo.ListPurchaseReturnSources(ctx, branchID, search, month)
}

func (s Service) ListReturnableItems(ctx context.Context, purchaseID string) ([]domain.ReturnableItem, error) {
	if strings.TrimSpace(purchaseID) == "" {
		return nil, apperror.New(http.StatusBadRequest, "Get buy return items failed", "purchase_id is required")
	}
	return s.Repo.ListPurchaseReturnableItems(ctx, purchaseID)
}

func (s Service) Create(ctx context.Context, branchID, userID string, req domain.CreateRequest) (domain.Detail, error) {
	purchaseID := strings.TrimSpace(req.BuyReturn.PurchaseID)
	if purchaseID == "" {
		return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "purchase_id is required")
	}
	payment := strings.TrimSpace(req.BuyReturn.Payment)
	if payment == "" {
		payment = "paid_by_cash"
	}
	returnDate := s.Clock.Now()
	if strings.TrimSpace(req.BuyReturn.ReturnDate) != "" {
		parsed, err := time.Parse("2006-01-02", req.BuyReturn.ReturnDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "return_date must be in YYYY-MM-DD format")
		}
		returnDate = parsed
	}
	_, err := s.Repo.FindPurchaseByID(ctx, branchID, purchaseID)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Create buy return failed", "purchase not found")
	}
	headerID := s.IDs.New("BRT")
	header := domain.BuyReturn{ID: headerID, PurchaseID: purchaseID, ReturnDate: returnDate, BranchID: branchID, UserID: userID, Payment: payment, CreatedAt: s.Clock.Now(), UpdatedAt: s.Clock.Now()}
	items := make([]domain.Item, 0, len(req.BuyReturnItems))
	total := 0
	for _, in := range req.BuyReturnItems {
		purchaseItem, err := s.Repo.FindPurchaseItemByPurchaseAndProduct(ctx, purchaseID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "product not found in source purchase")
		}
		returnedQty, err := s.Repo.SumBuyReturnedQty(ctx, purchaseID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
		}
		if returnedQty+in.Qty > purchaseItem.Qty {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "return qty exceeds purchased qty")
		}
		expiredDate, err := time.Parse("2006-01-02", in.ExpiredDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create buy return failed", "expired_date must be in YYYY-MM-DD format")
		}
		prod, err := s.Repo.FindProductByID(ctx, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
		}
		actualQty := in.Qty
		if purchaseItem.UnitID != "" && prod.UnitID != "" && purchaseItem.UnitID != prod.UnitID {
			conv, err := s.Repo.FindConversion(ctx, in.ProductID, purchaseItem.UnitID, prod.UnitID, branchID)
			if err == nil && conv.Value > 0 {
				actualQty = in.Qty * conv.Value
			}
		}
		prod.Stock -= actualQty
		if prod.Stock < 0 {
			prod.Stock = 0
		}
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
		}
		subTotal := purchaseItem.Price * in.Qty
		total += subTotal
		items = append(items, domain.Item{ID: s.IDs.New("BRI"), BuyReturnID: headerID, ProductID: in.ProductID, Price: purchaseItem.Price, Qty: in.Qty, SubTotal: subTotal, ExpiredDate: expiredDate})
	}
	header.TotalReturn = total
	if err := s.Repo.CreateBuyReturn(ctx, header); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
	}
	if err := s.Repo.CreateBuyReturnItems(ctx, items); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
	}
	if err := s.Repo.CreateTransactionReport(ctx, header.ID, "buy_return", userID, branchID, total, payment, s.Clock.Now()); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create buy return failed", err.Error())
	}
	return s.GetByID(ctx, branchID, headerID)
}
