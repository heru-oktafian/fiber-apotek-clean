package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	exportshared "github.com/heru-oktafian/fiber-apotek-clean/internal/shared/export"
	anotherincomeusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/anotherincome"
	expenseusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/expense"
)

type ExportFinanceHandler struct {
	AnotherIncomes anotherincomeusecase.Service
	Expenses       expenseusecase.Service
}

func (h ExportFinanceHandler) AnotherIncomesExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.AnotherIncomes.ExportExcel(c.Context(), branchID, f.Month)
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportFinanceHandler) AnotherIncomesPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.AnotherIncomes.ExportPDF(c.Context(), branchID, f.Month)
	}, "application/pdf")
}

func (h ExportFinanceHandler) ExpensesExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.Expenses.ExportExcel(c.Context(), branchID, f.Month)
	}, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}

func (h ExportFinanceHandler) ExpensesPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string, f exportshared.Filters) ([]byte, string, error) {
		return h.Expenses.ExportPDF(c.Context(), branchID, f.Month)
	}, "application/pdf")
}

func (h ExportFinanceHandler) send(c *fiber.Ctx, fn func(branchID string, filters exportshared.Filters) ([]byte, string, error), contentType string) error {
	claims := c.Locals("claims").(auth.Claims)
	filters := exportshared.ParseFilters(c)
	data, filename, err := fn(claims.BranchID, filters)
	if err != nil {
		return presenter.Handle(c, err)
	}
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(data)
}
