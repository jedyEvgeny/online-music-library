package app

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type routeServer struct {
	Library    string
	Lyrics     string
	DeleteSong string
	UpdateSong string
	AddSong    string
	Swagger    string
}

type routeClient struct {
	GetSong string
}

const pathSwaggerJson = "./docs/swagger.json"

func newRouteServer() *routeServer {
	return &routeServer{
		Library:    "/list/",
		Lyrics:     "/lyrics/",
		DeleteSong: "/song-del/",
		UpdateSong: "/song-upd/",
		AddSong:    "/song-add",
		Swagger:    "/swagger/doc.json",
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

	http.HandleFunc(a.routeServer.Lyrics, a.endpoint.HandlerLiryc)
	http.HandleFunc(a.routeServer.DeleteSong, a.endpoint.HandlerDeleteSong)
	http.HandleFunc(a.routeServer.Library, a.endpoint.HandlerLibrary)
	http.HandleFunc(a.routeServer.UpdateSong, a.endpoint.HandlerPatchSong)

	http.HandleFunc(a.routeServer.Swagger, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, pathSwaggerJson)
	})
	http.Handle("/swagger/", httpSwagger.WrapHandler)
}
