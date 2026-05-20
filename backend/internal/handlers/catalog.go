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

type CategoryHandler struct {
	db *database.Database
}

func NewCategoryHandler(db *database.Database) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var categories []models.Category
	cursor, err := h.db.Collection("categories").Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch categories")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &categories); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode categories")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", categories)
}

func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var category models.Category
	err = h.db.Collection("categories").FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Category not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch category")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", category)
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("categories").InsertOne(ctx, category)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create category")
	}

	category.ID = result.InsertedID.(primitive.ObjectID)
	return utils.RespondSuccess(c, fiber.StatusCreated, "Category created successfully", category)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid category ID")
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

	result, err := h.db.Collection("categories").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update category")
	}

	if result.MatchedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Category not found")
	}

	return h.GetByID(c)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("categories").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to delete category")
	}

	if result.DeletedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Category not found")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Category deleted successfully", nil)
}

type BrandHandler struct {
	db *database.Database
}

func NewBrandHandler(db *database.Database) *BrandHandler {
	return &BrandHandler{db: db}
}

func (h *BrandHandler) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var brands []models.Brand
	cursor, err := h.db.Collection("brands").Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch brands")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &brands); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode brands")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", brands)
}

func (h *BrandHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid brand ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var brand models.Brand
	err = h.db.Collection("brands").FindOne(ctx, bson.M{"_id": objectID}).Decode(&brand)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Brand not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch brand")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", brand)
}

func (h *BrandHandler) Create(c *fiber.Ctx) error {
	var brand models.Brand
	if err := c.BodyParser(&brand); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	brand.ID = primitive.NewObjectID()
	brand.CreatedAt = time.Now()
	brand.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("brands").InsertOne(ctx, brand)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create brand")
	}

	brand.ID = result.InsertedID.(primitive.ObjectID)
	return utils.RespondSuccess(c, fiber.StatusCreated, "Brand created successfully", brand)
}

func (h *BrandHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid brand ID")
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

	result, err := h.db.Collection("brands").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update brand")
	}

	if result.MatchedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Brand not found")
	}

	return h.GetByID(c)
}

func (h *BrandHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid brand ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("brands").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to delete brand")
	}

	if result.DeletedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Brand not found")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Brand deleted successfully", nil)
}

type SupplierHandler struct {
	db *database.Database
}

func NewSupplierHandler(db *database.Database) *SupplierHandler {
	return &SupplierHandler{db: db}
}

func (h *SupplierHandler) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var suppliers []models.Supplier
	cursor, err := h.db.Collection("suppliers").Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch suppliers")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &suppliers); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode suppliers")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", suppliers)
}

func (h *SupplierHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid supplier ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var supplier models.Supplier
	err = h.db.Collection("suppliers").FindOne(ctx, bson.M{"_id": objectID}).Decode(&supplier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Supplier not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch supplier")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", supplier)
}

func (h *SupplierHandler) Create(c *fiber.Ctx) error {
	var supplier models.Supplier
	if err := c.BodyParser(&supplier); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	supplier.ID = primitive.NewObjectID()
	supplier.CreatedAt = time.Now()
	supplier.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("suppliers").InsertOne(ctx, supplier)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create supplier")
	}

	supplier.ID = result.InsertedID.(primitive.ObjectID)
	return utils.RespondSuccess(c, fiber.StatusCreated, "Supplier created successfully", supplier)
}

func (h *SupplierHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid supplier ID")
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

	result, err := h.db.Collection("suppliers").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update supplier")
	}

	if result.MatchedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Supplier not found")
	}

	return h.GetByID(c)
}

func (h *SupplierHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid supplier ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("suppliers").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to delete supplier")
	}

	if result.DeletedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Supplier not found")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Supplier deleted successfully", nil)
}
