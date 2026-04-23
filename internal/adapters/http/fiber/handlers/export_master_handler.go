package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/presenter"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/domain/auth"
	memberusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/member"
	membercategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/membercategory"
	productcategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/productcategory"
	supplierusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/supplier"
	suppliercategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/suppliercategory"
	unitusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/unit"
)

type ExportMasterHandler struct {
	ProductCategories  productcategoryusecase.Service
	Suppliers          supplierusecase.Service
	SupplierCategories suppliercategoryusecase.Service
	MemberCategories   membercategoryusecase.Service
	Members            memberusecase.Service
	UnitConversions    unitusecase.ConversionService
}

func (h ExportMasterHandler) ProductCategoriesExcel(c *fiber.Ctx) error { return h.send(c, func(branchID string) ([]byte, string, error) { return h.ProductCategories.ExportExcel(c.Context(), branchID) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportMasterHandler) ProductCategoriesPDF(c *fiber.Ctx) error   { return h.send(c, func(branchID string) ([]byte, string, error) { return h.ProductCategories.ExportPDF(c.Context(), branchID) }, "application/pdf") }
func (h ExportMasterHandler) SuppliersExcel(c *fiber.Ctx) error         { return h.send(c, func(branchID string) ([]byte, string, error) { return h.Suppliers.ExportExcel(c.Context(), branchID) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") }
func (h ExportMasterHandler) SuppliersPDF(c *fiber.Ctx) error           { return h.send(c, func(branchID string) ([]byte, string, error) { return h.Suppliers.ExportPDF(c.Context(), branchID) }, "application/pdf") }
func (h ExportMasterHandler) SupplierCategoriesExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.SupplierCategories.ExportExcel(c.Context(), branchID) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}
func (h ExportMasterHandler) SupplierCategoriesPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.SupplierCategories.ExportPDF(c.Context(), branchID) }, "application/pdf")
}
func (h ExportMasterHandler) MemberCategoriesExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.MemberCategories.ExportExcel(c.Context(), branchID) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}
func (h ExportMasterHandler) MemberCategoriesPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.MemberCategories.ExportPDF(c.Context(), branchID) }, "application/pdf")
}
func (h ExportMasterHandler) MembersExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.Members.ExportExcel(c.Context(), branchID) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}
func (h ExportMasterHandler) MembersPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.Members.ExportPDF(c.Context(), branchID) }, "application/pdf")
}
func (h ExportMasterHandler) UnitConversionsExcel(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.UnitConversions.ExportExcel(c.Context(), branchID) }, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
}
func (h ExportMasterHandler) UnitConversionsPDF(c *fiber.Ctx) error {
	return h.send(c, func(branchID string) ([]byte, string, error) { return h.UnitConversions.ExportPDF(c.Context(), branchID) }, "application/pdf")
}

func (h ExportMasterHandler) send(c *fiber.Ctx, fn func(branchID string) ([]byte, string, error), contentType string) error {
	claims := c.Locals("claims").(auth.Claims)
	data, filename, err := fn(claims.BranchID)
	if err != nil {
		return presenter.Handle(c, err)
	}
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(data)
}
