package auth

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	database "github.com/DedMoroz38/uni-dating-app/internal/db"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Register a new user
// @Description Register a new user with email and password. The email must be from the @soton.ac.uk domain.
// @Tags auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Register Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func Register(db database.Querier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
		}

		if !strings.HasSuffix(req.Email, "@soton.ac.uk") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid email domain"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
		}

		err = db.CreateUser(c.Context(), database.CreateUserParams{
			Username: req.Username,
			Email:    req.Email,
			Password: string(hashedPassword),
			CreatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			log.Println("Failed to create user", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user registered successfully"})
	}
}
