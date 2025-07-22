package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/DedMoroz38/uni-dating-app/internal/config"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "malformed authorization header"})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return config.JWTSecret, nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired JWT"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check expiry
			expiry := int64(claims["expiry"].(float64))
			if time.Now().Unix() > expiry {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
			}
			c.Locals("userID", claims["ID"])
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
}

func GetUserIDFromContext(c *fiber.Ctx) (int64, error) {
	idVal, ok := c.Locals("userID").(float64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusInternalServerError, "user ID not found in context or is of wrong type")
	}
	return int64(idVal), nil
}
