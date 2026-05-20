package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sulthancom/pos-backend/internal/models"
)

type JWTClaims struct {
	UserID string        `json:"user_id"`
	Email  string        `json:"email"`
	Role   models.Role   `json:"role"`
	Name   string        `json:"name"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User, secret string, expiryHours int) (string, error) {
	claims := JWTClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   user.Role,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func RespondSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func RespondError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func RespondPaginated(c *fiber.Ctx, data interface{}, total int64, page, limit int) error {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items":       data,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}
