package controllers

import (
	"database/sql"

	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/gofiber/fiber/v2"
)

type CourseController struct {
	DB *db.Queries
}

func (cc *CourseController) SeedCourses(c *fiber.Ctx) error {
	coursesToSeed := []string{"computer science", "mathematics", "physics"}
	for _, courseName := range coursesToSeed {
		_, err := cc.DB.CreateCourse(c.Context(), courseName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not seed course"})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "courses seeded successfully"})
}

func (cc *CourseController) GetCourses(c *fiber.Ctx) error {
	courses, err := cc.DB.ListCourses(c.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusOK).JSON([]db.Course{})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get courses"})
	}

	return c.JSON(courses)
}
