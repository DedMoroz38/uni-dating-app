package routes

import (
	"github.com/DedMoroz38/uni-dating-app/internal/controllers"
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, queries *db.Queries) {
	userRoutes := router.Group("/user")
	uc := controllers.UserController{
		DB: queries,
	}

	userRoutes.Get("/me", uc.Me)
	userRoutes.Post("/images", uc.UploadImages)
	userRoutes.Get("/card/get_random", uc.GetRandomUser)
	userRoutes.Post("/like/:userId", uc.LikeUser)

	userRoutes.Post("/seed", uc.Seed)
}
