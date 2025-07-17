package main

import (
	"context"
	"log"
	"os"

	"github.com/DedMoroz38/uni-dating-app/cmd/server/routes"
	_ "github.com/DedMoroz38/uni-dating-app/docs"
	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/gofiber/fiber/v2"
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

	routes.ApiRoutes(app, queries)

	log.Fatal(app.Listen(":3000"))
}
