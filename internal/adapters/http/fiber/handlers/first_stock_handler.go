package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/firststock"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	firststockusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/firststock"
)

type FirstStockHandler struct {
	Service firststockusecase.Service
}

func (h FirstStockHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, firststock.ListRequest{Search: c.Query("search"), Month: c.Query("month"), Page: page, Limit: limit})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "First stocks retrieved successfully", result.Items, result.Meta)
}

func (h FirstStockHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req firststock.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), claims.BranchID, claims.Subject, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "First stock created successfully", item)
}

func (h FirstStockHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req firststock.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "First stock updated successfully", item)
}

func (h FirstStockHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.Delete(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "First stock deleted successfully", fiber.Map{"id": c.Params("id")})
}

func (h FirstStockHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	item, err := h.Service.GetDetail(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "First stock retrieved successfully", item)
}

func (h FirstStockHandler) ListItems(c *fiber.Ctx) error {
	items, err := h.Service.ListItems(c.Context(), c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Items retrieved successfully", items)
}

func (h FirstStockHandler) CreateItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req firststock.CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.CreateItem(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Item added successfully", item)
}

func (h FirstStockHandler) UpdateItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req firststock.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.UpdateItem(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Item updated successfully", item)
}

func (h FirstStockHandler) DeleteItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.DeleteItem(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Item deleted successfully", fiber.Map{"id": c.Params("id")})
}
