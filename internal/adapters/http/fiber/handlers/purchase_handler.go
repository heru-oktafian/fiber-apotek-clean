package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/purchase"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	purchaseusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/purchase"
)

type PurchaseHandler struct {
	Service purchaseusecase.Service
}

func (h PurchaseHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, purchase.ListRequest{Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Purchases retrieved successfully", result.Items, result.Meta)
}

func (h PurchaseHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	item, err := h.Service.GetByID(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Purchase retrieved successfully", item)
}

func (h PurchaseHandler) Create(c *fiber.Ctx) error {
	var req purchase.CreatePurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}

	claims := c.Locals("claims").(auth.Claims)
	purchaseEntity, items, err := h.Service.CreateTransaction(c.Context(), claims.BranchID, claims.Subject, req)
	if err != nil {
		return presenter.Handle(c, err)
	}

	return response.JSON(c, fiber.StatusOK, "Purchase created successfully", fiber.Map{
		"purchase":       purchaseEntity,
		"purchase_items": items,
	})
}

func (h PurchaseHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req purchase.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Purchase updated successfully", item)
}

func (h PurchaseHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.Delete(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Purchase deleted successfully", fiber.Map{"id": c.Params("id")})
}
