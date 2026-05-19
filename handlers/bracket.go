package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/snrkostik/cyberclass/templates"
)

type GenerateBracketRequest struct {
	TeamIDs []int64 `json:"team_ids"`
}

func (a *App) GenerateBracket(
	c *fiber.Ctx,
) error {

	tournamentID, err := strconv.ParseInt(
		c.Params("id"),
		10,
		64,
	)

	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid tournament id",
		)
	}

	var req GenerateBracketRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	err = a.bracketService.
		CreateSingleEliminationBracket(
			context.Background(),
			tournamentID,
			req.TeamIDs,
		)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func (a *App) GetTournamentMatches(
	c *fiber.Ctx,
) error {

	tournamentID, err := strconv.ParseInt(
		c.Params("id"),
		10,
		64,
	)

	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid tournament id",
		)
	}

	matches, err := a.store.
		GetMatchesByTournament(
			context.Background(),
			tournamentID,
		)

	if err != nil {
		return err
	}

	return c.JSON(matches)
}

func (a *App) BracketPartial(
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

	return templates.Bracket(
		tournamentID,
		rounds,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}

func (a *App) GenerateBracketFromTeams(
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

	if len(teams) != 16 {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"16 teams required",
		)
	}

	var teamIDs []int64

	for _, team := range teams {
		teamIDs = append(
			teamIDs,
			team.ID,
		)
	}

	err = a.bracketService.
		CreateSingleEliminationBracket(
			context.Background(),
			tournamentID,
			teamIDs,
		)

	if err != nil {
		return err
	}

	updatedTeams, err := a.store.
		GetTeamsByTournament(
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

	return templates.AdminContent(
		tournamentID,
		updatedTeams,
		rounds,
		true,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}
