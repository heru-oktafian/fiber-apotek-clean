package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/expense"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	expenseusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/expense"
)

type ExpenseHandler struct {
	Service expenseusecase.Service
}

func (h ExpenseHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, expense.ListRequest{
		Search: c.Query("search"),
		Month:  c.Query("month"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Expenses retrieved successfully", result.Items, result.Meta)
}

func (h ExpenseHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req expense.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), claims.BranchID, claims.Subject, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Expense created successfully", item)
}

func (h ExpenseHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req expense.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Expense updated successfully", item)
}

func (h ExpenseHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.Delete(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Expense deleted successfully", fiber.Map{"id": c.Params("id")})
}
