package handlers

import (
	"github.com/snrkostik/cyberclass/services"
	"github.com/snrkostik/cyberclass/store"
)

type App struct {
	store *store.SQLiteStore

	bracketService *services.BracketService
	matchService   *services.MatchService
}

func NewApp(
	store *store.SQLiteStore,
) *App {

	return &App{
		store: store,

		bracketService: services.NewBracketService(store),
		matchService:   services.NewMatchService(store),
	}
}
