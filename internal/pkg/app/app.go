package app

import "jedyEvgeny/online-music-library/internal/config"

type App struct {
	cfg         *config.Config
	routeClient *routeClient
	routeServer *routeServer
}

func New() *App {
	a := &App{}
	a.cfg = config.MustLoad()
	a.routeClient = newRouteClient()
	a.routeServer = newRouteServer()
	return a
}

func (a *App) Run() error {
	return nil
}
