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
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		SaleCombo(*fiber.Ctx) error
		PurchaseCombo(*fiber.Ctx) error
		OpnameCombo(*fiber.Ctx) error
	}
	Branch interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
	}
	UserBranch interface {
		List(*fiber.Ctx) error
		GetByKeys(*fiber.Ctx) error
		Create(*fiber.Ctx) error
	}
	User interface {
		List(*fiber.Ctx) error
		Detail(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
	}
	Supplier interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		Combo(*fiber.Ctx) error
	}
	Unit interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		Combo(*fiber.Ctx) error
	}
	ProductCategory interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		Combo(*fiber.Ctx) error
	}
	SupplierCategory interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		Combo(*fiber.Ctx) error
	}
	MemberCategory interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		Combo(*fiber.Ctx) error
	}
	Purchase interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
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
	app.Post("/api/branches", deps.AuthMiddleware, deps.Branch.Create)
	app.Delete("/api/branches/:id", deps.AuthMiddleware, deps.Branch.Delete)
	app.Get("/api/user-branches", deps.AuthMiddleware, deps.UserBranch.List)
	app.Get("/api/user-branches/:user_id/:branch_id", deps.AuthMiddleware, deps.UserBranch.GetByKeys)
	app.Post("/api/user-branches", deps.AuthMiddleware, deps.UserBranch.Create)
	app.Get("/api/users", deps.AuthMiddleware, deps.User.List)
	app.Get("/api/detail-users/:id", deps.AuthMiddleware, deps.User.Detail)
	app.Post("/api/users", deps.AuthMiddleware, deps.User.Create)
	app.Put("/api/users/:id", deps.AuthMiddleware, deps.User.Update)
	app.Get("/api/products", deps.AuthMiddleware, deps.Product.List)
	app.Post("/api/products", deps.AuthMiddleware, deps.Product.Create)
	app.Get("/api/products/:id", deps.AuthMiddleware, deps.Product.GetByID)
	app.Put("/api/products/:id", deps.AuthMiddleware, deps.Product.Update)
	app.Delete("/api/products/:id", deps.AuthMiddleware, deps.Product.Delete)
	app.Get("/api/sales-products-combo", deps.AuthMiddleware, deps.Product.SaleCombo)
	app.Get("/api/purchase-products-combo", deps.AuthMiddleware, deps.Product.PurchaseCombo)
	app.Get("/api/cmb-product-opname", deps.AuthMiddleware, deps.Product.OpnameCombo)
	app.Get("/api/suppliers", deps.AuthMiddleware, deps.Supplier.List)
	app.Get("/api/suppliers/:id", deps.AuthMiddleware, deps.Supplier.GetByID)
	app.Post("/api/suppliers", deps.AuthMiddleware, deps.Supplier.Create)
	app.Put("/api/suppliers/:id", deps.AuthMiddleware, deps.Supplier.Update)
	app.Delete("/api/suppliers/:id", deps.AuthMiddleware, deps.Supplier.Delete)
	app.Get("/api/suppliers-combo", deps.AuthMiddleware, deps.Supplier.Combo)
	app.Get("/api/units", deps.AuthMiddleware, deps.Unit.List)
	app.Get("/api/units/:id", deps.AuthMiddleware, deps.Unit.GetByID)
	app.Post("/api/units", deps.AuthMiddleware, deps.Unit.Create)
	app.Put("/api/units/:id", deps.AuthMiddleware, deps.Unit.Update)
	app.Delete("/api/units/:id", deps.AuthMiddleware, deps.Unit.Delete)
	app.Get("/api/cmb-units", deps.AuthMiddleware, deps.Unit.Combo)
	app.Get("/api/product-categories", deps.AuthMiddleware, deps.ProductCategory.List)
	app.Post("/api/product-categories", deps.AuthMiddleware, deps.ProductCategory.Create)
	app.Get("/api/product-categories/:id", deps.AuthMiddleware, deps.ProductCategory.GetByID)
	app.Put("/api/product-categories/:id", deps.AuthMiddleware, deps.ProductCategory.Update)
	app.Delete("/api/product-categories/:id", deps.AuthMiddleware, deps.ProductCategory.Delete)
	app.Get("/api/product-categories-combo", deps.AuthMiddleware, deps.ProductCategory.Combo)
	app.Get("/api/supplier-categories", deps.AuthMiddleware, deps.SupplierCategory.List)
	app.Post("/api/supplier-categories", deps.AuthMiddleware, deps.SupplierCategory.Create)
	app.Get("/api/supplier-categories/:id", deps.AuthMiddleware, deps.SupplierCategory.GetByID)
	app.Put("/api/supplier-categories/:id", deps.AuthMiddleware, deps.SupplierCategory.Update)
	app.Delete("/api/supplier-categories/:id", deps.AuthMiddleware, deps.SupplierCategory.Delete)
	app.Get("/api/supplier-categories-combo", deps.AuthMiddleware, deps.SupplierCategory.Combo)
	app.Get("/api/member-categories", deps.AuthMiddleware, deps.MemberCategory.List)
	app.Get("/api/member-categories/:id", deps.AuthMiddleware, deps.MemberCategory.GetByID)
	app.Post("/api/member-categories", deps.AuthMiddleware, deps.MemberCategory.Create)
	app.Put("/api/member-categories/:id", deps.AuthMiddleware, deps.MemberCategory.Update)
	app.Delete("/api/member-categories/:id", deps.AuthMiddleware, deps.MemberCategory.Delete)
	app.Get("/api/member-categories-combo", deps.AuthMiddleware, deps.MemberCategory.Combo)
	app.Get("/api/purchases", deps.AuthMiddleware, deps.Purchase.List)
	app.Get("/api/purchases/:id", deps.AuthMiddleware, deps.Purchase.GetByID)
	app.Post("/api/purchases", deps.AuthMiddleware, deps.Purchase.Create)
	app.Put("/api/purchases/:id", deps.AuthMiddleware, deps.Purchase.Update)
	app.Delete("/api/purchases/:id", deps.AuthMiddleware, deps.Purchase.Delete)
	app.Post("/api/sales", deps.AuthMiddleware, deps.Sale.Create)
	app.Post("/api/opnames", deps.AuthMiddleware, deps.Opname.Create)
	app.Get("/api/opnames/:id", deps.AuthMiddleware, deps.Opname.GetByID)
	app.Post("/api/opname-items", deps.AuthMiddleware, deps.Opname.CreateItem)
	app.Post("/api/opname-items-all", deps.AuthMiddleware, deps.Opname.GetItems)
}
