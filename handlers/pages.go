package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/snrkostik/cyberclass/templates"
)

func (a *App) AdminPage(
	c *fiber.Ctx,
) error {

	tournamentID, err := strconv.ParseInt(
		c.Params("id"),
		10,
		64,
	)
	if err != nil {
		return err
	}

	teams, err := a.store.GetTeamsByTournament(
		context.Background(),
		tournamentID,
	)
	if err != nil {
		return err
	}

	hasBracket, err := a.store.TournamentHasBracket(
		context.Background(),
		tournamentID,
	)
	if err != nil {
		return err
	}

	rounds, err := a.store.GetBracketView(
		context.Background(),
		tournamentID,
	)
	if err != nil {
		return err
	}

	c.Type("html")

	return templates.AdminPage(
		tournamentID,
		teams,
		rounds,
		hasBracket,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}
