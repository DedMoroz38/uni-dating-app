package controllers

import (
	"log"
	"strings"
	"time"

	"github.com/DedMoroz38/uni-dating-app/internal/config"
	database "github.com/DedMoroz38/uni-dating-app/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	DB database.Querier
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
	CourseID    int32  `json:"courseId"`
	Email       string `json:"email"`
	Password    string `json:"password"`
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

	dob, err := time.Parse(time.DateOnly, req.DateOfBirth)
	if err != nil {
		log.Println("Invalid date format for dateOfBirth", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date format for dateOfBirth"})
	}

	userId, err := ac.DB.CreateUserAndReturnID(c.Context(), database.CreateUserAndReturnIDParams{
		Name: req.Name,
		DateOfBirth: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
		Email:    req.Email,
		Password: string(hashedPassword),
		CourseID: req.CourseID,
	})
	if err != nil {
		log.Println("Failed to create user", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}

	token, err := createToken(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request"})
	}

	user, err := ac.DB.GetUserByEmail(c.Context(), req.Email)
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
