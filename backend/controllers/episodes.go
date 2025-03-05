package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/3blakejohnson/survivor-fantasy-league/backend/dao"
	"github.com/3blakejohnson/survivor-fantasy-league/backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetUser handles GET /episode/:season/:episode
func GetEpisode(dm dao.DAOManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		seasonStr := c.Params("season")
		episodeStr := c.Params("episode")
		seasonNum, err := strconv.ParseInt(seasonStr, 10, 64)
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "Failed to convert season to int",
			})
		}
		episodeNum, err := strconv.ParseInt(episodeStr, 10, 64)
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "Failed to convert episode to int",
			})
		}

		episode, err := dm.Episode().Get(seasonNum, episodeNum)
		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("Failed to fetch Episode: %s", err.Error()),
			})
		}
		return c.JSON(episode)
	}
}

type EpisodePayload struct {
	Title   string `json:"title"`
	Season  int64  `json:"season"`
	Episode int64  `json:"episode"`
	AirDate string `json:"air_date"` // JSON sends date as a string
}

// POST /episode
func CreateEpisode(dm dao.DAOManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var episodePayload EpisodePayload

		// Parse Request Body
		if err := c.BodyParser(&episodePayload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid request payload",
			})
		}

		parsedDate, err := time.Parse("2006-01-02", episodePayload.AirDate)
		if err != nil {
			fmt.Println("Error parsing AirDate:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid date format for air_date. Expected YYYY-MM-DD.",
			})
		}

		episode := models.Episode{
			Title:   episodePayload.Title,
			Season:  episodePayload.Season,
			Episode: episodePayload.Episode,
			AirDate: parsedDate,
		}

		err = dm.Episode().Create(episode)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create Episode",
			})
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"error":   "Successfully Created Episode",
			})
		}
	}
}
