package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/snrkostik/cyberclass/templates"
)

func (a *App) DisplayPage(
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

	rounds, err := a.store.GetBracketView(
		context.Background(),
		tournamentID,
	)
	if err != nil {
		return err
	}

	standings, err := a.standingsService.
		ComputeStandings(
			context.Background(),
			tournamentID,
		)

	if err != nil {
		return err
	}

	secondsLeft := a.timerService.SecondsLeft()

	c.Type("html")

	return templates.DisplayPage(
		tournamentID,
		rounds,
		standings,
		secondsLeft,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}

func (a *App) DisplayContent(
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

	rounds, err := a.store.GetBracketView(
		context.Background(),
		tournamentID,
	)
	if err != nil {
		return err
	}

	standings, err := a.standingsService.
		ComputeStandings(
			context.Background(),
			tournamentID,
		)
	if err != nil {
		return err
	}
	secondsLeft := a.timerService.SecondsLeft()

	c.Type("html")

	return templates.DisplayContent(
		tournamentID,
		rounds,
		standings,
		secondsLeft,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}
