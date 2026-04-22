package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/sale"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	saleusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/sale"
)

type SaleHandler struct {
	Service saleusecase.Service
}

func (h SaleHandler) List(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.Service.List(c.Context(), claims.BranchID, sale.ListRequest{Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSONWithMeta(c, fiber.StatusOK, "Sales retrieved successfully", result.Items, result.Meta)
}

func (h SaleHandler) GetByID(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	item, err := h.Service.GetByID(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Sale retrieved successfully", item)
}

func (h SaleHandler) Create(c *fiber.Ctx) error {
	var req sale.CreateSaleRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}

	claims := c.Locals("claims").(auth.Claims)
	saleEntity, items, err := h.Service.CreateTransaction(
		c.Context(),
		claims.BranchID,
		claims.Subject,
		claims.DefaultMember,
		claims.SubscriptionType,
		req,
	)
	if err != nil {
		return presenter.Handle(c, err)
	}

	return response.JSON(c, fiber.StatusOK, "Sale created successfully", fiber.Map{
		"sale":       saleEntity,
		"sale_items": items,
	})
}

func (h SaleHandler) Update(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req sale.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.Update(c.Context(), claims.BranchID, c.Params("id"), claims.DefaultMember, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Sale updated successfully", item)
}

func (h SaleHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.Delete(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Sale deleted successfully", fiber.Map{"id": c.Params("id")})
}

func (h SaleHandler) ListItems(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	items, err := h.Service.ListItems(c.Context(), claims.BranchID, c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Items retrieved successfully", items)
}

func (h SaleHandler) CreateItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req sale.CreateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.CreateItem(c.Context(), claims.BranchID, req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Item added successfully", item)
}

func (h SaleHandler) UpdateItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	var req sale.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	item, err := h.Service.UpdateItem(c.Context(), claims.BranchID, c.Params("id"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Item updated successfully", item)
}

func (h SaleHandler) DeleteItem(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	if err := h.Service.DeleteItem(c.Context(), claims.BranchID, c.Params("id")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Item deleted successfully", fiber.Map{"id": c.Params("id")})
}
