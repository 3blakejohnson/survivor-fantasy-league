package controllers

import (
	"strings"

	"github.com/3blakejohnson/survivor-fantasy-league/templates"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

// GetAppPage handles GET /app/:page
func GetAppPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.Params("page")
		var (
			content templ.Component
			title   string
		)

		switch page {
		case "", "home":
			title = "Home"
			content = templates.HomePage()
		case "login":
			title = "Login"
			content = templates.LoginPage()
		case "leagues":
			title = "Leagues"
			content = templates.LeagueTable()
		default:
			title = strings.Title(page)
			content = templates.MissingPage(page)
		}

		c.Set("Content-Type", "text/html")
		if c.Get("HX-Request") == "true" {
			c.Set("HX-Target", "#app-main")
			return content.Render(c.Context(), c.Response().BodyWriter())
		}
		return templates.AppPage(title, content).Render(c.Context(), c.Response().BodyWriter())
	}
}
