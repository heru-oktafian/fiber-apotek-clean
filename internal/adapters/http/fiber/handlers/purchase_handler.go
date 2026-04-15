package handlers

import (
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
