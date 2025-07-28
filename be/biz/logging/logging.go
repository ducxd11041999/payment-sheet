package middlewarelogging

import (
	"log"
	authenhandler "my-source/sheet-payment/be/biz/auth"
	"my-source/sheet-payment/be/repository"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Logger struct {
	repo repository.ILogging
}

func NewLogger(repo repository.ILogging) *Logger {
	return &Logger{repo: repo}
}

func (lg *Logger) LogUserActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		user := "anonymous"
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(authenhandler.JwtSecret), nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if username, ok := claims["username"].(string); ok {
						user = username
					}
				}
			}
		}

		log.Printf("[User: %s] %s %s at %s", user, c.Method(), c.Path(), start.Format(time.RFC3339))
		lWrite := repository.UserLog{
			Username: user,
			Method:   c.Method(),
			Path:     c.Path(),
		}

		_ = lg.repo.Write(lWrite)
		return c.Next()
	}
}

func (lg *Logger) GetLogs(c *fiber.Ctx) error {
	logs, err := lg.repo.GetAllLogs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch logs",
		})
	}

	return c.JSON(logs)
}
