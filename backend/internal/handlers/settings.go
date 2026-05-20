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

type SettingsHandler struct {
	db *database.Database
}

func NewSettingsHandler(db *database.Database) *SettingsHandler {
	return &SettingsHandler{db: db}
}

func (h *SettingsHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var settings models.StoreSettings
	err := h.db.Collection("settings").FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			defaultSettings := models.StoreSettings{
				ID:        primitive.NewObjectID(),
				Name:      "SulthanCom",
				Address:   "Jl. Raya Komputer No. 123, Surabaya",
				Phone:     "031-12345678",
				Email:     "info@sulthancom.com",
				TaxRate:   0,
				Currency:  "IDR",
				Theme:     "light",
				UpdatedAt: time.Now(),
			}
			h.db.Collection("settings").InsertOne(ctx, defaultSettings)
			return utils.RespondSuccess(c, fiber.StatusOK, "Success", defaultSettings)
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch settings")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", settings)
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	var settings models.StoreSettings
	if err := c.BodyParser(&settings); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	settings.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("settings").UpdateOne(ctx, bson.M{}, bson.M{"$set": settings}, options.Update().SetUpsert(true))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update settings")
	}

	if result.UpsertedCount > 0 {
		settings.ID = result.UpsertedID.(primitive.ObjectID)
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Settings updated successfully", settings)
}
