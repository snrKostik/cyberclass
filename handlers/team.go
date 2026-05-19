package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/snrkostik/cyberclass/models"
	"github.com/snrkostik/cyberclass/templates"
)

func (a *App) AddTeam(
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

	name := c.FormValue("name")
	slogan := c.FormValue("slogan")

	if name == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"name required",
		)
	}

	team := &models.Team{
		Name:   name,
		Slogan: slogan,
	}

	err = a.store.CreateTeam(
		context.Background(),
		team,
	)

	if err != nil {
		return err
	}

	teams, err := a.store.
		GetTeamsByTournament(
			context.Background(),
			tournamentID,
		)

	if err != nil {
		return err
	}

	seedValue := len(teams) + 1

	err = a.store.AddTeamToTournament(
		context.Background(),
		&models.TournamentTeam{
			TournamentID: tournamentID,
			TeamID:       team.ID,
			Seed:         &seedValue,
			JoinedAt:     time.Now().Unix(),
		},
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

	c.Type("html")

	return templates.AdminContent(
		tournamentID,
		updatedTeams,
		nil,
		false,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}

func (a *App) DeleteTeam(
	c *fiber.Ctx,
) error {

	tournamentID, err := strconv.ParseInt(
		c.Params("tournamentID"),
		10,
		64,
	)

	if err != nil {
		return err
	}

	teamID, err := strconv.ParseInt(
		c.Params("teamID"),
		10,
		64,
	)

	if err != nil {
		return err
	}

	hasBracket, err := a.store.
		TournamentHasBracket(
			context.Background(),
			tournamentID,
		)

	if err != nil {
		return err
	}

	if hasBracket {

		return fiber.NewError(
			fiber.StatusBadRequest,
			"cannot delete teams after bracket generation",
		)
	}

	err = a.store.DeleteTeam(
		context.Background(),
		teamID,
	)

	if err != nil {
		return err
	}

	teams, err := a.store.
		GetTeamsByTournament(
			context.Background(),
			tournamentID,
		)

	if err != nil {
		return err
	}

	c.Type("html")

	return templates.AdminContent(
		tournamentID,
		teams,
		nil,
		false,
	).Render(
		context.Background(),
		c.Response().BodyWriter(),
	)
}
