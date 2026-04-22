package expense

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/expense"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
)

type Service struct {
	Repo  ports.ExpenseRepository
	IDs   ports.IDGenerator
	Clock ports.Clock
}

func (s Service) List(ctx context.Context, branchID string, req expense.ListRequest) (expense.ListResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	return s.Repo.ListExpenses(ctx, branchID, req)
}

func (s Service) Create(ctx context.Context, branchID, userID string, req expense.CreateRequest) (expense.Expense, error) {
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.ExpenseDate))
	if err != nil {
		return expense.Expense{}, apperror.New(http.StatusBadRequest, "Create expense failed", "invalid date format, use YYYY-MM-DD")
	}
	now := s.Clock.Now()
	item := expense.Expense{
		ID:           s.IDs.New("EXP"),
		Description:  strings.TrimSpace(req.Description),
		ExpenseDate:  parsedDate,
		BranchID:     branchID,
		UserID:       userID,
		TotalExpense: req.TotalExpense,
		Payment:      common.PaymentStatus(strings.TrimSpace(req.Payment)),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if item.Payment == "" {
		item.Payment = common.PaymentCash
	}
	if err := s.Repo.CreateExpense(ctx, item); err != nil {
		return expense.Expense{}, apperror.New(http.StatusInternalServerError, "Create expense failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "expense", item.UserID, item.BranchID, item.TotalExpense, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return expense.Expense{}, apperror.New(http.StatusInternalServerError, "Create expense failed", err.Error())
	}
	return item, nil
}

func (s Service) Update(ctx context.Context, branchID, id string, req expense.UpdateRequest) (expense.Expense, error) {
	item, err := s.Repo.FindExpenseByID(ctx, branchID, id)
	if err != nil {
		return expense.Expense{}, apperror.New(http.StatusNotFound, "Update expense failed", "expense not found")
	}
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.ExpenseDate))
	if err != nil {
		return expense.Expense{}, apperror.New(http.StatusBadRequest, "Update expense failed", "invalid date format, use YYYY-MM-DD")
	}
	item.Description = strings.TrimSpace(req.Description)
	item.ExpenseDate = parsedDate
	item.TotalExpense = req.TotalExpense
	item.Payment = common.PaymentStatus(strings.TrimSpace(req.Payment))
	if item.Payment == "" {
		item.Payment = common.PaymentCash
	}
	item.UpdatedAt = s.Clock.Now()
	if err := s.Repo.UpdateExpense(ctx, item); err != nil {
		return expense.Expense{}, apperror.New(http.StatusInternalServerError, "Update expense failed", err.Error())
	}
	if err := s.Repo.UpsertTransactionReport(ctx, item.ID, "expense", item.UserID, item.BranchID, item.TotalExpense, string(item.Payment), item.CreatedAt, item.UpdatedAt); err != nil {
		return expense.Expense{}, apperror.New(http.StatusInternalServerError, "Update expense failed", err.Error())
	}
	return item, nil
}

func (s Service) Delete(ctx context.Context, branchID, id string) error {
	item, err := s.Repo.FindExpenseByID(ctx, branchID, id)
	if err != nil {
		return apperror.New(http.StatusNotFound, "Delete expense failed", "expense not found")
	}
	if err := s.Repo.DeleteTransactionReport(ctx, item.ID, "expense"); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete expense failed", err.Error())
	}
	if err := s.Repo.DeleteExpense(ctx, branchID, id); err != nil {
		return apperror.New(http.StatusInternalServerError, "Delete expense failed", err.Error())
	}
	return nil
}
