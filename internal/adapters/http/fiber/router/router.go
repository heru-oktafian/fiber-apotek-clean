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
	UnitConversion interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
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
	Member interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		Combo(*fiber.Ctx) error
	}
	AnotherIncome interface {
		List(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
	}
	Expense interface {
		List(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
	}
	FirstStock interface {
		List(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		ListItems(*fiber.Ctx) error
		CreateItem(*fiber.Ctx) error
		UpdateItem(*fiber.Ctx) error
		DeleteItem(*fiber.Ctx) error
	}
	BuyReturn interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		PurchaseSources(*fiber.Ctx) error
		ReturnableItems(*fiber.Ctx) error
	}
	SaleReturn interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		SaleSources(*fiber.Ctx) error
		ReturnableItems(*fiber.Ctx) error
	}
	DuplicateReceipt interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		ListDetailSummaries(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		ListItems(*fiber.Ctx) error
		CreateItem(*fiber.Ctx) error
		UpdateItem(*fiber.Ctx) error
		DeleteItem(*fiber.Ctx) error
	}
	Purchase interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		ListItems(*fiber.Ctx) error
		CreateItem(*fiber.Ctx) error
		UpdateItem(*fiber.Ctx) error
		DeleteItem(*fiber.Ctx) error
	}
	Sale interface {
		List(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		Create(*fiber.Ctx) error
		Update(*fiber.Ctx) error
		Delete(*fiber.Ctx) error
		ListItems(*fiber.Ctx) error
		CreateItem(*fiber.Ctx) error
		UpdateItem(*fiber.Ctx) error
		DeleteItem(*fiber.Ctx) error
	}
	Opname interface {
		Create(*fiber.Ctx) error
		CreateItem(*fiber.Ctx) error
		GetByID(*fiber.Ctx) error
		GetItems(*fiber.Ctx) error
	}
	Export interface {
		ProductsExcel(*fiber.Ctx) error
		ProductsPDF(*fiber.Ctx) error
		UnitsExcel(*fiber.Ctx) error
		UnitsPDF(*fiber.Ctx) error
		ProductCategoriesExcel(*fiber.Ctx) error
		ProductCategoriesPDF(*fiber.Ctx) error
		SuppliersExcel(*fiber.Ctx) error
		SuppliersPDF(*fiber.Ctx) error
		SupplierCategoriesExcel(*fiber.Ctx) error
		SupplierCategoriesPDF(*fiber.Ctx) error
		MemberCategoriesExcel(*fiber.Ctx) error
		MemberCategoriesPDF(*fiber.Ctx) error
		MembersExcel(*fiber.Ctx) error
		MembersPDF(*fiber.Ctx) error
		UnitConversionsExcel(*fiber.Ctx) error
		UnitConversionsPDF(*fiber.Ctx) error
		PurchasesExcel(*fiber.Ctx) error
		PurchasesPDF(*fiber.Ctx) error
		PurchaseItemsExcel(*fiber.Ctx) error
		PurchaseItemsPDF(*fiber.Ctx) error
		SalesExcel(*fiber.Ctx) error
		SalesPDF(*fiber.Ctx) error
		SaleItemsExcel(*fiber.Ctx) error
		SaleItemsPDF(*fiber.Ctx) error
		DuplicateReceiptsExcel(*fiber.Ctx) error
		DuplicateReceiptsPDF(*fiber.Ctx) error
		DuplicateReceiptItemsExcel(*fiber.Ctx) error
		DuplicateReceiptItemsPDF(*fiber.Ctx) error
		AnotherIncomesExcel(*fiber.Ctx) error
		AnotherIncomesPDF(*fiber.Ctx) error
		ExpensesExcel(*fiber.Ctx) error
		ExpensesPDF(*fiber.Ctx) error
		FirstStocksExcel(*fiber.Ctx) error
		FirstStocksPDF(*fiber.Ctx) error
		FirstStockItemsExcel(*fiber.Ctx) error
		FirstStockItemsPDF(*fiber.Ctx) error
		BuyReturnsExcel(*fiber.Ctx) error
		BuyReturnsPDF(*fiber.Ctx) error
		BuyReturnItemsExcel(*fiber.Ctx) error
		BuyReturnItemsPDF(*fiber.Ctx) error
		SaleReturnsExcel(*fiber.Ctx) error
		SaleReturnsPDF(*fiber.Ctx) error
		SaleReturnItemsExcel(*fiber.Ctx) error
		SaleReturnItemsPDF(*fiber.Ctx) error
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
	app.Get("/api/products/excel", deps.AuthMiddleware, deps.Export.ProductsExcel)
	app.Get("/api/products/pdf", deps.AuthMiddleware, deps.Export.ProductsPDF)
	app.Post("/api/products", deps.AuthMiddleware, deps.Product.Create)
	app.Get("/api/products/:id", deps.AuthMiddleware, deps.Product.GetByID)
	app.Put("/api/products/:id", deps.AuthMiddleware, deps.Product.Update)
	app.Delete("/api/products/:id", deps.AuthMiddleware, deps.Product.Delete)
	app.Get("/api/sales-products-combo", deps.AuthMiddleware, deps.Product.SaleCombo)
	app.Get("/api/purchase-products-combo", deps.AuthMiddleware, deps.Product.PurchaseCombo)
	app.Get("/api/cmb-product-opname", deps.AuthMiddleware, deps.Product.OpnameCombo)
	app.Get("/api/suppliers", deps.AuthMiddleware, deps.Supplier.List)
	app.Get("/api/suppliers/excel", deps.AuthMiddleware, deps.Export.SuppliersExcel)
	app.Get("/api/suppliers/pdf", deps.AuthMiddleware, deps.Export.SuppliersPDF)
	app.Get("/api/suppliers/:id", deps.AuthMiddleware, deps.Supplier.GetByID)
	app.Post("/api/suppliers", deps.AuthMiddleware, deps.Supplier.Create)
	app.Put("/api/suppliers/:id", deps.AuthMiddleware, deps.Supplier.Update)
	app.Delete("/api/suppliers/:id", deps.AuthMiddleware, deps.Supplier.Delete)
	app.Get("/api/suppliers-combo", deps.AuthMiddleware, deps.Supplier.Combo)
	app.Get("/api/units", deps.AuthMiddleware, deps.Unit.List)
	app.Get("/api/units/excel", deps.AuthMiddleware, deps.Export.UnitsExcel)
	app.Get("/api/units/pdf", deps.AuthMiddleware, deps.Export.UnitsPDF)
	app.Get("/api/units/:id", deps.AuthMiddleware, deps.Unit.GetByID)
	app.Post("/api/units", deps.AuthMiddleware, deps.Unit.Create)
	app.Put("/api/units/:id", deps.AuthMiddleware, deps.Unit.Update)
	app.Delete("/api/units/:id", deps.AuthMiddleware, deps.Unit.Delete)
	app.Get("/api/cmb-units", deps.AuthMiddleware, deps.Unit.Combo)
	app.Get("/api/unit-conversions", deps.AuthMiddleware, deps.UnitConversion.List)
	app.Get("/api/unit-conversions/excel", deps.AuthMiddleware, deps.Export.UnitConversionsExcel)
	app.Get("/api/unit-conversions/pdf", deps.AuthMiddleware, deps.Export.UnitConversionsPDF)
	app.Get("/api/unit-conversions/:id", deps.AuthMiddleware, deps.UnitConversion.GetByID)
	app.Post("/api/unit-conversions", deps.AuthMiddleware, deps.UnitConversion.Create)
	app.Put("/api/unit-conversions/:id", deps.AuthMiddleware, deps.UnitConversion.Update)
	app.Delete("/api/unit-conversions/:id", deps.AuthMiddleware, deps.UnitConversion.Delete)
	app.Get("/api/product-categories", deps.AuthMiddleware, deps.ProductCategory.List)
	app.Get("/api/product-categories/excel", deps.AuthMiddleware, deps.Export.ProductCategoriesExcel)
	app.Get("/api/product-categories/pdf", deps.AuthMiddleware, deps.Export.ProductCategoriesPDF)
	app.Post("/api/product-categories", deps.AuthMiddleware, deps.ProductCategory.Create)
	app.Get("/api/product-categories/:id", deps.AuthMiddleware, deps.ProductCategory.GetByID)
	app.Put("/api/product-categories/:id", deps.AuthMiddleware, deps.ProductCategory.Update)
	app.Delete("/api/product-categories/:id", deps.AuthMiddleware, deps.ProductCategory.Delete)
	app.Get("/api/product-categories-combo", deps.AuthMiddleware, deps.ProductCategory.Combo)
	app.Get("/api/supplier-categories", deps.AuthMiddleware, deps.SupplierCategory.List)
	app.Get("/api/supplier-categories/excel", deps.AuthMiddleware, deps.Export.SupplierCategoriesExcel)
	app.Get("/api/supplier-categories/pdf", deps.AuthMiddleware, deps.Export.SupplierCategoriesPDF)
	app.Post("/api/supplier-categories", deps.AuthMiddleware, deps.SupplierCategory.Create)
	app.Get("/api/supplier-categories/:id", deps.AuthMiddleware, deps.SupplierCategory.GetByID)
	app.Put("/api/supplier-categories/:id", deps.AuthMiddleware, deps.SupplierCategory.Update)
	app.Delete("/api/supplier-categories/:id", deps.AuthMiddleware, deps.SupplierCategory.Delete)
	app.Get("/api/supplier-categories-combo", deps.AuthMiddleware, deps.SupplierCategory.Combo)
	app.Get("/api/member-categories", deps.AuthMiddleware, deps.MemberCategory.List)
	app.Get("/api/members", deps.AuthMiddleware, deps.Member.List)
	app.Get("/api/members/excel", deps.AuthMiddleware, deps.Export.MembersExcel)
	app.Get("/api/members/pdf", deps.AuthMiddleware, deps.Export.MembersPDF)
	app.Get("/api/members/:id", deps.AuthMiddleware, deps.Member.GetByID)
	app.Post("/api/members", deps.AuthMiddleware, deps.Member.Create)
	app.Put("/api/members/:id", deps.AuthMiddleware, deps.Member.Update)
	app.Delete("/api/members/:id", deps.AuthMiddleware, deps.Member.Delete)
	app.Get("/api/members-combo", deps.AuthMiddleware, deps.Member.Combo)
	app.Get("/api/another-incomes", deps.AuthMiddleware, deps.AnotherIncome.List)
	app.Get("/api/another-incomes/excel", deps.AuthMiddleware, deps.Export.AnotherIncomesExcel)
	app.Get("/api/another-incomes/pdf", deps.AuthMiddleware, deps.Export.AnotherIncomesPDF)
	app.Post("/api/another-incomes", deps.AuthMiddleware, deps.AnotherIncome.Create)
	app.Put("/api/another-incomes/:id", deps.AuthMiddleware, deps.AnotherIncome.Update)
	app.Delete("/api/another-incomes/:id", deps.AuthMiddleware, deps.AnotherIncome.Delete)
	app.Get("/api/expenses", deps.AuthMiddleware, deps.Expense.List)
	app.Get("/api/expenses/excel", deps.AuthMiddleware, deps.Export.ExpensesExcel)
	app.Get("/api/expenses/pdf", deps.AuthMiddleware, deps.Export.ExpensesPDF)
	app.Post("/api/expenses", deps.AuthMiddleware, deps.Expense.Create)
	app.Put("/api/expenses/:id", deps.AuthMiddleware, deps.Expense.Update)
	app.Delete("/api/expenses/:id", deps.AuthMiddleware, deps.Expense.Delete)
	app.Get("/api/first-stocks", deps.AuthMiddleware, deps.FirstStock.List)
	app.Get("/api/buy-returns", deps.AuthMiddleware, deps.BuyReturn.List)
	app.Get("/api/buy-returns/excel", deps.AuthMiddleware, deps.Export.BuyReturnsExcel)
	app.Get("/api/buy-returns/pdf", deps.AuthMiddleware, deps.Export.BuyReturnsPDF)
	app.Post("/api/buy-returns", deps.AuthMiddleware, deps.BuyReturn.Create)
	app.Get("/api/buy-returns/:id", deps.AuthMiddleware, deps.BuyReturn.GetByID)
	app.Get("/api/buy-return-items/excel", deps.AuthMiddleware, deps.Export.BuyReturnItemsExcel)
	app.Get("/api/buy-return-items/pdf", deps.AuthMiddleware, deps.Export.BuyReturnItemsPDF)
	app.Get("/api/cmb-purchases", deps.AuthMiddleware, deps.BuyReturn.PurchaseSources)
	app.Get("/api/cmb-prod-buy-returns", deps.AuthMiddleware, deps.BuyReturn.ReturnableItems)
	app.Get("/api/sale-returns", deps.AuthMiddleware, deps.SaleReturn.List)
	app.Get("/api/sale-returns/excel", deps.AuthMiddleware, deps.Export.SaleReturnsExcel)
	app.Get("/api/sale-returns/pdf", deps.AuthMiddleware, deps.Export.SaleReturnsPDF)
	app.Post("/api/sale-returns", deps.AuthMiddleware, deps.SaleReturn.Create)
	app.Get("/api/sale-returns/:id", deps.AuthMiddleware, deps.SaleReturn.GetByID)
	app.Get("/api/sale-return-items/excel", deps.AuthMiddleware, deps.Export.SaleReturnItemsExcel)
	app.Get("/api/sale-return-items/pdf", deps.AuthMiddleware, deps.Export.SaleReturnItemsPDF)
	app.Get("/api/cmb-sales", deps.AuthMiddleware, deps.SaleReturn.SaleSources)
	app.Get("/api/cmb-prod-sale-returns", deps.AuthMiddleware, deps.SaleReturn.ReturnableItems)
	app.Get("/api/duplicate-receipts", deps.AuthMiddleware, deps.DuplicateReceipt.List)
	app.Get("/api/duplicate-receipts-details", deps.AuthMiddleware, deps.DuplicateReceipt.ListDetailSummaries)
	app.Get("/api/duplicate-receipts/excel", deps.AuthMiddleware, deps.Export.DuplicateReceiptsExcel)
	app.Get("/api/duplicate-receipts/pdf", deps.AuthMiddleware, deps.Export.DuplicateReceiptsPDF)
	app.Post("/api/duplicate-receipts", deps.AuthMiddleware, deps.DuplicateReceipt.Create)
	app.Get("/api/duplicate-receipts/:id", deps.AuthMiddleware, deps.DuplicateReceipt.GetByID)
	app.Put("/api/duplicate-receipts/:id", deps.AuthMiddleware, deps.DuplicateReceipt.Update)
	app.Delete("/api/duplicate-receipts/:id", deps.AuthMiddleware, deps.DuplicateReceipt.Delete)
	app.Get("/api/duplicate-receipts-items/all/:id", deps.AuthMiddleware, deps.DuplicateReceipt.ListItems)
	app.Get("/api/duplicate-receipts-items/excel", deps.AuthMiddleware, deps.Export.DuplicateReceiptItemsExcel)
	app.Get("/api/duplicate-receipts-items/pdf", deps.AuthMiddleware, deps.Export.DuplicateReceiptItemsPDF)
	app.Post("/api/duplicate-receipts-items", deps.AuthMiddleware, deps.DuplicateReceipt.CreateItem)
	app.Put("/api/duplicate-receipts-items/:id", deps.AuthMiddleware, deps.DuplicateReceipt.UpdateItem)
	app.Delete("/api/duplicate-receipts-items/:id", deps.AuthMiddleware, deps.DuplicateReceipt.DeleteItem)
	app.Get("/api/first-stocks/excel", deps.AuthMiddleware, deps.Export.FirstStocksExcel)
	app.Get("/api/first-stocks/pdf", deps.AuthMiddleware, deps.Export.FirstStocksPDF)
	app.Post("/api/first-stocks", deps.AuthMiddleware, deps.FirstStock.Create)
	app.Put("/api/first-stocks/:id", deps.AuthMiddleware, deps.FirstStock.Update)
	app.Delete("/api/first-stocks/:id", deps.AuthMiddleware, deps.FirstStock.Delete)
	app.Get("/api/first-stock-with-items/:id", deps.AuthMiddleware, deps.FirstStock.GetByID)
	app.Get("/api/first-stock-items/:id", deps.AuthMiddleware, deps.FirstStock.ListItems)
	app.Get("/api/first-stock-items/excel", deps.AuthMiddleware, deps.Export.FirstStockItemsExcel)
	app.Get("/api/first-stock-items/pdf", deps.AuthMiddleware, deps.Export.FirstStockItemsPDF)
	app.Post("/api/first-stock-items", deps.AuthMiddleware, deps.FirstStock.CreateItem)
	app.Put("/api/first-stock-items/:id", deps.AuthMiddleware, deps.FirstStock.UpdateItem)
	app.Delete("/api/first-stock-items/:id", deps.AuthMiddleware, deps.FirstStock.DeleteItem)
	app.Get("/api/member-categories/excel", deps.AuthMiddleware, deps.Export.MemberCategoriesExcel)
	app.Get("/api/member-categories/pdf", deps.AuthMiddleware, deps.Export.MemberCategoriesPDF)
	app.Get("/api/member-categories/:id", deps.AuthMiddleware, deps.MemberCategory.GetByID)
	app.Post("/api/member-categories", deps.AuthMiddleware, deps.MemberCategory.Create)
	app.Put("/api/member-categories/:id", deps.AuthMiddleware, deps.MemberCategory.Update)
	app.Delete("/api/member-categories/:id", deps.AuthMiddleware, deps.MemberCategory.Delete)
	app.Get("/api/member-categories-combo", deps.AuthMiddleware, deps.MemberCategory.Combo)
	app.Get("/api/purchases", deps.AuthMiddleware, deps.Purchase.List)
	app.Get("/api/purchases/excel", deps.AuthMiddleware, deps.Export.PurchasesExcel)
	app.Get("/api/purchases/pdf", deps.AuthMiddleware, deps.Export.PurchasesPDF)
	app.Get("/api/purchase-items/excel", deps.AuthMiddleware, deps.Export.PurchaseItemsExcel)
	app.Get("/api/purchase-items/pdf", deps.AuthMiddleware, deps.Export.PurchaseItemsPDF)
	app.Get("/api/purchases/:id", deps.AuthMiddleware, deps.Purchase.GetByID)
	app.Post("/api/purchases", deps.AuthMiddleware, deps.Purchase.Create)
	app.Put("/api/purchases/:id", deps.AuthMiddleware, deps.Purchase.Update)
	app.Delete("/api/purchases/:id", deps.AuthMiddleware, deps.Purchase.Delete)
	app.Get("/api/purchase-items/all/:id", deps.AuthMiddleware, deps.Purchase.ListItems)
	app.Post("/api/purchase-items", deps.AuthMiddleware, deps.Purchase.CreateItem)
	app.Put("/api/purchase-items/:id", deps.AuthMiddleware, deps.Purchase.UpdateItem)
	app.Delete("/api/purchase-items/:id", deps.AuthMiddleware, deps.Purchase.DeleteItem)
	app.Get("/api/sales", deps.AuthMiddleware, deps.Sale.List)
	app.Get("/api/sales/excel", deps.AuthMiddleware, deps.Export.SalesExcel)
	app.Get("/api/sales/pdf", deps.AuthMiddleware, deps.Export.SalesPDF)
	app.Get("/api/sale-items/excel", deps.AuthMiddleware, deps.Export.SaleItemsExcel)
	app.Get("/api/sale-items/pdf", deps.AuthMiddleware, deps.Export.SaleItemsPDF)
	app.Get("/api/sales/:id", deps.AuthMiddleware, deps.Sale.GetByID)
	app.Post("/api/sales", deps.AuthMiddleware, deps.Sale.Create)
	app.Put("/api/sales/:id", deps.AuthMiddleware, deps.Sale.Update)
	app.Delete("/api/sales/:id", deps.AuthMiddleware, deps.Sale.Delete)
	app.Get("/api/sale-items/all/:id", deps.AuthMiddleware, deps.Sale.ListItems)
	app.Post("/api/sale-items", deps.AuthMiddleware, deps.Sale.CreateItem)
	app.Put("/api/sale-items/:id", deps.AuthMiddleware, deps.Sale.UpdateItem)
	app.Delete("/api/sale-items/:id", deps.AuthMiddleware, deps.Sale.DeleteItem)
	app.Post("/api/opnames", deps.AuthMiddleware, deps.Opname.Create)
	app.Get("/api/opnames/:id", deps.AuthMiddleware, deps.Opname.GetByID)
	app.Post("/api/opname-items", deps.AuthMiddleware, deps.Opname.CreateItem)
	app.Post("/api/opname-items-all", deps.AuthMiddleware, deps.Opname.GetItems)
}
