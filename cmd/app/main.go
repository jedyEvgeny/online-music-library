package main

import (
	"jedyEvgeny/online-music-library/internal/pkg/app"
	"log"
)

// @title						Онлайн-библиотека музыки
// @version					1.0
// @description				Проект для Effective Mobile
// @contact.name				Евгений
// @contact.url				https://github.com/jedyEvgeny
// @contact.email				KEF1991@yandex.ru
// @license.name				MIT
// @license.url				http://opensource.org/licenses/MIT
// @host						localhost:8080
// @BasePath					/
// @accept						json
// @produce					json text/plain
// @schemes					http
// @externalDocs.description	"Readme на GitHub"
// @externalDocs.url			https://github.com/jedyEvgeny/online-music-library/blob/main/README.MD
// @x-name						{"environment": "production", "version": "1.0.0", "team": "backend"}
// @tag.name					Music-library
// @tag.description			Хранилище информации о музыкальных произведениях
// @tag.docs.url				https://t.me/EvKly
// @tag.docs.description		Связь с автором в Телеграм
func main() {
	a := app.New()
	err := a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
