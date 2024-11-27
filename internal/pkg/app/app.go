package app

import (
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/client"
	"jedyEvgeny/online-music-library/internal/app/endpoint"
	"jedyEvgeny/online-music-library/internal/app/service"
	"jedyEvgeny/online-music-library/internal/config"
	storage "jedyEvgeny/online-music-library/internal/storage/postgres"
	"jedyEvgeny/online-music-library/pkg/logger"
	"net/http"
	"strings"
)

type App struct {
	cfg         *config.Config
	log         *logger.Logger
	db          *storage.DataBase
	client      *client.Client
	service     *service.Service
	endpoint    *endpoint.Endpoint
	routeClient *routeClient
	routeServer *routeServer
}

func New() *App {
	a := &App{}
	a.cfg = config.MustLoad()
	a.log = logger.New(a.cfg.Server.LogLevel)

	a.db = storage.MustNew(a.cfg, a.log)

	a.routeClient = newRouteClient()
	a.routeServer = newRouteServer()

	a.client = client.New(a.cfg.Client.Host, a.cfg.Client.Port, a.log, a.routeClient.GetSong)
	a.service = service.New(a.log, a.db, a.client, a.db)
	a.endpoint = endpoint.New(a.log, a.service, a.service)

	return a
}

func (a *App) Run() error {
	defer func() {
		if err := a.db.Close(); err != nil {
			a.log.Info(fmt.Sprintf("Ошибка при закрытии базы данных: %v", err))
		}
	}()

	a.configureRoutes()

	a.log.Info(fmt.Sprintf("Запустили сервер на хосте: %s и порту: %s\n%s\n",
		a.cfg.Server.Host, a.cfg.Server.Port, strings.Repeat("-", 78)))

	err := http.ListenAndServe(a.serverAdress(), nil)
	if err != nil {
		return fmt.Errorf("ошибка прослушивания порта: %w", err)
	}
	return nil
}

func (a *App) serverAdress() string {
	return a.cfg.Server.Host + ":" + a.cfg.Server.Port
}
