package main

import (
	"context"
	"log"
	"os"

	_ "github.com/DedMoroz38/uni-dating-app/docs"
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/DedMoroz38/uni-dating-app/internal/middleware"
	"github.com/DedMoroz38/uni-dating-app/internal/routes"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Uni Dating App API
// @version 1.0
// @description This is a dating app for university students.
// @host localhost:3000
// @BasePath /api/v1
func main() {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE environment variable is not set")
	}

	dbpool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	queries := db.New(dbpool)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Static("/uploads", "./uploads")

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.AuthRoutes(v1, queries)

	api.Use(middleware.Protected())

	routes.UserRoutes(v1, queries)

	log.Fatal(app.Listen(":3000"))
}
