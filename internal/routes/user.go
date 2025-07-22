package routes

import (
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/DedMoroz38/uni-dating-app/internal/services"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router, queries *db.Queries) {
	userRoutes := router.Group("/user")
	userController := services.NewUserController(queries)

	userRoutes.Get("/me", userController.Me)
	userRoutes.Post("/images", userController.UploadImages)
	userRoutes.Get("/card/get_random", userController.GetRandomUser)
	userRoutes.Post("/like/:userId", userController.LikeUser)
	userRoutes.Post("/seed", userController.Seed)
}
