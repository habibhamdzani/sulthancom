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
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *database.Database
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	page, limit := utils.ParsePaginationQuery(c.Query("page"), c.Query("limit"))
	search := c.Query("search")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": primitive.Regex{Pattern: search, Options: "i"}},
				{"email": primitive.Regex{Pattern: search, Options: "i"}},
			},
		}
	}

	total, err := h.db.Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to count users")
	}

	var users []models.User
	cursor, err := h.db.Collection("users").Find(ctx, filter, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch users")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode users")
	}

	for i := range users {
		users[i].Password = ""
	}

	return utils.RespondPaginated(c, users, total, page, limit)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = h.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "User not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch user")
	}

	user.Password = ""
	return utils.RespondSuccess(c, fiber.StatusOK, "Success", user)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Role     models.Role `json:"role"`
		Phone    string      `json:"phone"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing models.User
	err := h.db.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing)
	if err == nil {
		return utils.RespondError(c, fiber.StatusConflict, "Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	now := time.Now()
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      req.Role,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := h.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	user.Password = ""

	return utils.RespondSuccess(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req struct {
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Role     models.Role `json:"role"`
		Phone    string      `json:"phone"`
		IsActive *bool       `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.Role != "" {
		update["role"] = req.Role
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}

	result, err := h.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update user")
	}

	if result.MatchedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "User not found")
	}

	return h.GetByID(c)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("users").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to delete user")
	}

	if result.DeletedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "User not found")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "User deleted successfully", nil)
}
