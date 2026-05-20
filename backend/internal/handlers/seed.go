package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sulthancom/pos-backend/internal/database"
	"github.com/sulthancom/pos-backend/internal/models"
	"github.com/sulthancom/pos-backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SeedHandler struct {
	db *database.Database
}

func NewSeedHandler(db *database.Database) *SeedHandler {
	return &SeedHandler{db: db}
}

func (h *SeedHandler) Seed(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, _ := h.db.Collection("users").CountDocuments(ctx, bson.M{})
	if count > 0 {
		return utils.RespondError(c, fiber.StatusConflict, "Database already seeded")
	}

	users := []interface{}{
		models.User{
			ID:        primitive.NewObjectID(),
			Name:      "Sulthan Admin",
			Email:     "admin@sulthancom.com",
			Password:  "$2a$10$WEPZQgJlQozbRN28GV14V.XZ7Sa9Lykm8As7YGyOdjOYEUtOx0N2a",
			Role:      models.RoleSuperAdmin,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		models.User{
			ID:        primitive.NewObjectID(),
			Name:      "Budi Inventory",
			Email:     "budi@sulthancom.com",
			Password:  "$2a$10$WEPZQgJlQozbRN28GV14V.XZ7Sa9Lykm8As7YGyOdjOYEUtOx0N2a",
			Role:      models.RoleAdmin,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		models.User{
			ID:        primitive.NewObjectID(),
			Name:      "Ani Kasir",
			Email:     "ani@sulthancom.com",
			Password:  "$2a$10$WEPZQgJlQozbRN28GV14V.XZ7Sa9Lykm8As7YGyOdjOYEUtOx0N2a",
			Role:      models.RoleKasir,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		models.User{
			ID:        primitive.NewObjectID(),
			Name:      "Rudi Teknisi",
			Email:     "rudi@sulthancom.com",
			Password:  "$2a$10$WEPZQgJlQozbRN28GV14V.XZ7Sa9Lykm8As7YGyOdjOYEUtOx0N2a",
			Role:      models.RoleTeknisi,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	h.db.Collection("users").InsertMany(ctx, users)

	categories := []interface{}{
		models.Category{ID: primitive.NewObjectID(), Name: "Laptop", Icon: "laptop", Color: "#3B82F6", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Category{ID: primitive.NewObjectID(), Name: "PC Rakitan", Icon: "monitor", Color: "#8B5CF6", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Category{ID: primitive.NewObjectID(), Name: "Printer", Icon: "printer", Color: "#10B981", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Category{ID: primitive.NewObjectID(), Name: "Monitor", Icon: "tv", Color: "#F59E0B", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Category{ID: primitive.NewObjectID(), Name: "Sparepart", Icon: "cpu", Color: "#EF4444", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Category{ID: primitive.NewObjectID(), Name: "Networking", Icon: "wifi", Color: "#06B6D4", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Category{ID: primitive.NewObjectID(), Name: "Aksesoris", Icon: "headphones", Color: "#EC4899", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Category{ID: primitive.NewObjectID(), Name: "Service", Icon: "wrench", Color: "#6366F1", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	h.db.Collection("categories").InsertMany(ctx, categories)

	brands := []interface{}{
		models.Brand{ID: primitive.NewObjectID(), Name: "ASUS", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Brand{ID: primitive.NewObjectID(), Name: "Lenovo", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Brand{ID: primitive.NewObjectID(), Name: "HP", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Brand{ID: primitive.NewObjectID(), Name: "Dell", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Brand{ID: primitive.NewObjectID(), Name: "Acer", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Brand{ID: primitive.NewObjectID(), Name: "MSI", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Brand{ID: primitive.NewObjectID(), Name: "Logitech", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		models.Brand{ID: primitive.NewObjectID(), Name: "Epson", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	h.db.Collection("brands").InsertMany(ctx, brands)

	settings := models.StoreSettings{
		ID:            primitive.NewObjectID(),
		Name:          "SulthanCom",
		Address:       "Jl. Raya Komputer No. 123, Surabaya",
		Phone:         "031-12345678",
		Email:         "info@sulthancom.com",
		Website:       "sulthancom.com",
		TaxRate:       0,
		Currency:      "IDR",
		InvoiceFooter: "Terima kasih atas kunjungan Anda. Garansi berlaku sesuai ketentuan.",
		Theme:         "light",
		UpdatedAt:     time.Now(),
	}

	h.db.Collection("settings").InsertOne(ctx, settings)

	h.db.Collection("products").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"sku": 1},
		Options: options.Index().SetUnique(true),
	})

	h.db.Collection("products").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"barcode": 1},
		Options: options.Index().SetUnique(true).SetSparse(true),
	})

	h.db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})

	return utils.RespondSuccess(c, fiber.StatusOK, "Database seeded successfully", bson.M{
		"users":      len(users),
		"categories": len(categories),
		"brands":     len(brands),
	})
}
