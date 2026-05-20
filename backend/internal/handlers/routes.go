package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sulthancom/pos-backend/internal/config"
	"github.com/sulthancom/pos-backend/internal/database"
	"github.com/sulthancom/pos-backend/internal/middleware"
	"github.com/sulthancom/pos-backend/internal/models"
)

func Setup(app *fiber.App, db *database.Database, redis *database.Redis, cfg *config.Config) {
	authHandler := NewAuthHandler(db, redis, cfg)
	userHandler := NewUserHandler(db)
	productHandler := NewProductHandler(db)
	categoryHandler := NewCategoryHandler(db)
	brandHandler := NewBrandHandler(db)
	supplierHandler := NewSupplierHandler(db)
	customerHandler := NewCustomerHandler(db)
	transactionHandler := NewTransactionHandler(db)
	serviceHandler := NewServiceHandler(db)
	receivableHandler := NewReceivableHandler(db)
	inventoryHandler := NewInventoryHandler(db)
	reportHandler := NewReportHandler(db)
	settingsHandler := NewSettingsHandler(db)
	seedHandler := NewSeedHandler(db)

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "sulthancom-pos-api",
			"time":    time.Now().Unix(),
		})
	})

	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/register", authHandler.Register)
	api.Post("/seed", seedHandler.Seed)

	auth := api.Group("")
	auth.Use(middleware.JWTProtected(cfg))

	auth.Get("/auth/me", authHandler.GetMe)
	auth.Post("/auth/logout", authHandler.Logout)

	users := auth.Group("/users")
	users.Use(middleware.RequireRole(models.RoleSuperAdmin))
	users.Get("/", userHandler.GetAll)
	users.Get("/:id", userHandler.GetByID)
	users.Post("/", userHandler.Create)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	products := auth.Group("/products")
	products.Get("/", productHandler.GetAll)
	products.Get("/low-stock", productHandler.GetLowStock)
	products.Get("/:id", productHandler.GetByID)
	products.Post("/", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), productHandler.Create)
	products.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), productHandler.Update)
	products.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), productHandler.Delete)

	categories := auth.Group("/categories")
	categories.Get("/", categoryHandler.GetAll)
	categories.Get("/:id", categoryHandler.GetByID)
	categories.Post("/", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), categoryHandler.Create)
	categories.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), categoryHandler.Update)
	categories.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), categoryHandler.Delete)

	brands := auth.Group("/brands")
	brands.Get("/", brandHandler.GetAll)
	brands.Get("/:id", brandHandler.GetByID)
	brands.Post("/", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), brandHandler.Create)
	brands.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), brandHandler.Update)
	brands.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), brandHandler.Delete)

	suppliers := auth.Group("/suppliers")
	suppliers.Get("/", supplierHandler.GetAll)
	suppliers.Get("/:id", supplierHandler.GetByID)
	suppliers.Post("/", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), supplierHandler.Create)
	suppliers.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), supplierHandler.Update)
	suppliers.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin), supplierHandler.Delete)

	customers := auth.Group("/customers")
	customers.Get("/", customerHandler.GetAll)
	customers.Get("/:id", customerHandler.GetByID)
	customers.Post("/", customerHandler.Create)
	customers.Put("/:id", customerHandler.Update)
	customers.Delete("/:id", customerHandler.Delete)

	transactions := auth.Group("/transactions")
	transactions.Get("/", transactionHandler.GetAll)
	transactions.Get("/stats", transactionHandler.GetStats)
	transactions.Get("/:id", transactionHandler.GetByID)
	transactions.Post("/", middleware.RequireRole(models.RoleSuperAdmin, models.RoleKasir), transactionHandler.Create)

	services := auth.Group("/services")
	services.Get("/", serviceHandler.GetAll)
	services.Get("/:id", serviceHandler.GetByID)
	services.Post("/", middleware.RequireRole(models.RoleSuperAdmin, models.RoleTeknisi), serviceHandler.Create)
	services.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleTeknisi), serviceHandler.Update)
	services.Post("/:id/notes", middleware.RequireRole(models.RoleSuperAdmin, models.RoleTeknisi), serviceHandler.AddNote)
	services.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin, models.RoleTeknisi), serviceHandler.Delete)

	receivables := auth.Group("/receivables")
	receivables.Get("/", receivableHandler.GetAll)
	receivables.Get("/stats", receivableHandler.GetStats)
	receivables.Get("/:id", receivableHandler.GetByID)
	receivables.Post("/:id/payment", middleware.RequireRole(models.RoleSuperAdmin, models.RoleKasir), receivableHandler.MakePayment)

	inventory := auth.Group("/inventory")
	inventory.Use(middleware.RequireRole(models.RoleSuperAdmin, models.RoleAdmin))
	inventory.Post("/stock-in", inventoryHandler.StockIn)
	inventory.Post("/stock-out", inventoryHandler.StockOut)
	inventory.Post("/adjustment", inventoryHandler.StockAdjustment)
	inventory.Get("/movements", inventoryHandler.GetMovements)

	reports := auth.Group("/reports")
	reports.Get("/sales", reportHandler.SalesReport)
	reports.Get("/products", reportHandler.ProductReport)
	reports.Get("/services", reportHandler.ServiceReport)
	reports.Get("/financial", reportHandler.FinancialReport)

	settings := auth.Group("/settings")
	settings.Use(middleware.RequireRole(models.RoleSuperAdmin))
	settings.Get("/", settingsHandler.Get)
	settings.Put("/", settingsHandler.Update)
}
