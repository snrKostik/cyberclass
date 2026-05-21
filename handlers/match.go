package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/snrkostik/cyberclass/templates"
)

type CompleteMatchRequest struct {
	Score1 int `json:"score1"`
	Score2 int `json:"score2"`
}

func (a *App) CompleteMatch(
	c *fiber.Ctx,
) error {

	matchID, err := strconv.ParseInt(
		c.Params("id"),
		10,
		64,
	)

	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid match id",
		)
	}

	var req CompleteMatchRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	err = a.matchService.CompleteMatch(
		context.Background(),
		matchID,
		req.Score1,
		req.Score2,
	)
	a.timerService.Start30Seconds()

	if err != nil {
		return err
	}

	match, err := a.store.GetMatchByID(
		context.Background(),
		matchID,
	)

	if err != nil {
		return err
	}

	rounds, err := a.store.GetBracketView(
		context.Background(),
		match.TournamentID,
	)

	if err != nil {
		return err
	}

	c.Type("html")

	return templates.Bracket(
		match.TournamentID,
		rounds,
		true,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}
