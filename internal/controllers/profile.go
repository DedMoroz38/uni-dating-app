package controllers

import (
	"encoding/json"
	"strings"

	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/DedMoroz38/uni-dating-app/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	DB db.Querier
}

func (pc *ProfileController) GetAll(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user ID from context"})
	}

	user, err := pc.DB.GetUserByIDWithImages(c.Context(), int64(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user"})
	}

	var images []struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
	}
	if user.Images != nil {
		imagesJSON, err := json.Marshal(user.Images)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to re-marshal images"})
		}

		if err := json.Unmarshal(imagesJSON, &images); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse image data"})
		}
	}

	for i := range images {
		images[i].URL = strings.TrimPrefix(images[i].URL, "./")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"images": images,
	})
}
