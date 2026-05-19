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

	if name == "" {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"name required",
		)
	}

	team := &models.Team{
		Name: name,

		CreatedAt: time.Now().Unix(),
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
