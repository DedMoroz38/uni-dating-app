package routes

import (
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(app *fiber.App, queries *db.Queries) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/swagger/*", swagger.HandlerDefault)

	AuthRoutes(v1, queries)
}
