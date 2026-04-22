package router

import "github.com/gofiber/fiber/v2"

type Dependencies struct {
	Auth interface {
		Login(*fiber.Ctx) error
		ListBranches(*fiber.Ctx) error
		Menus(*fiber.Ctx) error
		Logout(*fiber.Ctx) error
		SetBranch(*fiber.Ctx) error
		Profile(*fiber.Ctx) error
	}
	Product interface {
		Create(*fiber.Ctx) error
		SaleCombo(*fiber.Ctx) error
		PurchaseCombo(*fiber.Ctx) error
		OpnameCombo(*fiber.Ctx) error
	}
	Branch interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
	}
	UserBranch interface {
		List(*fiber.Ctx) error
		GetByKeys(*fiber.Ctx) error
	}
	User interface {
		List(*fiber.Ctx) error
		Detail(*fiber.Ctx) error
	}
	Purchase interface {
		Create(*fiber.Ctx) error
	}
	Sale interface {
		Create(*fiber.Ctx) error
	}
	Opname interface {
		Create(*fiber.Ctx) error
		CreateItem(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		GetItems(*fiber.Ctx) error
	}
	AuthMiddleware fiber.Handler
}

func Register(app *fiber.App, deps Dependencies) {
	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "ok"}) })
	app.Post("/api/login", deps.Auth.Login)
	app.Get("/api/list_branches", deps.AuthMiddleware, deps.Auth.ListBranches)
	app.Get("/api/menus", deps.AuthMiddleware, deps.Auth.Menus)
	app.Post("/api/set_branch", deps.AuthMiddleware, deps.Auth.SetBranch)
	app.Get("/api/profile", deps.AuthMiddleware, deps.Auth.Profile)
	app.Post("/api/logout", deps.AuthMiddleware, deps.Auth.Logout)
	app.Get("/api/branches", deps.AuthMiddleware, deps.Branch.List)
	app.Get("/api/branches/:id", deps.AuthMiddleware, deps.Branch.GetByID)
	app.Get("/api/user-branches", deps.AuthMiddleware, deps.UserBranch.List)
	app.Get("/api/user-branches/:user_id/:branch_id", deps.AuthMiddleware, deps.UserBranch.GetByKeys)
	app.Get("/api/users", deps.AuthMiddleware, deps.User.List)
	app.Get("/api/detail-users/:id", deps.AuthMiddleware, deps.User.Detail)
	app.Post("/api/products", deps.AuthMiddleware, deps.Product.Create)
	app.Get("/api/sales-products-combo", deps.AuthMiddleware, deps.Product.SaleCombo)
	app.Get("/api/purchase-products-combo", deps.AuthMiddleware, deps.Product.PurchaseCombo)
	app.Get("/api/cmb-product-opname", deps.AuthMiddleware, deps.Product.OpnameCombo)
	app.Post("/api/purchases", deps.AuthMiddleware, deps.Purchase.Create)
	app.Post("/api/sales", deps.AuthMiddleware, deps.Sale.Create)
	app.Post("/api/opnames", deps.AuthMiddleware, deps.Opname.Create)
	app.Get("/api/opnames/:id", deps.AuthMiddleware, deps.Opname.GetByID)
	app.Post("/api/opname-items", deps.AuthMiddleware, deps.Opname.CreateItem)
	app.Post("/api/opname-items-all", deps.AuthMiddleware, deps.Opname.GetItems)
}
