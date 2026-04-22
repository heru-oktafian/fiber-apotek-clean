package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	productusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/product"
	unitusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/unit"
)

type ExportHandler struct {
	Products productusecase.Service
	Units    unitusecase.MasterService
}

func (h ExportHandler) ProductsExcel(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	data, filename, err := h.Products.ExportExcel(c.Context(), claims.BranchID)
	if err != nil {
		return presenter.Handle(c, err)
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(data)
}

func (h ExportHandler) ProductsPDF(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	data, filename, err := h.Products.ExportPDF(c.Context(), claims.BranchID)
	if err != nil {
		return presenter.Handle(c, err)
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(data)
}

func (h ExportHandler) UnitsExcel(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	data, filename, err := h.Units.ExportExcel(c.Context(), claims.BranchID)
	if err != nil {
		return presenter.Handle(c, err)
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(data)
}

func (h ExportHandler) UnitsPDF(c *fiber.Ctx) error {
	claims := c.Locals("claims").(auth.Claims)
	data, filename, err := h.Units.ExportPDF(c.Context(), claims.BranchID)
	if err != nil {
		return presenter.Handle(c, err)
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(data)
}

func exportTimestampedFilename(base, ext string) string {
	return fmt.Sprintf("%s-%s.%s", base, time.Now().Format("2006-01-02-15-04-05"), ext)
}
