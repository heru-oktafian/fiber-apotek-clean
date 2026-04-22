package anotherincome

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/anotherincome"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Repo  ports.AnotherIncomeRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req anotherincome.ListRequest) (anotherincome.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	return s.Repo.ListAnotherIncomes(ctx, branchID, req)
}

func (s Service) Create(ctx context.Context, branchID, userID string, req anotherincome.CreateRequest) (anotherincome.AnotherIncome, error) {
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.IncomeDate))
	if err != nil {
		return anotherincome.AnotherIncome{}, apperror.New(http.StatusBadRequest, "Create another income failed", "invalid date format, use YYYY-MM-DD")
	}
	now := s.Clock.Now()
	item := anotherincome.AnotherIncome{
		ID:          s.IDs.New("ANI"),
		Description: strings.TrimSpace(req.Description),
		IncomeDate:  parsedDate,
		BranchID:    branchID,
		UserID:      userID,
		TotalIncome: req.TotalIncome,
		Payment:     common.PaymentStatus(strings.TrimSpace(req.Payment)),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if item.Payment == "" {
		item.Payment = common.PaymentCash
	}
	if err := s.Repo.CreateAnotherIncome(ctx, item); err != nil {
		return anotherincome.AnotherIncome{}, apperror.New(http.StatusInternalServerError, "Create another income failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "income", item.UserID, item.BranchID, item.TotalIncome, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return anotherincome.AnotherIncome{}, apperror.New(http.StatusInternalServerError, "Create another income failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id string, req anotherincome.UpdateRequest) (anotherincome.AnotherIncome, error) {
	item, err := s.Repo.FindAnotherIncomeByID(ctx, branchID, id)
	if err != nil {
		return anotherincome.AnotherIncome{}, apperror.New(http.StatusNotFound, "Update another income failed", "another income not found")
	}
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.IncomeDate))
	if err != nil {
		return anotherincome.AnotherIncome{}, apperror.New(http.StatusBadRequest, "Update another income failed", "invalid date format, use YYYY-MM-DD")
	}
	item.Description = strings.TrimSpace(req.Description)
	item.IncomeDate = parsedDate
	item.TotalIncome = req.TotalIncome
	item.Payment = common.PaymentStatus(strings.TrimSpace(req.Payment))
	if item.Payment == "" {
		item.Payment = common.PaymentCash
	}
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateAnotherIncome(ctx, item); err != nil {
		return anotherincome.AnotherIncome{}, apperror.New(http.StatusInternalServerError, "Update another income failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "income", item.UserID, item.BranchID, item.TotalIncome, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return anotherincome.AnotherIncome{}, apperror.New(http.StatusInternalServerError, "Update another income failed", err.Error())
	}
	return item, nil
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindAnotherIncomeByID(ctx, branchID, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete another income failed", "another income not found")
	}
	if err := s.Repo.DeleteTransactionReport(ctx, item.ID, "income"); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete another income failed", err.Error())
	}
	if err := s.Repo.DeleteAnotherIncome(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete another income failed", err.Error())
	}
	return nil
}
