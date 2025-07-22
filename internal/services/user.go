package services

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/DedMoroz38/uni-dating-app/internal/middleware"
)

type UserController struct {
	db db.Querier
}

func NewUserController(db db.Querier) *UserController {
	return &UserController{db: db}
}

func (uc *UserController) Me(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user retrieved"})
}

func (uc *UserController) UploadImages(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse form"})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no files uploaded"})
	}
	if len(files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot upload more than 10 images"})
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user ID from context"})
	}

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		path := filepath.Join("./uploads", filename)

		if err := c.SaveFile(file, path); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot save file"})
		}

		arg := db.CreateImageParams{
			UserID: int32(userID),
			Url:    path,
		}
		if err := uc.db.CreateImage(c.Context(), arg); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot save image to database"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "images uploaded successfully"})
}

func (uc *UserController) GetRandomUser(c *fiber.Ctx) error {
	user, err := uc.db.GetRandomUserWithImages(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get random user"})
	}

	var images []struct {
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

	imageURLs := make([]string, 0, len(images))
	for _, img := range images {
		urlPath := strings.TrimPrefix(img.URL, "./")
		imageURLs = append(imageURLs, urlPath)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"images":   imageURLs,
	})
}

func (uc *UserController) LikeUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	if err := uc.db.LikeUser(c.Context(), int32(id)); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot like user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("user %d liked", id)})
}
