package authenhandler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"my-source/sheet-payment/be/repository"
)

var JwtSecret = []byte("j7akAxU")

type AuthHandler struct {
	UserRepo repository.IUserRepository
}

func NewAuthHandler(usr repository.IUserRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo: usr,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	user := &repository.User{
		ID:       uuid.New().String(),
		Username: body.Username,
		Password: string(hashed),
	}

	if err := h.UserRepo.Create(user); err != nil {
		return fiber.NewError(fiber.StatusConflict, "username already exists")
	}
	return c.JSON(fiber.Map{"message": "user created"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	user, err := h.UserRepo.GetByUsername(body.Username)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return fiber.ErrUnauthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString(JwtSecret)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"token": tokenStr})
}
