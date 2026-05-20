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

type InventoryHandler struct {
	db *database.Database
}

func NewInventoryHandler(db *database.Database) *InventoryHandler {
	return &InventoryHandler{db: db}
}

func (h *InventoryHandler) StockIn(c *fiber.Ctx) error {
	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
		Notes     string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil || req.Quantity <= 0 {
		return utils.RespondError(c, fiber.StatusBadRequest, "Valid quantity is required")
	}

	productOID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err = h.db.Collection("products").FindOne(ctx, bson.M{"_id": productOID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch product")
	}

	previousStock := product.Stock
	product.Stock += req.Quantity

	_, err = h.db.Collection("products").UpdateOne(ctx, bson.M{"_id": productOID}, bson.M{
		"$set": bson.M{"stock": product.Stock, "updated_at": time.Now()},
	})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update stock")
	}

	movement := models.StockMovement{
		ID:            primitive.NewObjectID(),
		ProductID:     productOID,
		ProductName:   product.Name,
		Type:          models.StockIn,
		Quantity:      req.Quantity,
		PreviousStock: previousStock,
		NewStock:      product.Stock,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		CreatedBy:     c.Locals("user_name").(string),
	}

	h.db.Collection("stock_movements").InsertOne(ctx, movement)

	return utils.RespondSuccess(c, fiber.StatusOK, "Stock updated successfully", bson.M{
		"product_id":     productOID,
		"product_name":   product.Name,
		"previous_stock": previousStock,
		"new_stock":      product.Stock,
		"quantity_added": req.Quantity,
	})
}

func (h *InventoryHandler) StockOut(c *fiber.Ctx) error {
	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
		Notes     string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil || req.Quantity <= 0 {
		return utils.RespondError(c, fiber.StatusBadRequest, "Valid quantity is required")
	}

	productOID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err = h.db.Collection("products").FindOne(ctx, bson.M{"_id": productOID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch product")
	}

	if product.Stock < req.Quantity {
		return utils.RespondError(c, fiber.StatusBadRequest, "Insufficient stock")
	}

	previousStock := product.Stock
	product.Stock -= req.Quantity

	_, err = h.db.Collection("products").UpdateOne(ctx, bson.M{"_id": productOID}, bson.M{
		"$set": bson.M{"stock": product.Stock, "updated_at": time.Now()},
	})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update stock")
	}

	movement := models.StockMovement{
		ID:            primitive.NewObjectID(),
		ProductID:     productOID,
		ProductName:   product.Name,
		Type:          models.StockOut,
		Quantity:      req.Quantity,
		PreviousStock: previousStock,
		NewStock:      product.Stock,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		CreatedBy:     c.Locals("user_name").(string),
	}

	h.db.Collection("stock_movements").InsertOne(ctx, movement)

	return utils.RespondSuccess(c, fiber.StatusOK, "Stock updated successfully", bson.M{
		"product_id":       productOID,
		"product_name":     product.Name,
		"previous_stock":   previousStock,
		"new_stock":        product.Stock,
		"quantity_removed": req.Quantity,
	})
}

func (h *InventoryHandler) StockAdjustment(c *fiber.Ctx) error {
	var req struct {
		ProductID string `json:"product_id"`
		NewStock  int    `json:"new_stock"`
		Notes     string `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	productOID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err = h.db.Collection("products").FindOne(ctx, bson.M{"_id": productOID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch product")
	}

	previousStock := product.Stock
	quantity := req.NewStock - previousStock

	_, err = h.db.Collection("products").UpdateOne(ctx, bson.M{"_id": productOID}, bson.M{
		"$set": bson.M{"stock": req.NewStock, "updated_at": time.Now()},
	})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update stock")
	}

	movement := models.StockMovement{
		ID:            primitive.NewObjectID(),
		ProductID:     productOID,
		ProductName:   product.Name,
		Type:          models.StockAdjustment,
		Quantity:      quantity,
		PreviousStock: previousStock,
		NewStock:      req.NewStock,
		Notes:         req.Notes,
		CreatedAt:     time.Now(),
		CreatedBy:     c.Locals("user_name").(string),
	}

	h.db.Collection("stock_movements").InsertOne(ctx, movement)

	return utils.RespondSuccess(c, fiber.StatusOK, "Stock adjusted successfully", bson.M{
		"product_id":     productOID,
		"product_name":   product.Name,
		"previous_stock": previousStock,
		"new_stock":      req.NewStock,
	})
}

func (h *InventoryHandler) GetMovements(c *fiber.Ctx) error {
	page, limit := utils.ParsePaginationQuery(c.Query("page"), c.Query("limit"))
	productID := c.Query("product_id")
	movementType := c.Query("type")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if productID != "" {
		if oid, err := primitive.ObjectIDFromHex(productID); err == nil {
			filter["product_id"] = oid
		}
	}
	if movementType != "" {
		filter["type"] = movementType
	}

	total, err := h.db.Collection("stock_movements").CountDocuments(ctx, filter)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to count movements")
	}

	var movements []models.StockMovement
	cursor, err := h.db.Collection("stock_movements").Find(ctx, filter, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch movements")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movements); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode movements")
	}

	return utils.RespondPaginated(c, movements, total, page, limit)
}
