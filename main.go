package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/snrkostik/cyberclass/handlers"
	"github.com/snrkostik/cyberclass/store"
)

func main() {

	db, err := store.Open("tournament.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	st := store.NewSQLiteStore(db)

	appHandler := handlers.NewApp(st)

	app := fiber.New()

	app.Get(
		"/",
		appHandler.HomePage,
	)

	app.Post(
		"/tournaments",
		appHandler.CreateTournament,
	)

	app.Get(
		"/tournaments/new",
		appHandler.NewTournamentPage,
	)

	app.Post(
		"/tournaments/:id/bracket",
		appHandler.GenerateBracket,
	)

	app.Post(
		"/tournaments/:id/generate",
		appHandler.GenerateBracketFromTeams,
	)

	app.Post(
		"/tournaments/:id/teams",
		appHandler.AddTeam,
	)

	app.Delete(
		"/tournaments/:tournamentID/teams/:teamID",
		appHandler.DeleteTeam,
	)

	app.Post(
		"/tournaments/:id/delete",
		appHandler.DeleteTournament,
	)

	app.Get(
		"/tournaments/:id/matches",
		appHandler.GetTournamentMatches,
	)

	app.Post(
		"/matches/:id/complete",
		appHandler.CompleteMatch,
	)

	app.Get(
		"/admin/:id",
		appHandler.AdminPage,
	)

	app.Get(
		"/display/:id",
		appHandler.DisplayPage,
	)

	app.Get(
		"/display/:id/content",
		appHandler.DisplayContent,
	)

	log.Fatal(app.Listen(":3000"))
}
