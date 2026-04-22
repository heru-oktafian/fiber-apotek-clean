package expense

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/common"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/expense"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/ports"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/apperror"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
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

func (s Service) ExportExcel(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, expense.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	f := exportshared.NewExcelFile("Expenses")
	sheet := "Expenses"
	f.SetCellValue(sheet, "A1", fmt.Sprintf("PENGELUARAN %s", month))
	headers := []string{"ID", "KETERANGAN", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		col, _ := exportshared.ExcelColumnName(i + 1)
		f.SetCellValue(sheet, fmt.Sprintf("%s3", col), h)
	}
	grandTotal := 0
	for i, item := range result.Items {
		row := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Description)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.ExpenseDate)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Payment)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.TotalExpense)
		grandTotal += item.TotalExpense
	}
	totalRow := len(result.Items) + 4
	f.SetCellValue(sheet, fmt.Sprintf("A%d", totalRow), "GRAND TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", totalRow), grandTotal)
	bytes, err := exportshared.WriteExcel(f)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export expenses excel failed", err.Error())
	}
	return bytes, fmt.Sprintf("expenses-%s.xlsx", time.Now().Format("2006-01-02-15-04-05")), nil
}

func (s Service) ExportPDF(ctx context.Context, branchID, month string) ([]byte, string, error) {
	result, err := s.List(ctx, branchID, expense.ListRequest{Month: month, Page: 1, Limit: 10000})
	if err != nil {
		return nil, "", err
	}
	pdf := exportshared.NewPDF("PENGELUARAN")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(277, 10, fmt.Sprintf("PENGELUARAN %s", month), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	widths := []float64{45, 110, 40, 40, 42}
	headers := []string{"ID", "KETERANGAN", "TANGGAL", "PEMBAYARAN", "TOTAL"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	grandTotal := 0
	for _, item := range result.Items {
		values := []string{item.ID, item.Description, item.ExpenseDate, item.Payment, fmt.Sprintf("%d", item.TotalExpense)}
		for i, v := range values {
			pdf.CellFormat(widths[i], 8, v, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
		grandTotal += item.TotalExpense
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(235, 8, "TOTAL", "1", 0, "C", false, 0, "")
	pdf.CellFormat(42, 8, fmt.Sprintf("%d", grandTotal), "1", 1, "R", false, 0, "")
	bytes, err := exportshared.WritePDF(pdf)
	if err != nil {
		return nil, "", apperror.New(http.StatusInternalServerError, "Export expenses pdf failed", err.Error())
	}
	return bytes, fmt.Sprintf("PENGELUARAN-%s.pdf", time.Now().Format("2006-01-02-15-04-05")), nil
}
