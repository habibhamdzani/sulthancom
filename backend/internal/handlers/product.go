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

type ProductHandler struct {
	db *database.Database
}

func NewProductHandler(db *database.Database) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	page, limit := utils.ParsePaginationQuery(c.Query("page"), c.Query("limit"))
	search := c.Query("search")
	categoryID := c.Query("category_id")
	brandID := c.Query("brand_id")
	productType := c.Query("type")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": primitive.Regex{Pattern: search, Options: "i"}},
			{"sku": primitive.Regex{Pattern: search, Options: "i"}},
			{"barcode": primitive.Regex{Pattern: search, Options: "i"}},
		}
	}
	if categoryID != "" {
		if oid, err := primitive.ObjectIDFromHex(categoryID); err == nil {
			filter["category_id"] = oid
		}
	}
	if brandID != "" {
		if oid, err := primitive.ObjectIDFromHex(brandID); err == nil {
			filter["brand_id"] = oid
		}
	}
	if productType != "" {
		filter["type"] = productType
	}

	total, err := h.db.Collection("products").CountDocuments(ctx, filter)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to count products")
	}

	var products []models.Product
	cursor, err := h.db.Collection("products").Find(ctx, filter, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch products")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode products")
	}

	return utils.RespondPaginated(c, products, total, page, limit)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err = h.db.Collection("products").FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Product not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch product")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	if product.Images == nil {
		product.Images = []string{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("products").InsertOne(ctx, product)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create product")
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return utils.RespondSuccess(c, fiber.StatusCreated, "Product created successfully", product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	updates["updated_at"] = time.Now()
	delete(updates, "id")
	delete(updates, "created_at")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("products").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update product")
	}

	if result.MatchedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Product not found")
	}

	return h.GetByID(c)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("products").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to delete product")
	}

	if result.DeletedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Product not found")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetLowStock(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var products []models.Product
	cursor, err := h.db.Collection("products").Find(ctx, bson.M{
		"$expr": bson.M{"$lte": []string{"$stock", "$min_stock"}},
	})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch low stock products")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode products")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", products)
}
