package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/duplicatereceipt"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	duplicatereceiptusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/duplicatereceipt"
)

type DuplicateReceiptHandler struct {
	Service duplicatereceiptusecase.Service
}

func (h DuplicateReceiptHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	result, err := h.Service.List(c.Context(), claims.BranchID, duplicatereceipt.ListRequest{Search: c.Query("search"), Month: c.Query("month"), Page: c.QueryInt("page", 1), Limit: c.QueryInt("limit", 10)})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Duplicate receipts retrieved successfully", result.Items, result.Meta)
}

func (h DuplicateReceiptHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	item, err := h.Service.GetByID(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt retrieved successfully", item)
}

func (h DuplicateReceiptHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req duplicatereceipt.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Create(c.Context(), claims.BranchID, claims.Subject, claims.DefaultMember, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt transaction created successfully", item)
}

func (h DuplicateReceiptHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req duplicatereceipt.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, c.Params("id"), claims.DefaultMember, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt updated successfully", item)
}

func (h DuplicateReceiptHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.Delete(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt deleted successfully", fiber.Map{"id": c.Params("id")})
}

func (h DuplicateReceiptHandler) ListItems(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.ListItems(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt items retrieved successfully", items)
}

func (h DuplicateReceiptHandler) CreateItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req duplicatereceipt.CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.CreateItem(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt item created successfully", item)
}

func (h DuplicateReceiptHandler) UpdateItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req duplicatereceipt.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.UpdateItem(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt item updated successfully", item)
}

func (h DuplicateReceiptHandler) DeleteItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.DeleteItem(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Duplicate receipt item deleted successfully", fiber.Map{"id": c.Params("id")})
}
