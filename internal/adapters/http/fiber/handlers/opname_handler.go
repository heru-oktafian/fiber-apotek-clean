package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/opname"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	opnameusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/opname"
)

type OpnameHandler struct {
	Service opnameusecase.Service
}

func (h OpnameHandler) Create(c *fiber.Ctx) error {
	var req opname.CreateOpnameRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}

	claims := c.Locals("claims").(auth.Claims)
	entity, err := h.Service.CreateHeader(c.Context(), claims.BranchID, claims.Subject, req)
	if err != nil {
		return presenter.Handle(c, err)
	}

	return response.JSON(c, fiber.StatusOK, "Opname created successfully", entity)
}

func (h OpnameHandler) CreateItem(c *fiber.Ctx) error {
	var req opname.CreateOpnameItemRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}

	item, err := h.Service.CreateItem(c.Context(), req)
	if err != nil {
		return presenter.Handle(c, err)
	}

	return response.JSON(c, fiber.StatusOK, "Opname item created successfully", item)
}

func (h OpnameHandler) GetByID(c *fiber.Ctx) error {
	detail, err := h.Service.GetDetail(c.Context(), c.Params("id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Opname retrieved successfully", detail)
}

func (h OpnameHandler) GetItems(c *fiber.Ctx) error {
	var payload struct {
		OpnameID string `json:"opname_id"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return presenter.Handle(c, err)
	}
	items, err := h.Service.GetItems(c.Context(), payload.OpnameID)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Opname items retrieved successfully", items)
}
