package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/response"
	authusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/auth"
)

type AuthHandler struct {
	Service authusecase.Service
}

func (h AuthHandler) ListBranches(c *fiber.Ctx) error {
	items, err := h.Service.ListBranches(c.Context(), c.Get("Authorization"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "User Branch found", items)
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	result, err := h.Service.Login(c.Context(), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Login successful", result.Token)
}

func (h AuthHandler) Menus(c *fiber.Ctx) error {
	items, message, err := h.Service.Menus(c.Context(), c.Get("Authorization"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, message, items)
}

func (h AuthHandler) Logout(c *fiber.Ctx) error {
	if err := h.Service.Logout(c.Context(), c.Get("Authorization")); err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Logout successful", "Logout successful")
}

func (h AuthHandler) SetBranch(c *fiber.Ctx) error {
	var req auth.BranchSelectionRequest
	if err := c.BodyParser(&req); err != nil {
		return presenter.Handle(c, err)
	}
	token, err := h.Service.SetBranch(c.Context(), c.Get("Authorization"), req)
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Branch set successfully", token)
}

func (h AuthHandler) Profile(c *fiber.Ctx) error {
	item, err := h.Service.Profile(c.Context(), c.Get("Authorization"))
	if err != nil {
		return presenter.Handle(c, err)
	}
	claims, _, err := h.Service.Tokens.Parse(strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer ")))
	if err != nil {
		return presenter.Handle(c, err)
	}
	return response.JSON(c, fiber.StatusOK, "Otoritas : "+claims.UserRole, item)
}
