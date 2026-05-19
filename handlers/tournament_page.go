package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/snrkostik/cyberclass/templates"
)

func (a *App) NewTournamentPage(
	c *fiber.Ctx,
) error {

	c.Type("html")

	return templates.NewTournamentPage().
		Render(
			context.Background(),
			c.Response().BodyWriter(),
		)
}
