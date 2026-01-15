package middleware

import (
	"errors"

	"github.com/3blakejohnson/survivor-fantasy-league/models"
	"github.com/3blakejohnson/survivor-fantasy-league/services"
	"github.com/gofiber/fiber/v2"
)

// Auth returns a Fiber middleware handler skeleton for authentication.
func Auth(as services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(models.ACCESS_TOKEN_COOKIE)
		if token == "" {
			token = c.Get(models.AUTH_HEADER)
		}

		claims, err := as.VerifyJWT(token)
		if err != nil {
			if errors.Is(err, models.ErrJWTExpired) {
				// TODO: refresh
			} else {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("access_level", claims.AccessLevel)
		return c.Next()
	}
}
