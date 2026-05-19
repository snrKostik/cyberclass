package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/snrkostik/cyberclass/templates"
)

func (a *App) HomePage(
	c *fiber.Ctx,
) error {

	tournaments, err := a.store.GetAllTournaments(
		context.Background(),
	)

	if err != nil {
		return err
	}

	c.Type("html")

	return templates.HomePage(
		tournaments,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}
