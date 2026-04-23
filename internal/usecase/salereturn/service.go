package salereturn

import (
	"context"
	"net/http"
	"strings"
	"time"

	domain "github.com/heru-oktafian/fiber-apotek-clean/internal/domain/salereturn"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Repo  ports.SaleReturnRepository
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
	return s.Repo.ListSaleReturns(ctx, branchID, req)
}

func (s Service) GetByID(ctx context.Context, branchID, id string) (domain.Detail, error) {
	header, err := s.Repo.FindSaleReturnByID(ctx, branchID, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Get sale return failed", "sale return not found")
	}
	items, err := s.Repo.FindSaleReturnItems(ctx, id)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Get sale return failed", err.Error())
	}
	formatted := make([]domain.FormattedItem, 0, len(items))
	for _, item := range items {
		formatted = append(formatted, domain.FormattedItem{
			ID:           item.ID,
			SaleReturnID: item.SaleReturnID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			UnitID:       item.UnitID,
			UnitName:     item.UnitName,
			Price:        item.Price,
			Qty:          item.Qty,
			SubTotal:     item.SubTotal,
			ExpiredDate:  item.ExpiredDate.Format("2006-01-02"),
		})
	}
	return domain.Detail{ID: header.ID, SaleID: header.SaleID, ReturnDate: header.ReturnDate.Format("2006-01-02"), TotalReturn: header.TotalReturn, Payment: header.Payment, Items: formatted}, nil
}

func (s Service) Create(ctx context.Context, branchID, userID string, req domain.CreateRequest) (domain.Detail, error) {
	saleID := strings.TrimSpace(req.SaleReturn.SaleID)
	if saleID == "" {
		return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "sale_id is required")
	}
	payment := strings.TrimSpace(req.SaleReturn.Payment)
	if payment == "" {
		payment = "paid_by_cash"
	}
	returnDate := s.Clock.Now()
	if strings.TrimSpace(req.SaleReturn.ReturnDate) != "" {
		parsed, err := time.Parse("2006-01-02", req.SaleReturn.ReturnDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "return_date must be in YYYY-MM-DD format")
		}
		returnDate = parsed
	}
	_, err := s.Repo.FindSaleByID(ctx, branchID, saleID)
	if err != nil {
		return domain.Detail{}, apperror.New(http.StatusNotFound, "Create sale return failed", "sale not found")
	}
	headerID := s.IDs.New("SRT")
	header := domain.SaleReturn{ID: headerID, SaleID: saleID, ReturnDate: returnDate, BranchID: branchID, UserID: userID, Payment: payment, CreatedAt: s.Clock.Now(), UpdatedAt: s.Clock.Now()}
	items := make([]domain.Item, 0, len(req.SaleReturnItems))
	total := 0
	for _, in := range req.SaleReturnItems {
		saleItem, err := s.Repo.FindSaleItemBySaleAndProduct(ctx, saleID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "product not found in source sale")
		}
		returnedQty, err := s.Repo.SumSaleReturnedQty(ctx, saleID, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
		}
		if returnedQty+in.Qty > saleItem.Qty {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "return qty exceeds sold qty")
		}
		expiredDate, err := time.Parse("2006-01-02", in.ExpiredDate)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusBadRequest, "Create sale return failed", "expired_date must be in YYYY-MM-DD format")
		}
		prod, err := s.Repo.FindProductByID(ctx, in.ProductID)
		if err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
		}
		prod.Stock += in.Qty
		if err := s.Repo.UpdateProduct(ctx, prod); err != nil {
			return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
		}
		subTotal := saleItem.Price * in.Qty
		total += subTotal
		items = append(items, domain.Item{ID: s.IDs.New("SRI"), SaleReturnID: headerID, ProductID: in.ProductID, Price: saleItem.Price, Qty: in.Qty, SubTotal: subTotal, ExpiredDate: expiredDate})
	}
	header.TotalReturn = total
	if err := s.Repo.CreateSaleReturn(ctx, header); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
	}
	if err := s.Repo.CreateSaleReturnItems(ctx, items); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
	}
	if err := s.Repo.CreateTransactionReport(ctx, header.ID, "sale_return", userID, branchID, total, payment, s.Clock.Now()); err != nil {
		return domain.Detail{}, apperror.New(http.StatusInternalServerError, "Create sale return failed", err.Error())
	}
	return s.GetByID(ctx, branchID, headerID)
}
