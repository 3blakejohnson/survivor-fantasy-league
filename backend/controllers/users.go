package controllers

import (
	"fmt"
	"strconv"

	"github.com/3blakejohnson/survivor-fantasy-league/backend/dao"
	"github.com/gofiber/fiber/v2"
)

// GetUsers handles GET /users
func GetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "List of users"})
}

// GetUser handles GET /users/:id
func GetUser(dm dao.DAOManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "Failed to convert id to int",
			})
		}

		user, err := dm.User().Get(intId)
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("Failed to fetch user: %s", err.Error()),
			})
		}
		return c.JSON(user)
	}
}
