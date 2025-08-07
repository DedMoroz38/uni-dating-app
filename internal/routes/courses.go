package routes

import (
	"github.com/DedMoroz38/uni-dating-app/internal/controllers"
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/gofiber/fiber/v2"
)

func CourseRoutes(router fiber.Router, queries *db.Queries) {
	courseRoutes := router.Group("/courses")
	cc := controllers.CourseController{
		DB: queries,
	}

	courseRoutes.Post("/seed", cc.SeedCourses)
	courseRoutes.Get("/", cc.GetCourses)
}
