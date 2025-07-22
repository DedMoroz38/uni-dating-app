package services

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/DedMoroz38/uni-dating-app/internal/config"
	database "github.com/DedMoroz38/uni-dating-app/internal/db"
)

type AuthController struct {
	db database.Querier
}

func NewAuthController(db database.Querier) *AuthController {
	return &AuthController{db: db}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
func (ac *AuthController) Register(c *fiber.Ctx) error {
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

	err = ac.db.CreateUser(c.Context(), database.CreateUserParams{
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

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}

	user, err := ac.db.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := createToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func createToken(ID int64) (string, error) {
	claims := jwt.MapClaims{
		"ID":      ID,
		"expiry":  time.Now().Add(config.JWTTokenExpiry).Unix(),
		"created": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret)
}
