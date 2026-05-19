package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/snrkostik/cyberclass/models"
)

type CreateTournamentRequest struct {
	Name string `json:"name"`
	Game string `json:"game"`
}

func (a *App) CreateTournament(
	c *fiber.Ctx,
) error {

	var req CreateTournamentRequest

	req.Name = c.FormValue("name")
	req.Game = c.FormValue("game")

	// if err := c.BodyParser(&req); err != nil {
	// 	return fiber.NewError(
	// 		fiber.StatusBadRequest,
	// 		err.Error(),
	// 	)
	// }

	tournament := &models.Tournament{
		Name: req.Name,
		Game: req.Game,

		Format: models.TournamentFormatSingleElimination,

		Status: models.TournamentStatusPending,

		CreatedAt: time.Now().Unix(),
	}

	err := a.store.CreateTournament(
		context.Background(),
		tournament,
	)

	if err != nil {
		return err
	}

	return c.Redirect(
		fmt.Sprintf(
			"/admin/%d",
			tournament.ID,
		),
	)
}

func (a *App) DeleteTournament(
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

	err = a.store.DeleteTournament(
		context.Background(),
		tournamentID,
	)

	if err != nil {
		return err
	}

	return c.Redirect("/")
}
