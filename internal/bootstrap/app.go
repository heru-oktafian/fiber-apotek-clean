package bootstrap

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	jwtadapter "github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/auth/jwt"
	redisadapter "github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/cache/redis"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/handlers"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/middleware"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/http/fiber/router"
	postgresadapter "github.com/heru-oktafian/fiber-apotek-clean/internal/adapters/persistence/postgres"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/clock"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/config"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/console"
	"github.com/heru-oktafian/fiber-apotek-clean/internal/shared/idgen"
	authusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/auth"
	anotherincomeusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/anotherincome"
	branchusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/branch"
	expenseusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/expense"
	opnameusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/opname"
	membercategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/membercategory"
	productusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/product"
	productcategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/productcategory"
	purchaseusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/purchase"
	saleusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/sale"
	supplierusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/supplier"
	suppliercategoryusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/suppliercategory"
	unitusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/unit"
	userbranchusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/userbranch"
	userusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/user"
)

type App struct {
	Fiber  *fiber.App
	Config config.Config
}

type bcryptComparer struct{}

type bcryptHasher struct{}

type exportBundle struct {
	base        handlers.ExportHandler
	master      handlers.ExportMasterHandler
	transaction handlers.ExportTransactionHandler
	finance     handlers.ExportFinanceHandler
}

func (e exportBundle) ProductsExcel(c *fiber.Ctx) error            { return e.base.ProductsExcel(c) }
func (e exportBundle) ProductsPDF(c *fiber.Ctx) error              { return e.base.ProductsPDF(c) }
func (e exportBundle) UnitsExcel(c *fiber.Ctx) error               { return e.base.UnitsExcel(c) }
func (e exportBundle) UnitsPDF(c *fiber.Ctx) error                 { return e.base.UnitsPDF(c) }
func (e exportBundle) ProductCategoriesExcel(c *fiber.Ctx) error   { return e.master.ProductCategoriesExcel(c) }
func (e exportBundle) ProductCategoriesPDF(c *fiber.Ctx) error     { return e.master.ProductCategoriesPDF(c) }
func (e exportBundle) SuppliersExcel(c *fiber.Ctx) error           { return e.master.SuppliersExcel(c) }
func (e exportBundle) SuppliersPDF(c *fiber.Ctx) error             { return e.master.SuppliersPDF(c) }
func (e exportBundle) SupplierCategoriesExcel(c *fiber.Ctx) error  { return e.master.SupplierCategoriesExcel(c) }
func (e exportBundle) SupplierCategoriesPDF(c *fiber.Ctx) error    { return e.master.SupplierCategoriesPDF(c) }
func (e exportBundle) MemberCategoriesExcel(c *fiber.Ctx) error    { return e.master.MemberCategoriesExcel(c) }
func (e exportBundle) MemberCategoriesPDF(c *fiber.Ctx) error      { return e.master.MemberCategoriesPDF(c) }
func (e exportBundle) PurchasesExcel(c *fiber.Ctx) error           { return e.transaction.PurchasesExcel(c) }
func (e exportBundle) PurchasesPDF(c *fiber.Ctx) error             { return e.transaction.PurchasesPDF(c) }
func (e exportBundle) PurchaseItemsExcel(c *fiber.Ctx) error       { return e.transaction.PurchaseItemsExcel(c) }
func (e exportBundle) PurchaseItemsPDF(c *fiber.Ctx) error         { return e.transaction.PurchaseItemsPDF(c) }
func (e exportBundle) SalesExcel(c *fiber.Ctx) error               { return e.transaction.SalesExcel(c) }
func (e exportBundle) SalesPDF(c *fiber.Ctx) error                 { return e.transaction.SalesPDF(c) }
func (e exportBundle) SaleItemsExcel(c *fiber.Ctx) error           { return e.transaction.SaleItemsExcel(c) }
func (e exportBundle) SaleItemsPDF(c *fiber.Ctx) error             { return e.transaction.SaleItemsPDF(c) }
func (e exportBundle) AnotherIncomesExcel(c *fiber.Ctx) error      { return e.finance.AnotherIncomesExcel(c) }
func (e exportBundle) AnotherIncomesPDF(c *fiber.Ctx) error        { return e.finance.AnotherIncomesPDF(c) }
func (e exportBundle) ExpensesExcel(c *fiber.Ctx) error            { return e.finance.ExpensesExcel(c) }
func (e exportBundle) ExpensesPDF(c *fiber.Ctx) error              { return e.finance.ExpensesPDF(c) }

