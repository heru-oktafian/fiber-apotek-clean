package handlers

import (
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
