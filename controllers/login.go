package controllers

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/3blakejohnson/survivor-fantasy-league/models"
	"github.com/3blakejohnson/survivor-fantasy-league/services"
	"github.com/gofiber/fiber/v2"
)

type LoginPayload struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// CreateUser handles POST /users or /signup
func CreateUser(as services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger, ok := c.Locals("logger").(*log.Logger)
		if !ok || logger == nil {
			logger = log.Default()
		}
		var payload models.User
		_ = c.BodyParser(&payload)
		if payload.Username == "" {
			payload.Username = c.FormValue("username")
		}
		if payload.PasswordHash == "" {
			payload.PasswordHash = c.FormValue("password")
		}
		if payload.Username == "" || payload.PasswordHash == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Username and password are required",
			})
		}

		_, err := as.CreateUser(payload)
		if err != nil {
			logger.Printf("create user failed: %s", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
		})
	}
}

// Login handles POST /login
func Login(as services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger, ok := c.Locals("logger").(*log.Logger)
		if !ok || logger == nil {
			logger = log.Default()
		}
		var payload LoginPayload
		_ = c.BodyParser(&payload)
		if payload.Username == "" {
			payload.Username = c.FormValue("username")
		}
		if payload.Password == "" {
			payload.Password = c.FormValue("password")
		}
		if payload.Username == "" || payload.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Username and password are required",
			})
		}

		user, err := as.ValidateLogin(payload.Username, payload.Password)
		if err != nil {
			logger.Printf("failed to validate login: %s", err.Error())
			if errors.Is(err, models.ErrInvalidCredentials) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"error":   "Invalid username or password",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to validate login",
			})
		}

		token, err := as.CreateJWT(strconv.FormatInt(user.ID, 10), "user", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create access token",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     models.ACCESS_TOKEN_COOKIE,
			Value:    token,
			Path:     "/",
			HTTPOnly: true,
			SameSite: "Lax",
		})

		accept := c.Get(fiber.HeaderAccept)
		contentType := c.Get(fiber.HeaderContentType)
		wantsJSON := strings.Contains(accept, fiber.MIMEApplicationJSON) ||
			strings.Contains(contentType, fiber.MIMEApplicationJSON)
		if wantsJSON {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
			})
		}

		return c.Redirect("/app/home", fiber.StatusFound)
	}
}
