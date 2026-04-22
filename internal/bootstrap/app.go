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
	branchusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/branch"
	opnameusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/opname"
	productusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/product"
	purchaseusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/purchase"
	saleusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/sale"
	userbranchusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/userbranch"
	userusecase "github.com/heru-oktafian/fiber-apotek-clean/internal/usecase/user"
)

type App struct {
	Fiber  *fiber.App
	Config config.Config
}

type bcryptComparer struct{}

type bcryptHasher struct{}

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
	branchHandler := handlers.BranchHandler{Service: branchusecase.Service{Branches: repos}}
	userBranchHandler := handlers.UserBranchHandler{Service: userbranchusecase.Service{Branches: repos}}
	userHandler := handlers.UserHandler{Service: userusecase.Service{Users: repos, Passwords: bcryptHasher{}, IDs: ids}}
	productHandler := handlers.ProductHandler{Service: productusecase.Service{Products: repos, IDs: ids}}
	purchaseHandler := handlers.PurchaseHandler{Service: purchaseusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	saleHandler := handlers.SaleHandler{Service: saleusecase.Service{Repo: repos, IDs: ids, Clock: clk}}
	opnameHandler := handlers.OpnameHandler{Service: opnameusecase.Service{Repo: repos, IDs: ids, Clock: clk}}

	app := fiber.New(fiber.Config{DisableStartupMessage: true, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second})
	app.Use(console.RequestLogger())
	authMw := middleware.RequireAuth(jwtSvc, blacklist)
	router.Register(app, router.Dependencies{Auth: authHandler, Branch: branchHandler, UserBranch: userBranchHandler, User: userHandler, Product: productHandler, Purchase: purchaseHandler, Sale: saleHandler, Opname: opnameHandler, AuthMiddleware: authMw})
	return &App{Fiber: app, Config: cfg}, nil
}
