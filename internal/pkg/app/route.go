package app

import "net/http"

type routeServer struct {
	GetListSongs string
	GetLyrics    string
	DeleteSong   string
	UpdateSong   string
	AddSong      string
}

type routeClient struct {
	GetSong string
}

func newRouteServer() *routeServer {
	return &routeServer{
		GetListSongs: "/songs-list/",
		GetLyrics:    "/song-find/",
		DeleteSong:   "/song-del",
		UpdateSong:   "/song-up/",
		AddSong:      "/song-add",
	}
}

func newRouteClient() *routeClient {
	return &routeClient{
		GetSong: "/info/",
	}
}

func (a *App) configureRoutes() {
	http.HandleFunc(a.routeServer.AddSong, a.endpoint.HandlerAddSong)
	http.HandleFunc(a.routeClient.GetSong, emulateResponseFromRemoteService)

	http.HandleFunc(a.routeServer.DeleteSong, a.endpoint.HandlerDeleteSong)
}
