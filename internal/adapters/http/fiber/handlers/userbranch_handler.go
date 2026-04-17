package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	userbranchusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/userbranch"
)

type UserBranchHandler struct {
	Service userbranchusecase.Service
}

func (h UserBranchHandler) List(c *fiber.Ctx) error {
	items, err := h.Service.List(c.Context())
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "UserBranches retrieved successfully", items)
}

func (h UserBranchHandler) GetByKeys(c *fiber.Ctx) error {
	items, err := h.Service.GetByKeys(c.Context(), c.Params("user_id"), c.Params("branch_id"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "UserBranch found", items)
}
