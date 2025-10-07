package routes

import (
	"github.com/DedMoroz38/uni-dating-app/internal/controllers"
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/gofiber/fiber/v2"
)

func ProfileRoutes(router fiber.Router, queries *db.Queries) {
	userRoutes := router.Group("/profile")
	uc := controllers.ProfileController{
		DB: queries,
	}

	userRoutes.Get("/", uc.GetAll)
}
