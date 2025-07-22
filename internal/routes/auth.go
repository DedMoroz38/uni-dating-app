package routes

import (
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/DedMoroz38/uni-dating-app/internal/services"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router, queries *db.Queries) {
	authRoutes := router.Group("/auth")
	authController := services.NewAuthController(queries)

	authRoutes.Post("/register", authController.Register)
	authRoutes.Post("/login", authController.Login)
}