func (bcryptComparer) Compare(hashed string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func (bcryptHasher) Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func New() (*App, error) {
	_ = validator.New()
	_ = godotenv.Load()
	cfg := config.Load()
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.TimezoneName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	redisDB, _ := strconv.Atoi(cfg.RedisDB)
	rdb := goredis.NewClient(&goredis.Options{Addr: cfg.RedisHost + ":" + cfg.RedisPort, Password: cfg.RedisPass, DB: redisDB})
	repos := postgresadapter.Repositories{DB: db}
	jwtSvc := jwtadapter.Service{Secret: []byte(cfg.JWTSecret)}
	blacklist := redisadapter.Blacklist{Client: rdb}
	clk := clock.RealClock{}
	ids := idgen.Generator{}

	authHandler := handlers.AuthHandler{Service: authusecase.Service{Users: repos, Branches: repos, Passwords: bcryptComparer{}, Tokens: jwtSvc, Blacklist: blacklist, Clock: clk}}
	branchHandler := handlers.BranchHandler{Service: branchusecase.Service{Branches: repos, IDs: ids}}
	userBranchHandler := handlers.UserBranchHandler{Service: userbranchusecase.Service{Branches: repos, Users: repos}}
	userHandler := handlers.UserHandler{Service: userusecase.Service{Users: repos, Passwords: bcryptHasher{}, IDs: ids}}
	productHandler := handlers.ProductHandler{Service: productusecase.Service{Products: repos, IDs: ids}}
	supplierHandler := handlers.SupplierHandler{Service: supplierusecase.Service{Suppliers: repos, IDs: ids}}
	unitHandler := handlers.UnitHandler{Service: unitusecase.MasterService{Units: repos, IDs: ids}}
	productCategoryHandler := handlers.ProductCategoryHandler{Service: productcategoryusecase.Service{Categories: repos}}
	supplierCategoryHandler := handlers.SupplierCategoryHandler{Service: suppliercategoryusecase.Service{Categories: repos}}
	memberCategoryHandler := handlers.MemberCategoryHandler{Service: membercategoryusecase.Service{Categories: repos}}
	anotherIncomeHandler := handlers.AnotherIncomeHandler{Service: anotherincomeusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	expenseHandler := handlers.ExpenseHandler{Service: expenseusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	purchaseHandler := handlers.PurchaseHandler{Service: purchaseusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	saleHandler := handlers.SaleHandler{Service: saleusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	opnameHandler := handlers.OpnameHandler{Service: opnameusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	exportHandler := handlers.ExportHandler{Products: productusecase.Service{Products: repos, IDs: ids}, Units: unitusecase.MasterService{Units: repos, IDs: ids}}
	exportMasterHandler := handlers.ExportMasterHandler{ProductCategories: productcategoryusecase.Service{Categories: repos}, Suppliers: supplierusecase.Service{Suppliers: repos, IDs: ids}, SupplierCategories: suppliercategoryusecase.Service{Categories: repos}, MemberCategories: membercategoryusecase.Service{Categories: repos}}
	exportTransactionHandler := handlers.ExportTransactionHandler{Purchases: purchaseusecase.Service{Repo: repos, IDs: ids, Clock: clk}, Sales: saleusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	exportFinanceHandler := handlers.ExportFinanceHandler{AnotherIncomes: anotherincomeusecase.Service{Repo: repos, IDs: ids, Clock: clk}, Expenses: expenseusecase.Service{Repo: repos, IDs: ids, Clock: clk}}

	app := fiber.New(fiber.Config{DisableStartupMessage: true, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second})
	app.Use(console.RequestLogger())
	authMw := middleware.RequireAuth(jwtSvc, blacklist)
	router.Register(app, router.Dependencies{Auth: authHandler, Branch: branchHandler, UserBranch: userBranchHandler, User: userHandler, Product: productHandler, Supplier: supplierHandler, Unit: unitHandler, ProductCategory: productCategoryHandler, SupplierCategory: supplierCategoryHandler, MemberCategory: memberCategoryHandler, AnotherIncome: anotherIncomeHandler, Expense: expenseHandler, Purchase: purchaseHandler, Sale: saleHandler, Opname: opnameHandler, Export: exportBundle{base: exportHandler, master: exportMasterHandler, transaction: exportTransactionHandler, finance: exportFinanceHandler}, AuthMiddleware: authMw})
	return &App{Fiber: app, Config: cfg}, nil
}
