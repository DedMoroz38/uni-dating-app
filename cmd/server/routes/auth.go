package routes

import (
	"github.com/DedMoroz38/uni-dating-app/api/v1/auth"
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router, queries *db.Queries) {
	authRoutes := router.Group("/auth")

	authRoutes.Post("/register", auth.Register(queries))
}
