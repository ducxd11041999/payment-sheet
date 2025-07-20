package middlewarelogging

import (
	"log"
	"my-source/sheet-payment/be/repository"
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
		if token := c.Locals("user"); token != nil {
			if claims, ok := token.(*jwt.Token); ok {
				if claimsMap, ok := claims.Claims.(jwt.MapClaims); ok {
					if username, ok := claimsMap["username"].(string); ok {
						user = username
					}
				}
			}
		}

		log.Printf("[User: %s] %s %s at %s", user, c.Method(), c.Path(), start.Format(time.RFC3339))
		lWrite := repository.UserLog{
			Username: c.Method(),
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
