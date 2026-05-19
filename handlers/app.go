package handlers

import (
	"github.com/snrkostik/cyberclass/services"
	"github.com/snrkostik/cyberclass/store"
)

type App struct {
	store *store.SQLiteStore

	bracketService   *services.BracketService
	matchService     *services.MatchService
	standingsService *services.StandingsService
	timerService     *services.TimerService
}

func NewApp(
	store *store.SQLiteStore,
) *App {

	return &App{
		store: store,

		bracketService:   services.NewBracketService(store),
		matchService:     services.NewMatchService(store),
		standingsService: services.NewStandingsService(store),
		timerService:     services.NewTimerService(),
	}
}
