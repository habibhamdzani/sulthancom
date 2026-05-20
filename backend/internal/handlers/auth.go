package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sulthancom/pos-backend/internal/config"
	"github.com/sulthancom/pos-backend/internal/database"
	"github.com/sulthancom/pos-backend/internal/models"
	"github.com/sulthancom/pos-backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db    *database.Database
	redis *database.Redis
	cfg   *config.Config
}

func NewAuthHandler(db *database.Database, redis *database.Redis, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, redis: redis, cfg: cfg}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := h.db.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusUnauthorized, "Invalid email or password")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	if !user.IsActive {
		return utils.RespondError(c, fiber.StatusForbidden, "Account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.RespondError(c, fiber.StatusUnauthorized, "Invalid email or password")
	}

	token, err := utils.GenerateJWT(&user, h.cfg.JWTSecret, h.cfg.JWTExpiry)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	user.Password = ""

	response := models.LoginResponse{
		Token:     token,
		ExpiresIn: h.cfg.JWTExpiry,
		User:      user,
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Login successful", response)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Name     string       `json:"name"`
		Email    string       `json:"email"`
		Password string       `json:"password"`
		Role     models.Role  `json:"role"`
		Phone    string       `json:"phone"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, "Name, email, and password are required")
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Role == "" {
		req.Role = models.RoleKasir
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing models.User
	err := h.db.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing)
	if err == nil {
		return utils.RespondError(c, fiber.StatusConflict, "Email already registered")
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

	return utils.RespondSuccess(c, fiber.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = h.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return utils.RespondError(c, fiber.StatusNotFound, "User not found")
	}

	user.Password = ""
	return utils.RespondSuccess(c, fiber.StatusOK, "Success", user)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := authHeader[len("Bearer "):]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = h.redis.Set(ctx, "blacklist:"+token, "1", time.Duration(h.cfg.JWTExpiry)*time.Hour)

	return utils.RespondSuccess(c, fiber.StatusOK, "Logout successful", nil)
}
